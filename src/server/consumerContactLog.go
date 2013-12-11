// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-22 17:59
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-22 17:59 black 创建文档
package server

import (
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func ConsumerContactLogAction(w http.ResponseWriter, r *http.Request) {

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

	consumerId := r.FormValue("consumerId")

	params := []interface{}{}

	sql := "select  c.id,c.create_time,e.really_name,c.note,c.type from consumer_contact_log c left join employee e on e.user_id=c.create_user where c.consumer_id=? "

	params = append(params, consumerId)

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

	sql += " order by c.id desc limit ?,?"

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

		for i := 0; i < 4; i++ {
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

func ConsumerContactLogSaveAction(w http.ResponseWriter, r *http.Request) {

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
	contactType := r.FormValue("contactType")
	consumerId := r.FormValue("consumer_id")
	note := r.FormValue("note")

	db := lessgo.GetMySQL()
	defer db.Close()

	if id == "" {

		getConsumerStatusSql := "select contact_status from consumer_new where id=? "
		lessgo.Log.Debug(getConsumerStatusSql)

		rows, err := db.Query(getConsumerStatusSql, consumerId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}

		consumerStatus := CONSUMER_STATUS_NO_CONTACT

		if rows.Next() {

			err = commonlib.PutRecord(rows, &consumerStatus)

			if err != nil {
				lessgo.Log.Error(err.Error())
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

		/**********************数据插入consumer表**********************/
		insertContactLogSql := "insert into consumer_contact_log(create_user,create_time,note,consumer_id,type) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertContactLogSql)
		stmt, err := tx.Prepare(insertContactLogSql)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), note, consumerId, contactType)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if consumerStatus == CONSUMER_STATUS_NO_CONTACT {
			updateConsumerStatsuSql := "update consumer_new set contact_status=? ,last_contact_time=? where id=? "
			lessgo.Log.Debug(updateConsumerStatsuSql)

			stmt, err = tx.Prepare(updateConsumerStatsuSql)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(CONSUMER_STATUS_WAIT, time.Now().Format("20060102150405"), consumerId)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			insertContactStatus := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"

			lessgo.Log.Debug(insertContactStatus)

			stmt, err = tx.Prepare(insertContactStatus)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(consumerId, employee.UserId, time.Now().Format("20060102150405"), CONSUMER_STATUS_NO_CONTACT, CONSUMER_STATUS_WAIT)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		} else {
			updateConsumerStatsuSql := "update consumer_new set last_contact_time=? where id=? "
			lessgo.Log.Debug(updateConsumerStatsuSql)

			stmt, err = tx.Prepare(updateConsumerStatsuSql)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(time.Now().Format("20060102150405"), consumerId)
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

		tx.Commit()

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	} else {

		sql := "update consumer_contact_log set note=?,type=? where id=? "

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

		_, err = stmt.Exec(note, contactType, id)

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

func ConsumerContactLogLoadAction(w http.ResponseWriter, r *http.Request) {

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

	loadFormObjects := []lessgo.LoadFormObject{}

	sql := "select type,note from consumer_contact_log where id=? "

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

	var contactType, note string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &contactType, &note)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	m["success"] = true

	h1 := lessgo.LoadFormObject{"contactType", contactType}
	h2 := lessgo.LoadFormObject{"note", note}
	h3 := lessgo.LoadFormObject{"id", id}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}
