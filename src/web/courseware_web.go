// Title：课件相关服务
//
// Description:
//
// Author:black
//
// Createtime:2013-11-11 13:41
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-11 13:41 black 创建文档
package web

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"logic"
	"net/http"
	"strconv"
	"text/template"
)

func CoursewareListAction(w http.ResponseWriter, r *http.Request) {

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
		lessgo.Log.Error(err.Error())
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

	kw := r.FormValue("kw-like")

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	pageData, err := logic.CoursewareList(_employee.CenterId, kw, pageNo, pageSize)

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

func CoursewareSaveAction(w http.ResponseWriter, r *http.Request) {
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
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	intro := r.FormValue("intro")

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	flag, msg, err := logic.SaveCourseware(id, _employee.CenterId, name, intro)

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
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + msg
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m, " ")

	return
}

func CoursewareLoadAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")

	loadFormObjects, err := logic.LoadCourseware(id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}

func CoursewareUploadCallBack(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	state := r.FormValue("state")
	suffix := r.FormValue("suffix")
	fileName := r.FormValue("fileName")
	srcFileName := r.FormValue("srcFileName")
	fileSize := r.FormValue("fileSize")

	if state == "SUCCESS" {
		m["success"] = true
		m["suffix"] = suffix
		m["fileName"] = fileName
		m["srcFileName"] = srcFileName
		m["fileSize"] = fileSize
	} else {
		m["success"] = false
		m["msg"] = "文件上传发生错误，请重试"
	}

	commonlib.OutputJson(w, m, " ")
	return
}
