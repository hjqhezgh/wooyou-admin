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
	"net/http"
	"logic"
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

	flag,msg,err := logic.SaveConsumer(paramsMap)

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

