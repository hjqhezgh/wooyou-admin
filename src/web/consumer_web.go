// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 16:29
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 16:29 black 创建文档
package web

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"logic"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func ConsumerSaveAction(w http.ResponseWriter, r *http.Request) {

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

	paramsMap := make(map[string]string)

	paramsMap["id"] = r.FormValue("id")
	paramsMap["phone"] = r.FormValue("phone")
	paramsMap["contactsName"] = r.FormValue("contactsName")
	paramsMap["homePhone"] = r.FormValue("homePhone")
	paramsMap["child"] = r.FormValue("child")
	paramsMap["year"] = r.FormValue("year")
	paramsMap["month"] = r.FormValue("month")
	paramsMap["birthday"] = r.FormValue("birthday")
	paramsMap["comeFromId"] = r.FormValue("come_from_id")
	paramsMap["centerId"] = r.FormValue("center_id")
	paramsMap["remark"] = r.FormValue("remark")
	paramsMap["parttimeName"] = r.FormValue("parttimeName")
	paramsMap["level"] = r.FormValue("level")
	paramsMap["createUser"] = employee.UserId

	flag, msg, err := logic.SaveConsumer(paramsMap)

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
		m["msg"] = "保存失败:" + msg
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m, " ")

	return
}

func ConsumerListAction(w http.ResponseWriter, r *http.Request) {

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

	dataType := ""

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "tmk" {
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
		}
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

	paramsMap := make(map[string]string)

	paramsMap["centerId1"] = r.FormValue("centerId-eq")
	paramsMap["centerId2"] = r.FormValue("centerId")
	paramsMap["tmkId1"] = r.FormValue("tmkId-eq")
	paramsMap["tmkId2"] = r.FormValue("tmkId")
	paramsMap["status"] = r.FormValue("status-eq")
	paramsMap["lastContractStartTime"] = r.FormValue("lastContractStartTime-ge")
	paramsMap["lastContractEndTime"] = r.FormValue("lastContractEndTime-le")
	paramsMap["kw"] = r.FormValue("kw-like")
	paramsMap["sort"] = r.FormValue("sort-eq")
	paramsMap["payStatus"] = r.FormValue("payStatus-eq")
	paramsMap["timeType"] = r.FormValue("timeType-eq")
	paramsMap["parttimeName"] = r.FormValue("partTimeName-eq")

	pageData, err := logic.ConsumerPage(paramsMap, dataType, employee.UserId, pageNo, pageSize)

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

func ConsumerPayAction(w http.ResponseWriter, r *http.Request) {
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

	consumerIds := r.FormValue("ids")
	payType := r.FormValue("status")

	flag, msg, err := logic.ConsumerPay(consumerIds,payType,employee.UserId)

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
