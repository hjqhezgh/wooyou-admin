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
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	CONSUMER_STATUS_NO_CONTACT = "1"
	CONSUMER_STATUS_WAIT = "2"
	CONSUMER_STATUS_ABANDON = "3"
	CONSUMER_STATUS_NO_SIGNIN = "4"
	CONSUMER_STATUS_SIGNIN = "5"
)

//客户分页数据服务
/*
select cons.id,ce.name as centerName,e.really_name,cont.name,cont.phone,cons.home_phone,cons.child,cons.contact_status,cons.parent_id
from
(select consumer_id,min(id) contacts_id from contacts group by consumer_id)a
left join consumer_new cons on cons.id=a.consumer_id
left join contacts cont on cont.id=a.contacts_id
left join center ce on ce.cid=cons.center_id
left join employee e on e.user_id=cons.current_tmk_id
*/
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
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" ||  roleCode == "tmk" {
			dataType = "all"
			break
		} else{
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

	centerId1 := r.FormValue("centerId-eq")
	centerId2 := r.FormValue("centerId")
	tmkId1 := r.FormValue("tmkId-eq")
	tmkId2 := r.FormValue("tmkId")
	status := r.FormValue("status-eq")
	lastContractStartTime := r.FormValue("lastContractStartTime-ge")
	lastContractEndTime := r.FormValue("lastContractEndTime-le")
	kw := r.FormValue("kw-like")
	sort := r.FormValue("sort-eq")
	payStatus := r.FormValue("payStatus-eq")
	timeType := r.FormValue("timeType-eq")

	params := []interface{}{}

	sql := " select cons.id,ce.name as centerName,e.really_name,cont.name,cont.phone,cons.home_phone,cons.child,cons.contact_status,cons.parent_id,cons.remark,cons.pay_status,cons.pay_time "
	sql += " from "
	sql += " (select c.consumer_id,min(c.id) contacts_id from contacts c  "
	if kw != "" {
		sql += "left join consumer_new b on b.id=c.consumer_id where c.phone like ? or c.name like ? or b.child like ? or b.remark like ? or b.home_phone like ? "
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}
	sql += " group by c.consumer_id)a "
	sql += " left join consumer_new cons on cons.id=a.consumer_id "
	sql += " left join contacts cont on cont.id=a.contacts_id "
	sql += " left join center ce on ce.cid=cons.center_id "
	sql += " left join employee e on e.user_id=cons.current_tmk_id where 1=1 "

	if status != "" {
		params = append(params, status)
		sql += " and cons.contact_status=? "
	}

	if lastContractStartTime!= "" && timeType=="1"  {
		params = append(params,lastContractStartTime)
		sql += " and cons.sign_in_time>=? "
	}

	if lastContractStartTime!= "" && timeType=="2"  {
		params = append(params,lastContractStartTime)
		sql += " and cons.pay_time>=? "
	}

	if lastContractEndTime!= "" && timeType=="1"{
		params = append(params,lastContractEndTime)
		sql += " and cons.sign_in_time<=? "
	}

	if lastContractEndTime!= "" && timeType=="2"  {
		params = append(params,lastContractEndTime)
		sql += " and cons.pay_time<=? "
	}

	if payStatus == "1" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=1 "
	}else if payStatus == "2" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=2 "
	}else if payStatus == "3" {
		sql += " and (cons.pay_time is null or cons.pay_time ='')"
	}

	if tmkId1 != ""{
		params = append(params,tmkId1)
		sql += " and cons.current_tmk_id=? "
	}

	if tmkId2 != ""{
		params = append(params,tmkId2)
		sql += " and cons.current_tmk_id=? "
	}

	if dataType == "center" {
		userId, _ := strconv.Atoi(employee.UserId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
		params = append(params, _employee.CenterId)
		sql += " and cons.center_id=? "
	}

	if centerId1 != "" && dataType == "all" {
		params = append(params, centerId1)
		sql += " and cons.center_id=? "
	}

	if centerId2 != "" && dataType == "all" {
		params = append(params, centerId2)
		sql += " and cons.center_id=? "
	}

	countSql := ""

	countSql = "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

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

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	currPageNo := pageNo

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	if sort == "" || sort == "create_time" {
		sql += " order by cons.id desc  limit ?,? "
	} else if sort == "last_time" {
		sql += " order by cons.last_contact_time desc limit ?,? "
	}

	lessgo.Log.Debug(sql)

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(sql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < 11; i++ {
			prop := new(lessgo.Prop)
			prop.Name = fmt.Sprint(i)
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}

	pageData := commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects)

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*lessgo.Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/entity_page.json")

}

//客户保存服务
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

	id := r.FormValue("id")
	phone := r.FormValue("phone")
	contactsName := r.FormValue("contactsName")
	homePhone := r.FormValue("homePhone")
	child := r.FormValue("child")
	year := r.FormValue("year")
	month := r.FormValue("month")
	birthday := r.FormValue("birthday")
	come_from_id := r.FormValue("come_from_id")
	center_id := r.FormValue("center_id")
	remark := r.FormValue("remark")

	db := lessgo.GetMySQL()
	defer db.Close()

	if id == "" {

		homePhoneFlag,err := CheckConsumerPhoneExist(homePhone)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		phoneFlag,err := CheckConsumerPhoneExist(phone)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		if phoneFlag{
			m["success"] = false
			m["code"] = 100
			m["msg"] = "联系人电话已经存在，无需重复录入"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if homePhoneFlag{
			m["success"] = false
			m["code"] = 100
			m["msg"] = "家庭电话已经在系统中存在，无需重复录入"
			commonlib.OutputJson(w, m, " ")
			return
		}

		/**********************获取最后联系状态**********************/
		contactStatusSql := " select seconds,start_time,end_time,`inout` from audio where remotephone=? "
		contactStatusParams := []interface {}{}
		contactStatusParams = append(contactStatusParams,phone)

		if homePhone!=""{
			contactStatusSql += " or remotephone=? "
			contactStatusParams = append(contactStatusParams,homePhone)
		}
		contactStatusSql += " order by start_time desc,aid desc "
		lessgo.Log.Debug(contactStatusSql)

		contactStatusRows, err := db.Query(contactStatusSql, contactStatusParams...)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		contactStatus := CONSUMER_STATUS_NO_CONTACT
		lastContactTime := ""

		for contactStatusRows.Next() {
			seconds := ""
			endtime := ""
			starttime := ""
			inout := ""

			err = commonlib.PutRecord(contactStatusRows, &seconds,&starttime,&endtime,&inout)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			if seconds=="" && endtime==""{//说明是在接电话的过程中添加用户一定是有效数据，直接取当前的信息最为最后联系信息
				contactStatus = CONSUMER_STATUS_WAIT
				lastContactTime = strings.Replace(starttime,"-","",-1)
				lastContactTime = strings.Replace(lastContactTime," ","",-1)
				lastContactTime = strings.Replace(lastContactTime,":","",-1)
				break
			}

			//说明这个是有效电话，直接取这次的通话信息
			if seconds!="00:00:00" && inout!="未接听"{
				contactStatus = CONSUMER_STATUS_WAIT
				lastContactTime = strings.Replace(starttime,"-","",-1)
				lastContactTime = strings.Replace(lastContactTime," ","",-1)
				lastContactTime = strings.Replace(lastContactTime,":","",-1)
				break
			}
		}

		tx, err := db.Begin()
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		/**********************数据插入consumer表**********************/
		insertConsumerSql := "insert into consumer_new(center_id,contact_status,home_phone,create_time,child,year,month,birthday,last_contact_time,come_from_id,remark,create_user_id) values(?,?,?,?,?,?,?,?,?,?,?,?)"
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

		insertConsumerRes, err := insertConsumerStmt.Exec(center_id,contactStatus,homePhone,time.Now().Format("20060102150405"),child,year,month,birthday,lastContactTime,come_from_id,remark,employee.UserId)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		consumerId,err := insertConsumerRes.LastInsertId()
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

		_, err = insertContactsStmt.Exec(contactsName,phone,"1",consumerId)
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
	} else {

		sql := "update consumer_new set remark=?,home_phone=?,child=?,year=?,month=?,birthday=?,come_from_id=? where id=? "

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

		_, err = stmt.Exec(remark, homePhone, child, year, month, birthday, come_from_id,id)

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

}

//客户读取服务
func ConsumerLoadAction(w http.ResponseWriter, r *http.Request) {

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

	id := r.FormValue("consumerId")
	aid := r.FormValue("aid")

	loadFormObjects := []lessgo.LoadFormObject{}

	if id!= ""{
		sql := "select center_id,home_phone,child,year,month,birthday,come_from_id,remark from consumer_new where id=? "

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

		var center_id, home_phone, child, year, month, birthday, come_from_id,remark string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &center_id, &home_phone, &child, &year, &month, &birthday, &come_from_id,&remark)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		m["success"] = true

		h1 := lessgo.LoadFormObject{"homePhone", home_phone}
		h2 := lessgo.LoadFormObject{"child", child}
		h3 := lessgo.LoadFormObject{"year", year}
		h4 := lessgo.LoadFormObject{"month", month}
		h5 := lessgo.LoadFormObject{"birthday", birthday}
		h6 := lessgo.LoadFormObject{"come_from_id", come_from_id}
		h7 := lessgo.LoadFormObject{"center_id", center_id}
		h8 := lessgo.LoadFormObject{"remark", remark}
		h9 := lessgo.LoadFormObject{"id", id}

		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
		loadFormObjects = append(loadFormObjects, h3)
		loadFormObjects = append(loadFormObjects, h4)
		loadFormObjects = append(loadFormObjects, h5)
		loadFormObjects = append(loadFormObjects, h6)
		loadFormObjects = append(loadFormObjects, h7)
		loadFormObjects = append(loadFormObjects, h8)
		loadFormObjects = append(loadFormObjects, h9)
	}else if aid!="" {
		sql := "select remotephone,cid from audio where aid=? "
		lessgo.Log.Debug(sql)

		db := lessgo.GetMySQL()
		defer db.Close()

		rows, err := db.Query(sql, aid)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		var remotephone,cid string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &remotephone, &cid)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}


		m["success"] = true
		h1 := lessgo.LoadFormObject{"phone", remotephone}
		h2 := lessgo.LoadFormObject{"center_id", cid}
		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
	}else{
		userId, _ := strconv.Atoi(employee.UserId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		h1 := lessgo.LoadFormObject{"center_id", _employee.CenterId}
		loadFormObjects = append(loadFormObjects, h1)
	}


	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}

func ConsumerStatusChangeAction(w http.ResponseWriter, r *http.Request) {
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

	id := r.FormValue("ids")
	status := r.FormValue("status")

	if strings.Contains(id,","){
		m["success"] = false
		m["code"] = 100
		m["msg"] = "只能操作一个客户"
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql := "select contact_status from consumer_new where id=?"

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	var oldStatus string

	if rows.Next() {
		err := commonlib.PutRecord(rows, &oldStatus)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql = "update consumer_new set contact_status=? where id=? "

	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(status, id)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql = "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"

	lessgo.Log.Debug(sql)

	stmt, err = tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(id, employee.UserId, time.Now().Format("20060102150405"), oldStatus, status)

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
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")

}

//查看通话记录
/*
select c.aid,ce.name as centerName,e.really_name,c.start_time,c.seconds,c.inout,c.localphone,c.filename,c.is_upload_finish,c.contractName,c.note from(
select  a.aid,a.start_time,a.seconds,a.inout,a.localphone,a.filename,a.is_upload_finish,contract_phone.name contractName,a.note,a.cid from audio a
inner join (select phone,name from contacts where consumer_id=1) contract_phone on
contract_phone.phone=a.remotephone
union
select  a.aid,a.start_time,a.seconds,a.inout,a.localphone,a.filename,a.is_upload_finish,'家庭电话' contractName,a.note,a.cid from audio a
inner join (select home_phone,child from consumer_new where id=1 and home_phone is not null and home_phone !='') b on
b.home_phone=a.remotephone
)c
left join employee e on e.center_id=c.cid and e.phone_in_center=c.localphone
left join center ce on ce.cid=e.center_id
*/
func ConsumerContactRecordListAction(w http.ResponseWriter, r *http.Request) {

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

	id := r.FormValue("consumerId")

	params := []interface{}{}

	sql := " select c.aid,ce.name as centerName,e.really_name,c.start_time,c.seconds,c.inout,c.localphone,c.filename,c.is_upload_finish,c.contractName,c.note,c.cid from( "
	sql += " select  a.aid,a.start_time,a.seconds,a.inout,a.localphone,a.filename,a.is_upload_finish,contract_phone.name contractName,a.note,a.cid from audio a "
	sql += " inner join (select phone,name from contacts where consumer_id=?) contract_phone on "
	sql += " contract_phone.phone=a.remotephone "
	sql += " union "
	sql += " select  a.aid,a.start_time,a.seconds,a.inout,a.localphone,a.filename,a.is_upload_finish,'家庭电话' contractName,a.note,a.cid from audio a "
	sql += " inner join (select home_phone,child from consumer_new where id=? and home_phone is not null and home_phone !='') b on "
	sql += " b.home_phone=a.remotephone "
	sql += " )c "
	sql += " left join employee e on e.center_id=c.cid and e.phone_in_center=c.localphone "
	sql += " left join center ce on ce.cid=e.center_id "

	params = append(params, id)
	params = append(params, id)

	countSql := ""

	countSql = "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

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

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	currPageNo := pageNo

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	sql += " order by c.start_time desc ,c.aid desc limit ?,?"

	lessgo.Log.Debug(sql)

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(sql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < 11; i++ {
			prop := new(lessgo.Prop)
			prop.Name = fmt.Sprint(i)
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}

	pageData := commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects)

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*lessgo.Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/entity_page.json")

}

//判断客户是否已经存在
func CheckConsumerPhoneExist(phone string) (bool,error) {

	if phone == ""{
		return false, nil
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select count(1) from consumer_new where home_phone=? ";

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false,err
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false,err
		}
	}

	if  num > 0 {
		return true ,nil
	}


	sql = "select count(1) from contacts where phone=? ";

	lessgo.Log.Debug(sql)

	rows, err = db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false,err
	}

	num = 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false,err
		}
	}

	if  num > 0 {
		return true ,nil
	}

	return false,nil
}

func TmkAllConsumerListAction(w http.ResponseWriter, r *http.Request) {

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

	dataType := ""

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "tmk" ||  roleCode == "yyzj"{
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
		}
	}

	centerId1 := r.FormValue("centerId-eq")
	status := r.FormValue("status-eq")
	lastContractStartTime := r.FormValue("lastContractStartTime-ge")
	lastContractEndTime := r.FormValue("lastContractEndTime-le")
	kw := r.FormValue("kw-like")

	params := []interface{}{}

	sql := " select cons.id,ce.name as centerName,e.really_name,cont.name,cont.phone,cons.home_phone,cons.child,cons.contact_status,cons.parent_id,cons.remark "
	sql += " from "
	sql += " (select c.consumer_id,min(c.id) contacts_id from contacts c  "
	if kw != "" {
		sql += "left join consumer_new b on b.id=c.consumer_id where c.phone like ? or c.name like ? or b.child like ? or b.remark like ? or b.home_phone like ? "
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}
	sql += " group by c.consumer_id)a "
	sql += " left join consumer_new cons on cons.id=a.consumer_id "
	sql += " left join contacts cont on cont.id=a.contacts_id "
	sql += " left join center ce on ce.cid=cons.center_id "
	sql += " left join employee e on e.user_id=cons.last_tmk_id where 1=1 and cons.is_own_by_tmk=2 "

	if dataType=="center" {
		userId, _ := strconv.Atoi(employee.UserId)
		_employee, err := FindEmployeeById(userId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		params = append(params, _employee.CenterId)
		sql += " and cons.center_id=? "
	}

	if status != "" {
		params = append(params, status)
		sql += " and cons.contact_status=? "
	}

	if lastContractStartTime!= "" {
		params = append(params,lastContractStartTime)
		sql += " and cons.last_contact_time>=? "
	}

	if lastContractEndTime!= "" {
		params = append(params,lastContractEndTime)
		sql += " and cons.last_contact_time<=? "
	}

	if centerId1 != "" && dataType=="all" {
		params = append(params, centerId1)
		sql += " and cons.center_id=? "
	}

	countSql := ""

	countSql = "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

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

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	currPageNo := pageNo

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	sql += " order by cons.id desc  limit ?,? "

	lessgo.Log.Debug(sql)

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(sql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < 9; i++ {
			prop := new(lessgo.Prop)
			prop.Name = fmt.Sprint(i)
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}

	pageData := commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects)

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*lessgo.Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/entity_page.json")

}


func TmkInviteAction(w http.ResponseWriter, r *http.Request) {
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

	ids := r.FormValue("ids")

	if strings.Contains(ids,",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "每次只能邀约一个客户"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select count(1) from tmk_consumer where consumer_id=? "
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, ids)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	totalNum := 0

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if totalNum > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已被其他TMK邀约，请选择其他客户"
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql2 := "select count(1) from tmk_consumer tc left join consumer_new c on tc.consumer_id=c.id where c.contact_status=1 and tc.tmk_id=?  "
	lessgo.Log.Debug(sql2)

	rows2, err := db.Query(sql2,employee.UserId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if rows2.Next() {
		err = rows2.Scan(&totalNum)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	//这里将来要改成可配置的
	if totalNum > 4 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "您还有5个未联系的客户，请处理结束后再邀约新客户"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertSql := "insert into tmk_consumer(tmk_id,consumer_id,tmk_create_time) values (?,?,?)"
	lessgo.Log.Debug(insertSql)

	stmt, err := tx.Prepare(insertSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(employee.UserId,ids,time.Now().Format("20060102150405"))

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	updateSql := "update consumer_new set is_own_by_tmk=1,current_tmk_id=?,last_tmk_id=? where id=?"

	stmt1, err := tx.Prepare(updateSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt1.Exec(employee.UserId,employee.UserId,ids)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		tx.Rollback()
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m, " ")
	return
}

/*
select cons.*,cont.name,cont.phone
from
tmk_consumer tc
left join (select c.consumer_id,min(c.id) contacts_id from contacts c group by c.consumer_id) a on a.consumer_id=tc.consumer_id
left join contacts cont on cont.id=a.contacts_id
left join consumer_new cons on cons.id=a.consumer_id
*/
func TmkConsumerSelfListAction(w http.ResponseWriter, r *http.Request) {

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

	centerId1 := r.FormValue("centerId-eq")
	status := r.FormValue("status-eq")
	lastContractStartTime := r.FormValue("lastContractStartTime-ge")
	lastContractEndTime := r.FormValue("lastContractEndTime-le")
	kw := r.FormValue("kw-like")
	payStatus := r.FormValue("payStatus-eq")
	timeType := r.FormValue("timeType-eq")

	params := []interface{}{}

	sql := " select cons.id,ce.name as centerName,cont.name,cont.phone,cons.home_phone,cons.child,cons.contact_status,cons.parent_id,cons.remark,cons.center_id,cons.pay_status,cons.pay_time "
	sql += " from tmk_consumer tc"
	sql += " inner join (select c.consumer_id,min(c.id) contacts_id from contacts c  "
	if kw != "" {
		sql += "left join consumer_new b on b.id=c.consumer_id where c.phone like ? or c.name like ? or b.child like ? or b.remark like ? or b.home_phone like ? "
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
		params = append(params, "%"+kw+"%")
	}
	sql += " group by c.consumer_id)a on a.consumer_id=tc.consumer_id "
	sql += " left join consumer_new cons on cons.id=a.consumer_id "
	sql += " left join contacts cont on cont.id=a.contacts_id "
	sql += " left join center ce on ce.cid=cons.center_id "
	sql += " where tc.tmk_id= "+employee.UserId

	if status != "" {
		params = append(params, status)
		sql += " and cons.contact_status=? "
	}

	if lastContractStartTime!= "" && timeType=="1"  {
		params = append(params,lastContractStartTime)
		sql += " and cons.sign_in_time>=? "
	}

	if lastContractStartTime!= "" && timeType=="2"  {
		params = append(params,lastContractStartTime)
		sql += " and cons.pay_time>=? "
	}

	if lastContractEndTime!= "" && timeType=="1"{
		params = append(params,lastContractEndTime)
		sql += " and cons.sign_in_time<=? "
	}

	if lastContractEndTime!= "" && timeType=="2"  {
		params = append(params,lastContractEndTime)
		sql += " and cons.pay_time<=? "
	}

	if payStatus == "1" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=1 "
	}else if payStatus == "2" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=2 "
	}else if payStatus == "3" {
		sql += " and (cons.pay_time is null or cons.pay_time ='')"
	}

	if centerId1 != "" {
		params = append(params, centerId1)
		sql += " and cons.center_id=? "
	}

	countSql := ""

	countSql = "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

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

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	currPageNo := pageNo

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	sql += " order by cons.contact_status ,cons.last_contact_time desc limit ?,? "

	lessgo.Log.Debug(sql)

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(sql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < 11; i++ {
			prop := new(lessgo.Prop)
			prop.Name = fmt.Sprint(i)
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}

	pageData := commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects)

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*lessgo.Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/entity_page.json")

}

//客户缴费
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

	ids := r.FormValue("ids")
	status := r.FormValue("status")

	if strings.Contains(ids,",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "每次只能缴费一个客户"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	checkSql := "select count(1) from pay_log where consumer_id=? "
	rows, err := db.Query(checkSql, ids)

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

	if totalNum > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该用户已缴费，无需重复操作"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertSql := "insert into pay_log(consumer_id,pay_time,employee_id) values(?,?,?)"
	lessgo.Log.Debug(insertSql)

	stmt, err := tx.Prepare(insertSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(ids,time.Now().Format("20060102150405"),employee.UserId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	updateSql := "update consumer_new set pay_time=?,pay_status=? where id=?"
	lessgo.Log.Debug(updateSql)

	stmt, err = tx.Prepare(updateSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(time.Now().Format("20060102150405"),status,ids)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "缴费成功"
	commonlib.OutputJson(w, m, " ")
	return
}
