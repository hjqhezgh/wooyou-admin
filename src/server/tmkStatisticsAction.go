// Title：TMK运营报表
//
// Description:
//
// Author:black
//
// Createtime:2013-10-23 16:24
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-23 16:24 black 创建文档
package server

import (
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"strings"
)

//tmk业绩报表
/*
select em.user_id,aa.num as '电话数',bb.num as '名单数',cc.num as '邀约数',dd.num as '签到数',ee.num as '缴费数'
from employee em left join
(select count(1) num,e.user_id from audio a left join employee e on a.cid=e.center_id and a.localphone=e.phone_in_center
where e.user_id is not null group by e.user_id ) aa on em.user_id= aa.user_id
left join
(select count(1) num,tmk_id from tmk_consumer group by tmk_id )bb on em.user_id=bb.tmk_id
left join
(select count(1) num,create_user from wyclass_free_child where create_time>=? and create_time <=? group by create_user) cc on em.user_id=cc.tmk_id
left join
(select count(1) num,tmk_id from(select tc.tmk_id,tc.consumer_id,wfsi.sign_in_time from tmk_consumer tc left join wyclass_free_sign_in wfsi on tc.consumer_id=wfsi.consumer_id
where wfsi.sign_in_time is not null ) a group by tmk_id )dd on em.user_id=dd.tmk_id
left join
(select count(1) num,tmk_id from(select tc.tmk_id,tc.consumer_id,pl.pay_time from tmk_consumer tc left join pay_log pl on tc.consumer_id=pl.consumer_id
where pl.pay_time is not null) b group by tmk_id) ee on em.user_id=ee.tmk_id
where cc.num is not null
order by em.user_id
*/
func TmkStatisticsAction(w http.ResponseWriter, r *http.Request) {

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

	centerId := r.FormValue("centerId-eq")
	startTime := r.FormValue("startTime-ge")
	endTime := r.FormValue("endTime-le")
	employeeId := r.FormValue("employee_id-eq")

	dataType := ""

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" {
			dataType = "all"
			break
		} else if roleCode == "cd"{
			dataType = "center"
			break
		}else {
			dataType = "self"
			break
		}
	}

	params := []interface{}{}

	sql := "select em.user_id,em.really_name,aa.num as '电话数',bb.num as '名单数',cc.num as '邀约数',dd.num as '签到数',ee.num as '缴费数' "
	sql += " from employee em left join "
	sql += " (select count(1) num,e.user_id from audio a left join employee e on a.cid=e.center_id and a.localphone=e.phone_in_center "
	sql += " where e.user_id is not null and start_time >=? and start_time<=?  group by e.user_id ) aa on em.user_id= aa.user_id "
	sql += " left join "
	sql += " (select count(1) num,tmk_id from tmk_consumer where tmk_create_time >=? and tmk_create_time<=? group by tmk_id )bb on em.user_id=bb.tmk_id "
	sql += " left join "
	sql += " (select count(1) num,create_user from wyclass_free_child where create_time>=? and create_time <=? group by create_user) cc on em.user_id=cc.create_user "
	sql += " left join "
	sql += " (select count(1) num,tmk_id from(select tc.tmk_id,tc.consumer_id from tmk_consumer tc left join wyclass_free_sign_in wfsi on tc.consumer_id=wfsi.consumer_id "
	sql += " where wfsi.sign_in_time is not null and wfsi.sign_in_time>=? and wfsi.sign_in_time<=? ) a group by tmk_id )dd on em.user_id=dd.tmk_id "
	sql += " left join "
	sql += " (select count(1) num,tmk_id from(select tc.tmk_id,tc.consumer_id from tmk_consumer tc left join pay_log pl on tc.consumer_id=pl.consumer_id "
	sql += " where pl.pay_time is not null and pl.pay_time>=? and pl.pay_time<=? ) b group by tmk_id) ee on em.user_id=ee.tmk_id "
	sql += " where em.phone_in_center is not null and em.phone_in_center != '' "

	defaultStartTime := "20000101000000"
	defaultEndTime := "29991231000000"

	if startTime != ""{
		defaultStartTime = startTime
	}

	if endTime != ""{
		defaultEndTime = endTime
	}

	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)

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
		sql += " and em.center_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		sql += " and em.center_id=? "
	}

	if dataType == "self" {
		params = append(params, employee.UserId)
		sql += " and em.user_id=? "
	}

	if dataType != "self" && employeeId!= "" {
		params = append(params, employeeId)
		sql += " and em.user_id=? "
	}


	countSql := "select count(1) from (" + sql + ") num"

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

	sql += " order by em.user_id limit ?,?"

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(sql)

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

		for i := 0; i < 6; i++ {
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

//tmk运营报表详情
/*
select c.id,e.really_name,ce.name,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone,c.child,c.contact_status
from tmk_consumer tc
left join consumer c on tc.consumer_id=c.id
left join employee e on e.user_id=tc.employee_id
left join center ce on ce.cid=c.center_id
where tc.employee_id=?
*/
func TmkStatisticsDetailAction(w http.ResponseWriter, r *http.Request) {
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

	sql := " select c.id,e.user_id,e.really_name,ce.name,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone,c.child,c.contact_status "
	sql += " from tmk_consumer tc "
	sql += " left join consumer c on tc.consumer_id=c.id "
	sql += " left join employee e on e.user_id=tc.employee_id "
	sql += " left join center ce on ce.cid=c.center_id "
	sql += " where tc.employee_id=? "

	countSql := "select count(1) from tmk_consumer where employee_id=?"

	params = append(params, id)

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

	sql += " limit ?,?"

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(sql)

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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		model.Id = r.Intn(1000)
		model.Props = []*lessgo.Prop{}

		fillObjects := []interface{}{}

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
