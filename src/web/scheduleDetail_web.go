// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-15 16:51
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-15 16:51 black 创建文档
package web

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"logic"
	"net/http"
	"strconv"
	"strings"
)

func AddChildToClassAction(w http.ResponseWriter, r *http.Request) {

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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	flag, msg, err := logic.AddChildToClass(classId, scheduleId, ids, employee.UserId)

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

func AddChildToClassQuickAction(w http.ResponseWriter, r *http.Request) {

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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	consumerId := r.FormValue("consumerId")

	flag, msg, err := logic.AddChildToClassQuick(classId, scheduleId, consumerId, employee.UserId)

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

func ContractCheckInSaveAction(w http.ResponseWriter, r *http.Request) {

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
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	childId := r.FormValue("childId")
	scheduleId := r.FormValue("scheduleId")
	contractId := r.FormValue("contractId")
	actionType := r.FormValue("type")

	flag, msg, err := logic.ContractCheckIn(childId, scheduleId, contractId, actionType)

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

func CreateWeekScheduleAction(w http.ResponseWriter, r *http.Request) {
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

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := logic.FindEmployeeById(userId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	err = r.ParseForm()

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	firstDayOfWeek := r.FormValue("firstDayOfWeek")

	flag, msg, err := logic.CreateWeekSchedule(_employee.CenterId, firstDayOfWeek, employee.UserId)

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

func ClassScheduleDetailLeaveAction(w http.ResponseWriter, r *http.Request) {
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

	scheduleId := r.FormValue("scheduleId")
	childId := r.FormValue("childId")

	flag, msg, err := logic.ClassScheduleDetailLeave(childId, scheduleId, employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func ClassScheduleDetailTruantAction(w http.ResponseWriter, r *http.Request) {
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

	scheduleId := r.FormValue("scheduleId")
	childId := r.FormValue("childId")

	flag, msg, err := logic.ClassScheduleDetailTruant(childId, scheduleId, employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func ClassScheduleDetailSignInAction(w http.ResponseWriter, r *http.Request) {
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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "tmk" {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "对不起，您没有权限进行签到"
			commonlib.OutputJson(w, m, " ")
			return
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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	flag, msg, err := logic.ClassScheduleDetailSignIn(ids, scheduleId,classId, employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func ChildSignInWithoutClassAction(w http.ResponseWriter, r *http.Request) {
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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "tmk" {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "对不起，您没有权限进行签到"
			commonlib.OutputJson(w, m, " ")
			return
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

	consumerIds := r.FormValue("ids")

	flag, msg, err := logic.ChildSignInWithoutClass(consumerIds,employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func AddChildForNormalTempelateAction(w http.ResponseWriter, r *http.Request) {
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

	scheduleId := r.FormValue("scheduleId")
	childs := r.FormValue("ids")

	flag, msg, err := logic.AddChildForNormalTempelate(childs,scheduleId,employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func AddChildForNormalOnceAction(w http.ResponseWriter, r *http.Request) {
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

	scheduleId := r.FormValue("scheduleId")
	childs := r.FormValue("ids")

	flag, msg, err := logic.AddChildForNormalOnce(childs,scheduleId,employee.UserId)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

	return
}

func RemoveChildFromScheduleAction(w http.ResponseWriter, r *http.Request) {
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

	roleCodes := strings.Split(employee.RoleCode, ",")

	dataType := ""

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "cd" {
			dataType = "all"
			break
		} else {
			dataType = "self"
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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	childId := r.FormValue("childId")

	flag, msg, err := logic.RemoveChildFromSchedule(childId,scheduleId,classId,dataType,employee.UserId)

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

func ChangeClassScheduleAction(w http.ResponseWriter, r *http.Request) {

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

	oldScheduleId := r.FormValue("oldScheduleId")
	newScheduleId := r.FormValue("newScheduleId")
	childId := r.FormValue("childId")

	flag, msg, err := logic.ChangeClassSchedule(childId, newScheduleId, oldScheduleId, employee.UserId)

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
		m["msg"] = "调班失败:" + msg
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m, " ")

	return
}
