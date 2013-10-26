// Title：将客户分配给CallCenter
//
// Description:
//
// Author:black
//
// Createtime:2013-10-15 15:32
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-15 15:32 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"strconv"
	"time"
)

func SendToCallCenterAction(w http.ResponseWriter, r *http.Request) {

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

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql := "insert into tmk_consumer(consumer_id,center_id,cd_create_time) values(?,?,?)"

	stmt, err := db.Prepare(sql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	params := []interface{}{}
	params = append(params, id)
	params = append(params, _employee.CenterId)
	params = append(params, time.Now().Format("20060102150405"))

	res, err := stmt.Exec(params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = res.LastInsertId()

	m["success"] = "分配成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func AllSendToCallCenter(w http.ResponseWriter, r *http.Request) {
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

	db := lessgo.GetMySQL()
	defer db.Close()

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	params := []interface{}{}

	sql := "select c.id from consumer c  where 1=1 and c.center_id=? and c.id not in (select consumer_id from tmk_consumer) "
	insertSql := "insert into tmk_consumer(consumer_id,center_id,cd_create_time) values(?,?,?)"

	lessgo.Log.Debug(sql)

	params = append(params, _employee.CenterId)

	rows, err := db.Query(sql, params...)

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		var id string

		err := rows.Scan(&id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		stmt, err := tx.Prepare(insertSql)

		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(id, _employee.CenterId, time.Now().Format("20060102150405"))

		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx.Commit()

	m["success"] = "分配成功"
	commonlib.OutputJson(w, m, " ")

	return
}
