// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-17 15:58
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-17 15:58 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
)

/*
select st.id,st.room_id,st.time_id,st.week,csd.start_time,csd.course_id
			from class_schedule_detail csd
			left join schedule_template st on csd.center_id=st.center_id and csd.room_id=st.room_id and csd.time_id=st.time_id and csd.week=st.week
			where csd.id=?
*/
func getScheduleTmpByScheduelDetailId(scheduleId string) (map[string]string, error) {

	sql := `
			select st.id,st.room_id,st.time_id,st.week,csd.start_time,csd.course_id
			from class_schedule_detail csd
			left join schedule_template st on csd.center_id=st.center_id and csd.room_id=st.room_id and csd.time_id=st.time_id and csd.week=st.week
			where csd.id=?`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, scheduleId)

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

func getScheduleTmpsByCenterId(centerId string) ([]map[string]string, error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	getScheduleTmpSql := "select id,room_id,teacher_id,assistant_id,time_id,week,course_id from schedule_template where center_id=? "
	lessgo.Log.Debug(getScheduleTmpSql)

	rows, err := db.Query(getScheduleTmpSql, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var datas []map[string]string

	for rows.Next() {
		dataMap, err := lessgo.GetDataMap(rows)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}

		datas = append(datas, dataMap)
	}

	return datas, nil
}

func getChildAndContractByScheduleTempId(tempId string) ([]map[string]string, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select child_id,contract_id from schedule_template_child where schedule_template_id=? "
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, tempId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var datas []map[string]string

	for rows.Next() {
		dataMap, err := lessgo.GetDataMap(rows)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}

		datas = append(datas, dataMap)
	}

	return datas, nil
}

func checkScheduleTempChildExist(childId, scheduleTempId string) (flag bool, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := " select count(1) from schedule_template_child where child_id=? and  schedule_template_id=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, childId, scheduleTempId)

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

func insertScheduleTempChild(tx *sql.Tx, childId, scheduleTempId, contractId string) error {
	sql := "insert into schedule_template_child(schedule_template_id,child_id,contract_id) values(?,?,?) "
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(scheduleTempId, childId, contractId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func deleteScheduleTempChild(tx *sql.Tx, childId, scheduleTempId string) error {
	sql := "delete from schedule_template_child where schedule_template_id=? and child_id=? "
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(scheduleTempId, childId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}
