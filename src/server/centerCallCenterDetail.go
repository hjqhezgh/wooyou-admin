// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-10-16 14:24
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-16 14:24 black 创建文档
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
)

func CenterCallCenterDetailAction(w http.ResponseWriter, r *http.Request) {

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

	status := r.FormValue("status-eq")
	name := r.FormValue("name-eq")
	centerId := r.FormValue("id")

	dataType := ""

	roleIds := strings.Split(employee.RoleId, ",")

	for _, roleId := range roleIds {
		if roleId == "1" || roleId == "3" || roleId == "6" || roleId == "10" {
			dataType = "all"
			break
		} else if roleId == "2" {
			dataType = "center"
			break
		} else {
			dataType = "self"
		}
	}

	params := []interface{}{}

	sql := ""
	countSql := ""

	sql += "select tc.id,e.really_name,c.mother,c.mother_phone,c.father,c.father_phone,c.home_phone,c.child,c.contact_status,c.id as 'cid',c.parent_id from tmk_consumer tc left join consumer c on tc.consumer_id=c.id left join employee e on e.user_id = tc.employee_id where tc.center_id=? "

	if dataType == "all" {
		params = append(params, centerId)

		if status != "" {
			if status == "no employee" {
				sql += " and tc.employee_id is null or tc.employee_id=0 "
			} else {
				sql += " and c.contact_status=? "
				params = append(params, status)
			}
		}

	} else if dataType == "center" {
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

		if status != "" {
			if status == "no employee" {
				sql += " and tc.employee_id is null or tc.employee_id=0 "
			} else {
				sql += " and c.contact_status=? "
				params = append(params, status)
			}
		}

	} else {

		params = append(params, centerId)

		sql += " and tc.employee_id=? "
		params = append(params, employee.UserId)

		if status != "" && status != "no employee" {
			sql += " and c.contact_status=? "
			params = append(params, status)
		}
	}

	if name != "" {
		sql += " and (c.mother like ? or c.mother_phone like ? or c.father like ? or c.father_phone like ? or c.home_phone like ? or c.child like ? )"
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
	}

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

	sql += " order by tc.cd_create_time desc  limit ?,?"

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

		for i := 0; i < 10; i++ {
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
