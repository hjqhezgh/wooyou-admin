// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-18 16:12
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-18 16:12 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
)

const (
	SIGN_IN_SUCCESS = "1" //正常签到
	SIGN_IN_LEAVE   = "2" //请假
	SIGN_IN_TRUANT  = "3" //旷课
)

const (
	IS_FREE_YES = "1" //是试听课
	IS_FREE_NO  = "2" //不是试听课
)

//如果scheduleId=“”，则视为无班签到的判断
func checkSignInExist(childId, scheduleId string) (bool, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	dataSql := ""

	if scheduleId=="" {
		dataSql = "select count(1) from sign_in where child_id=? and schedule_detail_id is null "
	}else{
		dataSql = "select count(1) from sign_in where child_id=? and schedule_detail_id=? "
	}

	lessgo.Log.Debug(dataSql)

	var rows *sql.Rows
	var err error

	if scheduleId=="" {
		rows, err = db.Query(dataSql, childId)
	}else{
		rows, err = db.Query(dataSql, childId, scheduleId)
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	return false, nil
}

//scheduleId==""为无班签到
func insertSignIn(tx *sql.Tx, scheduleId, childId, signType, contractId, cardId, employeeId, isFree string) (err error) {

	sql := ""

	if scheduleId == ""{
		sql = "insert into sign_in(child_id,sign_time,type,contract_id,card_id,employee_id,is_free) values(?,?,?,?,?,?,?)"
	}else{
		sql = "insert into sign_in(child_id,sign_time,schedule_detail_id,type,contract_id,card_id,employee_id,is_free) values(?,?,?,?,?,?,?,?)"
	}

	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	if scheduleId == ""{
		_, err = stmt.Exec(childId, time.Now().Format("20060102150405"), signType, contractId, cardId, employeeId, isFree)
	}else{
		_, err = stmt.Exec(childId, time.Now().Format("20060102150405"), scheduleId, signType, contractId, cardId, employeeId, isFree)
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}
