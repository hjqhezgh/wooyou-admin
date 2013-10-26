// Title：电话录音备注信息
//
// Description:
//
// Author:black
//
// Createtime:2013-10-09 17:25
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-09 17:25 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
)

//加载通话记录备注
func AudioNoteLoadAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := lessgo.GetCurrentEmployee(r)

	if employee.UserId == "" {
		lessgo.Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select a.aid,a.note,a.start_time,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone from (select * from audio where aid = ?) a	left join consumer c on (a.remotephone=c.mother_phone and c.mother_phone!='' and c.mother_phone is not null ) or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null)"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	var aid, note, startTime, mother, motherPhone, father, fatherPhone, homePhone string

	if rows.Next() {
		err := commonlib.PutRecord(rows, &aid, &note, &startTime, &mother, &motherPhone, &father, &fatherPhone, &homePhone)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	m["success"] = true

	loadFormObjects := []lessgo.LoadFormObject{}

	h1 := lessgo.LoadFormObject{"aid", aid}
	h2 := lessgo.LoadFormObject{"father", father}
	h3 := lessgo.LoadFormObject{"fatherPhone", fatherPhone}
	h4 := lessgo.LoadFormObject{"mother", mother}
	h5 := lessgo.LoadFormObject{"motherPhone", motherPhone}
	h6 := lessgo.LoadFormObject{"homePhone", homePhone}
	h7 := lessgo.LoadFormObject{"note", note}
	h8 := lessgo.LoadFormObject{"startTime", startTime}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")

}

//保存通话记录备注
func AudioNoteSaveAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := lessgo.GetCurrentEmployee(r)

	if employee.UserId == "" {
		lessgo.Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("aid")
	note := r.FormValue("note")

	sql := "update audio set note=? where aid=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(note, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	commonlib.OutputJson(w, m, " ")
}
