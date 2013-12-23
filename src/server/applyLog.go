// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-12-23 15:17
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-12-23 15:17 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"time"
	"strings"
)

func ApplyLogAddToConsumerAction(w http.ResponseWriter, r *http.Request) {

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

	year := r.FormValue("year")
	month := r.FormValue("month")
	childName := r.FormValue("childName")
	phone := r.FormValue("phone")
	comeFrom := r.FormValue("comeFrom")
	centerId := r.FormValue("centerId")
	createTime := r.FormValue("createTime")
	applyType := r.FormValue("type")
	contactsName := childName+"家长"

	remark := ""

	if comeFrom=="网站"{
		comeFrom = "5"
		remark += "通过网站报名"
	}else if comeFrom=="微信"{
		comeFrom = "6"
		remark += "通过微信报名"
	}else{
		comeFrom = "8"
	}

	if applyType == "2"{
		remark += "美术课程"
	}else if applyType=="1"{
		remark += "英语课程"
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	homePhoneFlag, err := CheckConsumerPhoneExist(phone)
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	phoneFlag, err := CheckConsumerPhoneExist(phone)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if phoneFlag {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "联系人电话已经存在，无需重复录入"
		commonlib.OutputJson(w, m, " ")
		return
	}

	if homePhoneFlag {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "家庭电话已经在系统中存在，无需重复录入"
		commonlib.OutputJson(w, m, " ")
		return
	}

	contactStatus := CONSUMER_STATUS_NO_CONTACT
	lastContactTime := ""

	if remark != "" {
		contactStatus = CONSUMER_STATUS_WAIT
		lastContactTime := strings.Replace(createTime," ","",-1)
		lastContactTime = strings.Replace(lastContactTime,":","",-1)
		lastContactTime = strings.Replace(lastContactTime,"-","",-1)
	}

	/**********************数据插入consumer表**********************/
	insertConsumerSql := "insert into consumer_new(center_id,contact_status,home_phone,create_time,child,year,month,birthday,last_contact_time,come_from_id,remark,create_user_id,parttime_name) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	lessgo.Log.Debug(insertConsumerSql)
	insertConsumerStmt, err := tx.Prepare(insertConsumerSql)
	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertConsumerRes, err := insertConsumerStmt.Exec(centerId, contactStatus, "", time.Now().Format("20060102150405"), childName, year, month, "", lastContactTime, comeFrom, remark, employee.UserId, "")
	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	consumerId, err := insertConsumerRes.LastInsertId()
	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	/***数据插入联系人表**/
	insertContactsSql := "insert into contacts(name,phone,is_default,consumer_id) values(?,?,?,?)"
	lessgo.Log.Debug(insertContactsSql)

	insertContactsStmt, err := tx.Prepare(insertContactsSql)
	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = insertContactsStmt.Exec(contactsName, phone, "1", consumerId)
	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	//如果添加阶段就有备注，那就要开始添加接触记录
	if remark != "" {
		insertContactLogSql := "insert into consumer_contact_log(create_user,create_time,note,consumer_id,type) values(?,?,?,?,?) "

		insertContactLogStmt, err := tx.Prepare(insertContactLogSql)
		if err != nil {
			tx.Rollback()

			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = insertContactLogStmt.Exec(employee.UserId, time.Now().Format("20060102150405"), remark, consumerId, 3)
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

	/******************插入parent表*****************/
	checkParentSql := "select pid from parent where telephone=? "
	lessgo.Log.Debug(checkParentSql)

	rows, err := db.Query(checkParentSql, phone)
	pid := 0

	if rows.Next() {
		err := commonlib.PutRecord(rows, &pid)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if pid != 0 {
		updateConsumerSql := "update consumer_new set parent_id=? where id=? "
		lessgo.Log.Debug(updateConsumerSql)
		stmt, err := tx.Prepare(updateConsumerSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(pid, consumerId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	} else {
		insertParentSql := "insert into parent(name,password,telephone,reg_date,come_form) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertParentSql)
		stmt, err := tx.Prepare(insertParentSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		newParentName := childName + "家长"
		if contactsName != "" {
			newParentName = contactsName
		}

		res, err := stmt.Exec(newParentName, "123456", phone, time.Now().Format("20060102150405"), 2)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		newParentId, err := res.LastInsertId()
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertChildSql := "insert into child(name,pid,sex,birthday,center_id) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertChildSql)

		stmt, err = tx.Prepare(insertChildSql)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		childBirthday := "20090101"
		if year != "" && month != "" {
			childBirthday = year + month + "01"
		} else if year != "" {
			childBirthday = year + "0101"
		}

		_, err = stmt.Exec(childName, newParentId, 1, childBirthday, centerId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		updateConsumerSql := "update consumer_new set parent_id=? where id=? "
		lessgo.Log.Debug(updateConsumerSql)

		stmt, err = tx.Prepare(updateConsumerSql)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(newParentId, consumerId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx.Commit()

	m["success"] = true
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
}
