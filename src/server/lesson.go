// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-12 10:15
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-12 10:15 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
)

type Lesson struct {
	Lid        int
	CourseId   int
	Caption    string
	TeacherId  int
	KeyStone   string
	LessonTime int
	OrderNo    int
	IsLast     string
}

func LessonByClassIdAction(w http.ResponseWriter, r *http.Request) {

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

	sql := "select les.lid,les.caption from lesson les left join course cour on cour.cid=les.course_id left join wyclass wy on wy.course_id=cour.cid where wy.class_id=? "

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

		err = commonlib.PutRecord(rows, &model.Value, &model.Desc)

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

func FindLessonById(id string) (Lesson, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select lid,course_id,caption,teacher_id,keystone,lesson_time,order_no,is_last from lesson where lid=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		return Lesson{}, err
	}

	lesson := Lesson{}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &lesson.Lid, &lesson.CourseId, &lesson.Caption, &lesson.TeacherId, &lesson.KeyStone, &lesson.LessonTime, &lesson.OrderNo, &lesson.IsLast)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return Lesson{}, err
		}
	}

	return lesson, nil
}
