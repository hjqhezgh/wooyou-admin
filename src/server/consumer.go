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
		} else if roleId == "2" {
			dataType = "center"
			break
		} else if roleId == "" {
			dataType = "self"
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
	mother := r.FormValue("mother")
	motherPhone := r.FormValue("motherPhone")
	father := r.FormValue("father")
	fatherPhone := r.FormValue("fatherPhone")
	homePhone := r.FormValue("homePhone")
	child := r.FormValue("child")
	comeFromId := r.FormValue("come_from_id")
	centerId := r.FormValue("center_id")

	db := lessgo.GetMySQL()
	defer db.Close()

	if id == "" {
		sql := "insert into consumer(father,father_phone,mother,mother_phone,home_phone,employee_id,child,come_from_id,center_id) values(?,?,?,?,?,?,?,?,?)"

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

		_, err = stmt.Exec(father, fatherPhone, mother, motherPhone, homePhone, employee.UserId, child, comeFromId, centerId)

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
	} else {
		sql := "update consumer set father=?,father_phone=?,mother=?,mother_phone=?,home_phone=?,child=?,come_from_id=?,center_id=? where id=? "

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

		_, err = stmt.Exec(father, fatherPhone, mother, motherPhone, homePhone, child, comeFromId,centerId, id)

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

	id := r.FormValue("id")

	sql := "select father,father_phone,mother,mother_phone,home_phone,child,come_from_id,center_id from consumer where id=? "

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

	var father, fatherPhone, mother, motherPhone, homePhone, child, comeFrom,centerId string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &father, &fatherPhone, &mother, &motherPhone, &homePhone, &child, &comeFrom,&centerId)

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
	h6 := lessgo.LoadFormObject{"homePhone", homePhone}
	h7 := lessgo.LoadFormObject{"child", child}
	h8 := lessgo.LoadFormObject{"come_from_id", comeFrom}
	h9 := lessgo.LoadFormObject{"center_id", centerId}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)
	loadFormObjects = append(loadFormObjects, h9)

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
select ce.name,e.really_name,a.start_time,a.seconds,a.inout,a.localphone,a.filename,a.is_upload_finish
from audio a
left join consumer c on
(c.mother_phone=a.remotephone and c.mother_phone!='' and c.mother_phone is not null)
or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null)
left join employee e on
a.localphone=e.phone_in_center
left join center ce on
ce.cid=e.center_id
where c.id=1
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

	id := r.FormValue("id")

	params := []interface{}{}

	sql := " select a.aid,ce.name,e.really_name,a.start_time,a.seconds,a.inout,a.note,a.localphone,a.filename,a.is_upload_finish from audio a "
	sql += " left join consumer c on  "
	sql += " (c.mother_phone=a.remotephone and c.mother_phone!='' and c.mother_phone is not null) "
	sql += " or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null) "
	sql += " or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null) "
	sql += " left join employee e on "
	sql += " a.localphone=e.phone_in_center "
	sql += " left join center ce on "
	sql += " ce.cid=e.center_id "
	sql += " where c.id=? "

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

	sql += " order by a.aid desc limit ?,?"

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
