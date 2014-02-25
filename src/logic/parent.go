// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 17:53
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 17:53 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
	"strconv"
	"fmt"
)

//根据电话获取parent表的id
func getParentIdByPhone(phone string) (int64, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select pid from parent where telephone=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, phone)

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

func insertParent(tx *sql.Tx, name, password, telephone, comeForm string) (id int64, err error) {

	sql := "insert into parent(name,password,telephone,reg_date,come_form) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(name, password, telephone, time.Now().Format("20060102150405"), comeForm)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newParentId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newParentId, nil
}

func MemberPage(dataType,centerId,kw,employeeId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
	if centerId == "1" {
		centerId = "7"
	}

	params := []interface{}{}

	dataSql := `
				select p.pid id,ce.name centerName,p.telephone,p.father_name fatherName,p.father_phone fatherPhone,mother_name motherName,mother_phone motherPhone,b.childs
				from parent p
				left join consumer_new cons on cons.parent_id=p.pid
				left join center ce on cons.center_id=ce.cid
				left join (select pid,GROUP_CONCAT(name   ORDER BY cid DESC SEPARATOR ',') childs from child group by pid) b on b.pid=p.pid
				where p.is_member=1
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
		dataSql += " and cons.center_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		dataSql += " and cons.center_id=? "
	}

	if kw != "" {
		dataSql += " and (b.childs like ? or p.telephone like ? or p.father_name like ? or p.father_phone like ? or p.mother_name like ? or p.mother_phone like ?)"
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}

	countSql := "select count(1) from (" + dataSql + ") num"
	lessgo.Log.Debug(countSql)
	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, params)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql += " order by p.pid desc limit ?,? "
	lessgo.Log.Debug(dataSql)
	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, params)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	return pageData, nil
}

func updateParent(tx *sql.Tx, parentDataMap map[string]interface{}, id string) error {
	sql := "update parent set %v where pid=?"
	params := []interface{}{}

	setSql := ""

	for key, value := range parentDataMap {
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
