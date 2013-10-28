// Title：客户列表数据
//
// Description:
//
// Author:black
//
// Createtime:2013-09-26 15:50
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-09-26 15:50 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"time"
)

func AddAppAccountLoadAction(w http.ResponseWriter, r *http.Request) {

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

	sql := "select father,father_phone,mother,mother_phone,child,center_id from consumer where id=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	var father, fatherPhone, mother, motherPhone, child,centerId string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &father, &fatherPhone, &mother, &motherPhone, &child,&centerId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	m["success"] = true

	loadFormObjects := []lessgo.LoadFormObject{}

	h1 := lessgo.LoadFormObject{"id", id}
	h2 := lessgo.LoadFormObject{"father", father}
	h3 := lessgo.LoadFormObject{"fatherPhone", fatherPhone}
	h4 := lessgo.LoadFormObject{"mother", mother}
	h5 := lessgo.LoadFormObject{"motherPhone", motherPhone}
	h6 := lessgo.LoadFormObject{"child", child}
	h7 := lessgo.LoadFormObject{"childName", child}
	h8 := lessgo.LoadFormObject{"account", ""}
	h9 := lessgo.LoadFormObject{"userName", ""}
	h10 := lessgo.LoadFormObject{"centerId", centerId}

	if motherPhone!="" {
		h8 = lessgo.LoadFormObject{"account", motherPhone}
		h9 = lessgo.LoadFormObject{"userName", mother}
	}else if fatherPhone!= ""{
		h8 = lessgo.LoadFormObject{"account", fatherPhone}
		h9 = lessgo.LoadFormObject{"userName", father}
	}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)
	loadFormObjects = append(loadFormObjects, h9)
	loadFormObjects = append(loadFormObjects, h10)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")

}

func AddAppAccountSaveAction(w http.ResponseWriter, r *http.Request) {

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
	account := r.FormValue("account")
	pwd := r.FormValue("pwd")
	userName := r.FormValue("userName")
	childName := r.FormValue("childName")
	childBirthday := r.FormValue("childBirthday")
	childSex := r.FormValue("childSex")
	centerId := r.FormValue("centerId")

	db := lessgo.GetMySQL()
	defer db.Close()

	searchSql := "select count(1) from parent where telephone=?"

	lessgo.Log.Debug(searchSql)

	rows, err := db.Query(searchSql, account)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	totalNum := 0

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if totalNum>0{
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该手机号已经被注册，请更换一个注册手机号"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertParentSql := "insert into parent(name,password,telephone,reg_date) values(?,?,?,?)"

	lessgo.Log.Debug(insertParentSql)

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	stmt1, err := tx.Prepare(insertParentSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	res, err := stmt1.Exec(userName, pwd,account, time.Now().Format("20060102150405"))

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	parentId, err := res.LastInsertId()

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	updateConsumerSql := "update consumer set parent_id=? where id=?"

	lessgo.Log.Debug(updateConsumerSql)

	stmt2, err := tx.Prepare(updateConsumerSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt2.Exec(parentId,id)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertChildSql := "insert into child(name,pid,sex,birthday,center_id) values(?,?,?,?,?)"

	lessgo.Log.Debug(insertChildSql)

	stmt3, err := tx.Prepare(insertChildSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	childBirthdayDate,err :=  time.ParseInLocation("2006-01-02", childBirthday, time.Local)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "小孩子生日格式错误"
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt3.Exec(childName,parentId,childSex,childBirthdayDate.Format("20060102"),centerId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	commonlib.OutputJson(w, m, " ")
	return
}
