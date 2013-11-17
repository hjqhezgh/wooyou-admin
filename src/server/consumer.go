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

	roleIds := strings.Split(employee.RoleId, ",")

	for _, roleId := range roleIds {
		if roleId == "1" || roleId == "3" || roleId == "6" || roleId == "10" || roleId == "11" {
			dataType = "all"
			break
		} else {//if roleId == "2" {
			dataType = "center"
			break
		} /*else if roleId == "" {
			dataType = "self"
		}*/
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

	mySort := r.FormValue("mySort-eq")
	name := r.FormValue("name-like")
	centerId := r.FormValue("cid-eq")

	params := []interface{}{}

	sql := "select c.id,ce.name,e.really_name,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone,c.child,a.num,a.maxtime,c.contact_status,c.parent_id from consumer c left join (select count(1) num,max(start_time) maxtime, remotephone from audio group by remotephone) a on (c.mother_phone=a.remotephone and c.mother_phone!='' and c.mother_phone is not null) or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null) left join employee e on e.user_id=c.employee_id left join center ce on ce.cid=c.center_id where 1=1 "

	if name != "" {
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		sql += " and c.mother like ? or c.father like ? or c.child like ? "
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
		sql += " and ce.cid=? "
	}

	if dataType == "self" {
		params = append(params, employee.UserId)
		sql += " and e.user_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		sql += " and c.center_id=? "
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

	if mySort == "" || mySort == "time" {
		sql += " order by a.maxtime desc,c.id desc  limit ?,?"
	} else if mySort == "frequency" {
		sql += " order by a.num desc,c.id desc  limit ?,?"
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

		for i := 0; i < 12; i++ {
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

	id := r.FormValue("id")
	status := r.FormValue("status")

	sql := "select contact_status from consumer where id=?"

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

	sql = "update consumer set contact_status=? where id=? "

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

	m["success"] = "客户状态更改成功"
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
