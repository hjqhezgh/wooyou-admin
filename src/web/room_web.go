// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-10 10:15
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-10 10:15 black 创建文档
package web

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"logic"
	"net/http"
	"strconv"
	"text/template"
)

func TimeSectionSave11Action(w http.ResponseWriter, r *http.Request) {
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

	startTime := r.FormValue("start_time")
	endTime := r.FormValue("end_time")

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	flag, msg, err := logic.InsertTimeSection(_employee.CenterId, startTime, endTime)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if !flag {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "操作失败:" + msg
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func RoomListAction(w http.ResponseWriter, r *http.Request) {
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

	pageNoString := r.FormValue("page")
	pageNo := 1
	if pageNoString != "" {
		pageNo, err = strconv.Atoi(pageNoString)
		if err != nil {
			pageNo = 1
			lessgo.Log.Warn("错误的pageNo:", pageNo)
		}
	}

	pageSizeString := r.FormValue("rows")
	pageSize := 10
	if pageSizeString != "" {
		pageSize, err = strconv.Atoi(pageSizeString)
		if err != nil {
			lessgo.Log.Warn("错误的pageSize:", pageSize)
		}
	}

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	pageData, err := logic.RoomPage(_employee.CenterId, pageNo, pageSize)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1

	commonlib.RenderTemplate(w, r, "page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/page.json")
}

func RoomSaveAction(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	channel := r.FormValue("channel")

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	flag, msg, err := logic.InsertRoom(_employee.CenterId, name, channel)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if !flag {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "操作失败:" + msg
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}
