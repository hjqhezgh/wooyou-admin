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
	CONSUMER_STATUS_WAIT       = "2"
	CONSUMER_STATUS_ABANDON    = "3"
	CONSUMER_STATUS_NO_SIGNIN  = "4"
	CONSUMER_STATUS_SIGNIN     = "5"
)

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

	if id != "" {
		sql := "select center_id,home_phone,child,year,month,birthday,come_from_id,remark,parttime_name,level from consumer_new where id=? "

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

		var center_id, home_phone, child, year, month, birthday, come_from_id, remark, parttimeName, level string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &center_id, &home_phone, &child, &year, &month, &birthday, &come_from_id, &remark, &parttimeName, &level)

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
		h10 := lessgo.LoadFormObject{"parttimeName", parttimeName}
		h11 := lessgo.LoadFormObject{"level", level}

		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
		loadFormObjects = append(loadFormObjects, h3)
		loadFormObjects = append(loadFormObjects, h4)
		loadFormObjects = append(loadFormObjects, h5)
		loadFormObjects = append(loadFormObjects, h6)
		loadFormObjects = append(loadFormObjects, h7)
		loadFormObjects = append(loadFormObjects, h8)
		loadFormObjects = append(loadFormObjects, h9)
		loadFormObjects = append(loadFormObjects, h10)
		loadFormObjects = append(loadFormObjects, h11)
	} else if aid != "" {
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

		var remotephone, cid string

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
	} else {
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

	if strings.Contains(id, ",") {
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
func CheckConsumerPhoneExist(phone string) (bool, error) {

	if phone == "" {
		return false, nil
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select count(1) from consumer_new where home_phone=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	sql = "select count(1) from contacts where phone=? "

	lessgo.Log.Debug(sql)

	rows, err = db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num = 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	return false, nil
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
		if roleCode == "tmk" || roleCode == "yyzj" {
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
		}
	}

	centerId1 := r.FormValue("centerId-eq")
	status := r.FormValue("status-eq")
	//	lastContractStartTime := r.FormValue("lastContractStartTime-ge")
	//	lastContractEndTime := r.FormValue("lastContractEndTime-le")
	kw := r.FormValue("kw-like")
	partTimeName := r.FormValue("partTimeName-eq")
	level := r.FormValue("level-eq")

	params := []interface{}{}
	paramsForCount := []interface{}{}

	sql := `select cons.id,ce.name as centerName,e.really_name,cons.level,cont.name,cont.phone,cons.child,cons.birthday,cons.year,cons.contact_status,cons.parent_id,d.remark,cf.name comeFromName,cons.parttime_name
			from
			(select a.consumer_id,min(a.id) contacts_id from contacts a left join consumer_new b on a.consumer_id=b.id
			where (a.name like ? or b.child like ? or a.phone like ? or b.home_phone like ?) and b.is_own_by_tmk=2 and b.pay_time is null and b.contact_status!=5 `

	params = append(params, "%"+kw+"%")
	params = append(params, "%"+kw+"%")
	params = append(params, "%"+kw+"%")
	params = append(params, "%"+kw+"%")

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
		sql += " and b.center_id=? "
	}

	if status != "" {
		params = append(params, status)
		sql += " and b.contact_status=? "
	}

	if centerId1 != "" && dataType == "all" {
		params = append(params, centerId1)
		sql += " and b.center_id=? "
	}

	if partTimeName != "" {
		params = append(params, partTimeName)
		sql += " and b.parttime_name=? "
	}

	if level != "" {
		params = append(params, level)
		sql += " and b.level=? "
	}

	sql += `group by a.consumer_id order by b.id desc  limit ?,?) c
	left join consumer_new cons on cons.id=c.consumer_id
	left join contacts cont on cont.id=c.contacts_id
	left join center ce on ce.cid=cons.center_id
	left join come_from cf on cf.id=cons.come_from_id
	left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) d on d.consumer_id=cons.id
	left join employee e on e.user_id=cons.last_tmk_id	`

	countSql := `select count(1) from(select a.consumer_id,min(a.id) contacts_id from contacts a left join consumer_new b on a.consumer_id=b.id
	where (a.name like ? or b.child like ? or a.phone like ? or b.home_phone like ?) and b.is_own_by_tmk=2 and b.pay_time is null and b.contact_status!=5 `

	paramsForCount = append(paramsForCount, "%"+kw+"%")
	paramsForCount = append(paramsForCount, "%"+kw+"%")
	paramsForCount = append(paramsForCount, "%"+kw+"%")
	paramsForCount = append(paramsForCount, "%"+kw+"%")

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

		paramsForCount = append(paramsForCount, _employee.CenterId)
		countSql += " and b.center_id=? "
	}

	if status != "" {
		paramsForCount = append(paramsForCount, status)
		countSql += " and b.contact_status=? "
	}

	if centerId1 != "" && dataType == "all" {
		paramsForCount = append(paramsForCount, centerId1)
		countSql += " and b.center_id=? "
	}

	if partTimeName != "" {
		paramsForCount = append(paramsForCount, partTimeName)
		countSql += " and b.parttime_name=? "
	}

	if level != "" {
		paramsForCount = append(paramsForCount, level)
		countSql += " and b.level=? "
	}

	countSql += " group by a.consumer_id) aa "

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, paramsForCount...)

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

		for i := 0; i < 13; i++ {
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

	if strings.Contains(ids, ",") {
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

	rows2, err := db.Query(sql2, employee.UserId)

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

	_, err = stmt.Exec(employee.UserId, ids, time.Now().Format("20060102150405"))

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

	_, err = stmt1.Exec(employee.UserId, employee.UserId, ids)

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
	partTimeName := r.FormValue("partTimeName-eq")
	level := r.FormValue("level-eq")

	params := []interface{}{}

	sql := " select cons.id,ce.name as centerName,cont.name,cons.level,cont.phone,cons.child,cons.birthday,cons.year,cons.contact_status,cons.parent_id,b.remark,cons.center_id,cons.pay_status,cons.pay_time,cf.name comeFromName,cons.parttime_name "
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
	sql += " left join come_from cf on cf.id=cons.come_from_id "
	sql += " left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note)  ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) b on b.consumer_id=cons.id "
	sql += " where tc.tmk_id= " + employee.UserId

	if status != "" {
		params = append(params, status)
		sql += " and cons.contact_status=? "
	}

	if lastContractStartTime != "" && timeType == "1" {
		params = append(params, lastContractStartTime)
		sql += " and cons.sign_in_time>=? "
	}

	if lastContractStartTime != "" && timeType == "2" {
		params = append(params, lastContractStartTime)
		sql += " and cons.pay_time>=? "
	}

	if lastContractEndTime != "" && timeType == "1" {
		params = append(params, lastContractEndTime)
		sql += " and cons.sign_in_time<=? "
	}

	if lastContractEndTime != "" && timeType == "2" {
		params = append(params, lastContractEndTime)
		sql += " and cons.pay_time<=? "
	}

	if partTimeName != "" {
		params = append(params, partTimeName)
		sql += " and cons.parttime_name=? "
	}

	if payStatus == "1" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=1 "
	} else if payStatus == "2" {
		sql += " and cons.pay_time is not null and cons.pay_time != '' and cons.pay_status=2 "
	} else if payStatus == "3" {
		sql += " and (cons.pay_time is null or cons.pay_time ='')"
	}

	if centerId1 != "" {
		params = append(params, centerId1)
		sql += " and cons.center_id=? "
	}

	if level != "" {
		params = append(params, level)
		sql += " and cons.level=? "
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

	sql += " order by tc.tmk_create_time desc,cons.contact_status ,cons.last_contact_time desc limit ?,? "

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

		for i := 0; i < 15; i++ {
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

	if strings.Contains(ids, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "每次只能缴费一个客户"
		commonlib.OutputJson(w, m, " ")
		return
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

	_, err = stmt.Exec(ids, time.Now().Format("20060102150405"), employee.UserId)

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

	_, err = stmt.Exec(time.Now().Format("20060102150405"), status, ids)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	getOldStatusSql := "select contact_status from consumer_new where id=? "
	lessgo.Log.Debug(getOldStatusSql)

	rows, err := db.Query(getOldStatusSql, ids)
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	oldStatus := CONSUMER_STATUS_NO_CONTACT

	if rows.Next() {
		err := rows.Scan(&oldStatus)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if oldStatus != CONSUMER_STATUS_SIGNIN {
		updateConsumerSql := "update consumer_new set contact_status=?,sign_in_time=? where id=?"
		lessgo.Log.Debug(updateConsumerSql)

		stmt, err = tx.Prepare(updateConsumerSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(CONSUMER_STATUS_SIGNIN, time.Now().Format("20060102150405"), ids)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertContactsStatusLog := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertContactsStatusLog)

		stmt, err = tx.Prepare(insertContactsStatusLog)
		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(ids, employee.UserId, time.Now().Format("20060102150405"), oldStatus, CONSUMER_STATUS_SIGNIN)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		getChildId := "select ch.cid from consumer_new cons left join child ch on ch.pid=cons.parent_id where cons.id=?"
		lessgo.Log.Debug(getChildId)

		rows, err = db.Query(getChildId, ids)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		childId := 0

		if rows.Next() {
			err := rows.Scan(&childId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if childId == 0 {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "该客户没有与小孩子关联，无法缴费"
			commonlib.OutputJson(w, m, " ")
			return
		}

		getClassIdSql := "select wyclass_id,schedule_detail_id from schedule_detail_child where child_id=? and wyclass_id is not null order by id desc "
		rows, err = db.Query(getClassIdSql, childId)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		classId := 0
		scheduleId := 0

		if rows.Next() {
			err := rows.Scan(&classId, &scheduleId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if classId == 0 { //找不到班级，就进行无班签到
			checkExistSql := "select count(1) from sign_in where child_id=? and wyclass_id is null and schedule_detail_id is null "

			lessgo.Log.Debug(checkExistSql)

			rows, err = db.Query(checkExistSql, childId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			num := 0

			if rows.Next() {

				err = commonlib.PutRecord(rows, &num)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门"
					commonlib.OutputJson(w, m, " ")
					return
				}
			}

			if num == 0 {
				insertWFSSql := "insert into sign_in(child_id,sign_time,employee_id,type) values(?,?,?,?)"
				lessgo.Log.Debug(insertWFSSql)

				insertWFSStmt, err := tx.Prepare(insertWFSSql)
				if err != nil {
					tx.Rollback()

					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门"
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = insertWFSStmt.Exec(childId, time.Now().Format("20060102150405"), employee.UserId, 1)
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

		} else {
			checkSignSql := "select count(1) from sign_in where child_id=? and wyclass_id=? "
			rows, err = db.Query(checkSignSql, childId, classId)

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

			if totalNum == 0 {
				insertSignInSql := "insert into sign_in(child_id,sign_time,schedule_detail_id,type,wyclass_id,employee_id) values(?,?,?,?,?,?)"
				lessgo.Log.Debug(insertSignInSql)

				stmt, err = tx.Prepare(insertSignInSql)
				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = stmt.Exec(childId, time.Now().Format("20060102150405"), scheduleId, 1, classId, employee.UserId)

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

			}
		}
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "缴费成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func BackToAllConsumerAction(w http.ResponseWriter, r *http.Request) {
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

	if strings.Contains(id, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "只能操作一个客户"
		commonlib.OutputJson(w, m, " ")
		return
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

	sql := "update consumer_new set is_own_by_tmk=2,current_tmk_id=null where id=? "
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

	_, err = stmt.Exec(id)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	sql = "delete from tmk_consumer where tmk_id=? and consumer_id=? "

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

	_, err = stmt.Exec(employee.UserId, id)

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
