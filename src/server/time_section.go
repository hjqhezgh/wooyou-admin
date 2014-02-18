// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-12 00:25
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-12 00:25 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
)

type TimeSection struct {
	Id        int
	CenterId  int
	StartTime string
	EndTime   string
	LessonNo  int
}

func TimeSectionByCenterIdAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	centerId := r.FormValue("id")

	sql := "select id,start_time,end_time from time_section where center_id=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, centerId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	type Result struct {
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}

	objects := []*Result{}

	for rows.Next() {

		model := new(Result)
		var st, et string

		err = commonlib.PutRecord(rows, &model.Value, &st, &et)

		model.Desc = st + "-" + et

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}
	m["success"] = true
	m["datas"] = objects

	commonlib.OutputJson(w, m, " ")
	return
}

func FindTimeSectionById(id string) (TimeSection, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select id,center_id,start_time,end_time,lesson_no from time_section where id=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		return TimeSection{}, err
	}

	timeSection := TimeSection{}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &timeSection.Id, &timeSection.CenterId, &timeSection.StartTime, &timeSection.EndTime, &timeSection.LessonNo)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return TimeSection{}, err
		}
	}

	return timeSection, nil
}
