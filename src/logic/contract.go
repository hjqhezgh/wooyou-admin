// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-17 14:08
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-17 14:08 black 创建文档
package logic

import (
	"database/sql"
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
)

const (
	CONTRACT_STATUS_WSM = "1" //未上课
	CONTRACT_STATUS_SWZ = "2" //上课中
	CONTRACT_STATUS_YFQ = "3" //已废弃
	CONTRACT_STATUS_YJS = "4" //已结束
)

func ContractCheckIn(childId, scheduleId, contractId, actionType string) (flag bool, msg string, err error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if actionType == "temp" { //临时合同登记
		err = updateContractOfScheduleChild(tx, contractId, childId, scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	} else {
		scheduleTmpgetDataMap, err := getScheduleTmpByScheduelDetailId(scheduleId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if scheduleTmpgetDataMap["id"] == "" {
			return false, "该课表不是模板课表，无法登记跟班合同", nil
		}

		scheduleIds, err := getFurtherScheduleIds(scheduleTmpgetDataMap["time_id"], scheduleTmpgetDataMap["room_id"], scheduleTmpgetDataMap["week"], scheduleTmpgetDataMap["start_time"], scheduleTmpgetDataMap["course_id"])

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		for _, scheduleId := range scheduleIds {
			err = updateContractOfScheduleChild(tx, contractId, childId, scheduleId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}
		}

		err = updateContractOfScheduleTmpChild(tx, contractId, childId, scheduleTmpgetDataMap["id"])

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

/*
select id,child_id,apply_time,contract_no,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status,expire_date from contract where id=?
*/
func getContractById(id string) (map[string]string, error) {

	sql := `select id,child_id,apply_time,contract_no,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status,expire_date from contract where id=?`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var dataMap map[string]string

	if rows.Next() {
		dataMap, err = lessgo.GetDataMap(rows)
	}

	if err != nil {
		return nil, err
	}

	return dataMap, nil
}

//获取指定合同的有效上课次数，这里的有效指的是正常上课和旷课
func getVaildNumOfContract(id string) (int, error) {

	sql := `select count(1) from sign_in where contract_id=? and (type=1 or type=3) `

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	num := 0

	if rows.Next() {
		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return 0, err
		}
	}

	return num, nil
}

//校验合同是否还在三个月的有效期内
func checkContractValid(expireDateString string) (bool, error) {

	if expireDateString == "" { //没有过期时间的，暂时不处理
		return true, nil
	}

	now := time.Now()

	expireDate, err := time.ParseInLocation("20060102150405", expireDateString+"235959", time.Local)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	if now.After(expireDate) {
		return false, nil
	}

	return true, nil
}

/*
select id,child_id,apply_time,contract_no,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status,expire_date from contract where id=?
*/
func getContractIdByChildIdAndScheduleId(childId, scheduleId string) (contractId,isFree string,err error) {

	sql := `select contract_id,is_free from schedule_detail_child where child_id=? and schedule_detail_id=?`
	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, childId, scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return "0","", err
	}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &contractId, &isFree)
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		return "0","", err
	}

	return contractId,isFree, nil
}

func ContractList(childId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := "select count(1) from contract where child_id=?"
	lessgo.Log.Debug(countSql)
	countParams := []interface{}{childId}

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `
		select contr.id id,contr.contract_no contractNo,cour.name courseName,contr.apply_time applyTime,contr.price,contr.left_lesson_num totalNum,e.really_name employeeName,contr.type contractType,contr.status contractStatus,usedNum.num usedNum,leaveNum.num leaveNum,truantNum.num truantNum
		from contract contr
		left join (select count(1) num,contract_id from sign_in where child_id=? and type=1 group by contract_id) usedNum on usedNum.contract_id=contr.id
		left join (select count(1) num,contract_id from sign_in where child_id=? and type=2 group by contract_id) leaveNum on leaveNum.contract_id=contr.id
		left join (select count(1) num,contract_id from sign_in where child_id=? and type=3 group by contract_id) truantNum on truantNum.contract_id=contr.id
		left join employee e on e.user_id=contr.employee_id
		left join course cour on cour.cid=contr.course_id
		where contr.child_id=? limit ?,?
	`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{}
	dataParams = append(dataParams, childId)
	dataParams = append(dataParams, childId)
	dataParams = append(dataParams, childId)
	dataParams = append(dataParams, childId)
	dataParams = append(dataParams, (currPageNo-1)*pageSize)
	dataParams = append(dataParams, pageSize)

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func SaveContract(id, contractNo, price, courseId, courseNum, contractType, childId, expireDate, employeeId, centerId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if id == "" {
		parentId, err := getChildByParentId(childId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		_, err = insertContract(tx, childId, contractNo, fmt.Sprint(parentId), price, employeeId, centerId, courseId, courseNum, contractType, expireDate)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

	} else {

		contractDataMap := make(map[string]interface{})
		contractDataMap["contract_no"] = contractNo
		contractDataMap["price"] = price
		contractDataMap["course_id"] = courseId
		contractDataMap["left_lesson_num"] = courseNum
		contractDataMap["type"] = contractType
		contractDataMap["expire_date"] = expireDate

		err = updateContract(tx, contractDataMap, id)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

func insertContract(tx *sql.Tx, childId, contractNo, parentId, price, employeeId, centerId, courseId, courseNum, contractType, expireDate string) (int64, error) {

	sql := "insert into contract(child_id,apply_time,contract_no,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status,expire_date) values(?,?,?,?,?,?,?,?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(childId, time.Now().Format("20060102150405"), contractNo, parentId, price, employeeId, centerId, courseId, courseNum, contractType, CONTRACT_STATUS_WSM, expireDate)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	contractId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return contractId, nil
}

func updateContract(tx *sql.Tx, contractDataMap map[string]interface{}, id string) error {
	sql := "update contract set %v where id=?"
	params := []interface{}{}

	setSql := ""

	for key, value := range contractDataMap {
		setSql += key + "=?,"
		params = append(params, value)
	}

	params = append(params, id)

	setSql = commonlib.Substr(setSql, 0, len(setSql)-1)

	sql = fmt.Sprintf(sql, setSql)
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(params...)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

/*
select id,child_id,apply_time,contract_no,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status,expire_date from contract where id=?
*/
func GetContractById(id string) (map[string]string, error) {
	return getContractById(id)
}
