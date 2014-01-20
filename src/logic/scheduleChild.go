// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-19 23:16
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-19 23:16 black 创建文档
package logic

import (
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
	"database/sql"
)

func getScheduleChildByChildIdAndScheduleId(childId,scheduleId string) (map[string]string, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select id,schedule_detail_id,child_id,create_user,wyclass_id,contract_id,is_free from schedule_detail_child where child_id=? and schedule_detail_id=? "
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, childId,scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var dataMap map[string]string

	if rows.Next() {
		dataMap, err = lessgo.GetDataMap(rows)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}
	}

	return dataMap, nil
}

func deleteScheduleChild(tx *sql.Tx,childId,scheduleId string) error {
	sql := "delete from schedule_detail_child where child_id=? and schedule_detail_id=?"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(childId, scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func getNewestFreeScheduleIdByChildId(childId string) (scheduleId,classId string ,err error){

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select schedule_detail_id,wyclass_id from schedule_detail_child where child_id=? and is_free=1 order by id desc "
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, childId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return "","", err
	}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &scheduleId,&classId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return "","", err
		}
	}

	return scheduleId,classId,nil
}

