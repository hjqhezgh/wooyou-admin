// Title：分配tmk
//
// Description:
//
// Author:black
//
// Createtime:2013-10-16 15:41
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-16 15:41 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
)

func SendToTmkLoadAction(w http.ResponseWriter, r *http.Request) {
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

	sql := "select tc.id,e.user_id,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone,c.child,ce.name from tmk_consumer tc left join consumer c on tc.consumer_id=c.id left join employee e on e.user_id = tc.employee_id left join center ce on ce.cid=tc.center_id where tc.id=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	var tcid,employeeId,mother,motherPhone,father,fatherPhone,homePhone,child,centerName string

	if rows.Next() {
		err := commonlib.PutRecord(rows,&tcid,&employeeId,&mother,&motherPhone,&father,&fatherPhone,&homePhone,&child,&centerName)

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

	h1 := lessgo.LoadFormObject{"tcid", tcid}
	h2 := lessgo.LoadFormObject{"father", father}
	h3 := lessgo.LoadFormObject{"fatherPhone", fatherPhone}
	h4 := lessgo.LoadFormObject{"mother", mother}
	h5 := lessgo.LoadFormObject{"motherPhone", motherPhone}
	h6 := lessgo.LoadFormObject{"homePhone", homePhone}
	h7 := lessgo.LoadFormObject{"child", child}
	h8 := lessgo.LoadFormObject{"employeeId", employeeId}
	h9 := lessgo.LoadFormObject{"centerName", centerName}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)
	loadFormObjects = append(loadFormObjects, h9)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}

func SendToTmkSaveAction(w http.ResponseWriter, r *http.Request) {
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

	id := r.FormValue("tcid")
	employeeId := r.FormValue("employeeId")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "update tmk_consumer set employee_id=? where id=? "

	lessgo.Log.Debug(sql)


	stmt, err := db.Prepare(sql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(employeeId,id)

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
