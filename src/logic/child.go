// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-08 11:31
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-08 11:31 black 创建文档
package logic

import (
	"database/sql"
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"strconv"
	"time"
)

func insertChild(tx *sql.Tx, name, pid, sex, birthday, year, month, centerId string) (id int64, err error) {

	sql := "insert into child(name,pid,sex,birthday,center_id) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	childBirthday := "20090101"
	if birthday != "" {
		childBirthday = birthday
	} else {
		if year != "" && month != "" {
			monthInt, _ := strconv.Atoi(month)
			if monthInt > 9 {
				childBirthday = year + month + "01"
			} else {
				childBirthday = year + "0" + month + "01"
			}
		} else if year != "" {
			childBirthday = year + "0101"
		}
	}

	res, err := stmt.Exec(name, pid, sex, childBirthday, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	childId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return childId, nil
}

/*
select cid id,name,card_id cardId,pid,sex,birthday,center_id centerId,avatar from child where cid=?
*/
func getChildById(id string) (map[string]string, error) {

	sql := `select cid id,name,card_id cardId,pid,sex,birthday,center_id centerId,avatar
	        from child where cid=?`

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

func updateChild(tx *sql.Tx, childDataMap map[string]interface{}, id string) error {
	sql := "update child set %v where cid=?"
	params := []interface{}{}

	setSql := ""

	for key, value := range childDataMap {
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

//todo 目前只返回第一个孩子的id，逻辑有待优化
func getChildByParentId(pid string) (int64, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select cid from child where pid=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, pid)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	var id int64

	if rows.Next() {

		err = commonlib.PutRecord(rows, &id)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return 0, err
		}
	}

	return id, nil
}

//todo 目前只返回第一个孩子的id，逻辑有待优化
func getChildByConsumerId(consumerId string) (int64, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select ch.cid from child ch left join consumer_new cons on cons.parent_id=ch.pid where cons.id=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, consumerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	var id int64

	if rows.Next() {

		err = commonlib.PutRecord(rows, &id)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return 0, err
		}
	}

	return id, nil
}

func ChildInClassPage(scheduleId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from schedule_detail_child where schedule_detail_id=? "
	lessgo.Log.Debug(countSql)
	countParams := []interface{}{scheduleId}

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `
				select sdc.child_id id,ch.name childName,p.telephone phone,e.really_name tmkName,sdc.create_time inviteTime,si.sign_time signTime,sdc.wyclass_id classId,sdc.create_user inviteUser,ch.center_id centerId,sdc.sms_status smsStatus,d.remark remark,cons.level level,cons.id cosumerId,wc.code code ,ch.sex sex,ch.birthday birthday
				from schedule_detail_child sdc
				left join class_schedule_detail csd on csd.id=sdc.schedule_detail_id
				left join wyclass wc on wc.class_id=csd.class_id
	            left join child ch on ch.cid=sdc.child_id
	            left join parent p on p.pid=ch.pid
	            left join consumer_new cons on cons.parent_id=ch.pid
	            left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) d on d.consumer_id=cons.id
				left join employee e on e.user_id=sdc.create_user
				left join sign_in si on si.child_id=sdc.child_id and sdc.schedule_detail_id=si.schedule_detail_id
				where sdc.schedule_detail_id=? order by si.sid desc,sdc.id desc limit ?,?`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{}
	dataParams = append(dataParams, scheduleId)
	dataParams = append(dataParams, (currPageNo-1)*pageSize)
	dataParams = append(dataParams, pageSize)

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func ChildInCenterPage(centerId, kw string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from child ch left join parent p on p.pid=ch.pid where ch.center_id=? and (ch.name like ? or p.telephone like ?) "
	countParams := []interface{}{centerId, "%" + kw + "%", "%" + kw + "%"}
	lessgo.Log.Debug(countSql)

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `
				select c.cid id,c.name as childName,c.telephone phone,cons.level level,d.remark remark,cons.id consumerId from
				(select ch.cid,ch.name,p.telephone,ch.pid from child ch left join parent p on p.pid=ch.pid where ch.center_id=? and (ch.name like ? or p.telephone like ?) order by ch.cid desc limit ?,?) c
				left join consumer_new cons on cons.parent_id=c.pid
	            left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) d on d.consumer_id=cons.id`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{centerId, "%" + kw + "%", "%" + kw + "%", (currPageNo - 1) * pageSize, pageSize}

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func ChildInNormalSchedulePage(scheduleId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from schedule_detail_child where schedule_detail_id=?"
	countParams := []interface{}{scheduleId}
	lessgo.Log.Debug(countSql)

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `
				select sdc.child_id id,ch.name childName,p.telephone phone,si.type signType,si.sign_time signTime,cour.name courseName,contr.id as contractId,contr.contract_no contractNo,contr.apply_time applyTime,cons.id consumerId,cons.level level,d.remark,sdc.is_free isFree,contr.left_lesson_num totalNum,usedNum.num usedNum,csd.center_id centerId
	 		    from (select * from schedule_detail_child where schedule_detail_id=? order by id desc limit ?,?) sdc
	 			left join child ch on ch.cid=sdc.child_id
	 			left join parent p on p.pid=ch.pid
	 			left join consumer_new cons on cons.parent_id=ch.pid
	            left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) d on d.consumer_id=cons.id
	 			left join sign_in si on si.child_id=sdc.child_id and sdc.schedule_detail_id=si.schedule_detail_id
	 			left join class_schedule_detail csd on csd.id=sdc.schedule_detail_id
	 			left join contract contr on contr.id=sdc.contract_id
	 			left join course cour on cour.cid=contr.course_id
	 			left join (select count(1) num,contract_id from sign_in where type=1 or type=3 group by contract_id ) usedNum on usedNum.contract_id=contr.id
	 			order by si.sid desc,ch.cid desc
	 			`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{scheduleId, (currPageNo - 1) * pageSize, pageSize}

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func ChildPage(centerId, contractStatus, kw, dataType, employeeId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
	if centerId == "1" {
		centerId = "7"
	}

	params := []interface{}{}

	dataSql := `
				select ch.cid id,ce.name centerName,ch.name childName,p.telephone phone,p.password,ch.card_id cardId,ch.birthday,ch.sex,p.reg_date regTime,contract_num.num haveContract,totalNum.num totalNum,usedNum.num usedNum
				from child ch
				left join parent p on p.pid=ch.pid
				left join center ce on ce.cid=ch.center_id
				left join (select sum(left_lesson_num) num,child_id from contract group by child_id) totalNum on totalNum.child_id=ch.cid
				left join (select count(1) num,child_id from sign_in where (type=1 or type=3) and contract_id!=0 and contract_id is not null group by child_id) usedNum on usedNum.child_id=ch.cid
				left join (select child_id,count(1) num from contract where price >0 group by child_id) contract_num on contract_num.child_id=ch.cid where 1=1
	`

	if dataType == "center" {
		userId, _ := strconv.Atoi(employeeId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}

		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if _employee.CenterId == "1" {
			_employee.CenterId = "7"
		}

		params = append(params, _employee.CenterId)

		dataSql += " and ch.center_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		dataSql += " and ch.center_id=? "
	}

	if contractStatus != "" {
		dataSql += " and contract_num.num is not null "
	}

	if kw != "" {
		dataSql += " and (ch.name like ? "
		dataSql += " or p.telephone like ?) "
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}

	countSql := "select count(1) from (" + dataSql + ") num"
	lessgo.Log.Debug(countSql)
	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, params)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql += " order by ch.cid desc limit ?,? "

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(dataSql)
	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, params)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func ChildPay(childId, scheduleId, classId, payType, employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	flag, msg, err = childPay(tx, childId, "", scheduleId, classId, payType, employeeId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if !flag {
		return false, msg, nil
	}

	tx.Commit()

	return true, "", nil
}

func childPay(tx *sql.Tx, childId, consumerId, scheduleId, classId, payType, employeeId string) (flag bool, msg string, err error) {

	var consumerDataMap map[string]string

	if consumerId == "" {
		consumerDataMap, err = getConsumerByChildId(childId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
		consumerId = consumerDataMap["id"]
	} else {
		consumerDataMap, err = getConsumerById(consumerId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	if childId == "" {
		childIdInt, err := getChildByConsumerId(consumerId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
		childId = fmt.Sprint(childIdInt)
	}

	err = insertPayLog(tx, consumerId, employeeId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if payType == "all" {
		payType = "2"
	} else if payType == "part" {
		payType = "1"
	}

	consumerUpdateMap := make(map[string]interface{})
	consumerUpdateMap["pay_status"] = payType
	consumerUpdateMap["pay_time"] = time.Now().Format("20060102150405")

	err = updateConsumer(tx, consumerUpdateMap, consumerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if scheduleId == "" { //是在客户主界面点击的缴费
		if consumerDataMap["contactStatus"] != CONSUMER_STATUS_SIGNIN {
			consumerUpdateMap = make(map[string]interface{})
			consumerUpdateMap["contact_status"] = CONSUMER_STATUS_SIGNIN
			consumerUpdateMap["sign_in_time"] = time.Now().Format("20060102150405")

			err = updateConsumer(tx, consumerUpdateMap, consumerId)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			_, err = insertConsumerStatusLog(tx, consumerId, employeeId, consumerDataMap["contactStatus"], CONSUMER_STATUS_SIGNIN)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			scheduleId, classId, err = getNewestFreeScheduleIdByChildId(childId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			if scheduleId == "" { //进行无班签到
				err = insertSignIn(tx, "", fmt.Sprint(childId), SIGN_IN_SUCCESS, "0", "", employeeId, IS_FREE_YES)

				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}

				insertConsumerContactsLog(tx, employeeId, "无班签到", consumerId, CONTACTS_LOG_TYPE_SYSTEM)

				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}
			} else {
				err = insertSignIn(tx, scheduleId, childId, SIGN_IN_SUCCESS, "0", "", employeeId, IS_FREE_YES)
				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}

				note := ""

				if classId == "" { //常规课
					scheduleDataMap, err := getScheduleDetailId(scheduleId)

					if err != nil {
						return false, "", err
					}
					note = fmt.Sprintf("签到常规课[时间]%v[教室]%v[课程]%v中", scheduleDataMap["startTime"], scheduleDataMap["roomName"], scheduleDataMap["courseName"])
				} else {
					classDataMap, err := getWyClassById(classId)

					if err != nil {
						lessgo.Log.Error(err.Error())
						return false, "", err
					}

					note = "签到" + classDataMap["start_time"] + classDataMap["name"]
				}

				insertConsumerContactsLog(tx, employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}
			}

		}
	} else { //在排课界面的缴费
		err = insertSignIn(tx, scheduleId, childId, SIGN_IN_SUCCESS, "0", "", employeeId, IS_FREE_YES)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		consumerUpdateMap = make(map[string]interface{})
		consumerUpdateMap["contact_status"] = CONSUMER_STATUS_SIGNIN
		consumerUpdateMap["sign_in_time"] = time.Now().Format("20060102150405")

		err = updateConsumer(tx, consumerUpdateMap, consumerId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		_, err = insertConsumerStatusLog(tx, consumerId, employeeId, consumerDataMap["contactStatus"], CONSUMER_STATUS_SIGNIN)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		note := ""

		if classId == "" { //常规课
			scheduleDataMap, err := getScheduleDetailId(scheduleId)

			if err != nil {
				return false, "", err
			}
			note = fmt.Sprintf("签到常规课[时间]%v[教室]%v[课程]%v中", scheduleDataMap["startTime"], scheduleDataMap["roomName"], scheduleDataMap["courseName"])
		} else {
			classDataMap, err := getWyClassById(classId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			note = "签到" + classDataMap["start_time"] + classDataMap["name"]
		}

		insertConsumerContactsLog(tx, employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

func PotentialChildPage(centerId, kw, dataType, employeeId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
	if centerId == "1" {
		centerId = "7"
	}

	params := []interface{}{}

	dataSql := `
				select ch.cid id,ce.name centerName,ch.name childName,p.telephone phone,ch.birthday,ch.sex
				from child ch
				left join parent p on p.pid=ch.pid
				left join center ce on ce.cid=ch.center_id
				left join (select child_id,count(1) num from contract where price >0 group by child_id) contract_num on contract_num.child_id=ch.cid where 1=1 and (contract_num.num=0 or contract_num.num is null)
	`

	if dataType == "center" {
		userId, _ := strconv.Atoi(employeeId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}

		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if _employee.CenterId == "1" {
			_employee.CenterId = "7"
		}

		params = append(params, _employee.CenterId)

		dataSql += " and ch.center_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		dataSql += " and ch.center_id=? "
	}

	if kw != "" {
		dataSql += " and (ch.name like ? "
		dataSql += " or p.telephone like ?) "
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}

	countSql := "select count(1) from (" + dataSql + ") num"
	lessgo.Log.Debug(countSql)
	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, params)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql += " order by ch.cid desc limit ?,? "

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(dataSql)
	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, params)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func ChildInParentPage(parentId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	params := []interface{}{}

	dataSql := `
				select ch.cid id,ch.name childName,ch.card_id cardId,ch.birthday,ch.sex,totalNum.num totalNum,usedNum.num usedNum
				from child ch
				left join (select sum(left_lesson_num) num,child_id from contract group by child_id) totalNum on totalNum.child_id=ch.cid
				left join (select count(1) num,child_id from sign_in where (type=1 or type=3) and contract_id!=0 and contract_id is not null group by child_id) usedNum on usedNum.child_id=ch.cid
				where 1=1 and ch.pid=?
	`

	params = append(params, parentId)

	countSql := "select count(1) from (" + dataSql + ") num"
	lessgo.Log.Debug(countSql)
	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, params)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql += " order by ch.cid desc limit ?,? "

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(dataSql)
	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, params)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}
