// Title：顾问通话详情列表
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
	"text/template"
	//"strings"
)



//顾问分页数据服务
func VideoListAction(w http.ResponseWriter, r *http.Request) {

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

	roleIds := strings.Split(employee.RoleId, ",")

	for _, roleId := range roleIds {
		if roleId == "1" || roleId == "3" || roleId == "6" || roleId=="10"{
			dataType = "all"
			break
		} else if roleId == "2" {
			dataType = "center"
			break
		} else{
			dataType = "self"
		}
	}

	//fmt.Println("dataType:" ,dataType, ",roleIds:", roleIds)
	//r.FormValue("cid-eq")
	//fmt.Println(r.Form)

	sql := ""

	params := []interface{}{}

	if dataType == "all" {
		sql += `select v.vid,v.cid,ce.name as centername,r.name rname,c.course_id,co.name,c.teacher_id,e.really_name,v.start_time,v.end_time from
				video v left join class_schedule_detail c on v.schedule_detail_id=c.id left join room r on r.rid=v.rid left join employee e
					on c.teacher_id=e.user_id left join course co on c.course_id=co.cid left join center ce on v.cid=ce.cid
						where data_rel_status=2`
	} else if dataType == "center" {
		sql +=`select v.vid,v.cid,ce.name as centername,r.name rname,c.course_id,co.name,c.teacher_id,e.really_name,v.start_time,v.end_time from
				video v left join class_schedule_detail c on v.schedule_detail_id=c.id left join room r on r.rid=v.rid left join employee e
					on c.teacher_id=e.user_id left join course co on c.course_id=co.cid left join center ce on v.cid=ce.cid
						where v.cid=? and data_rel_status=2`
		userId, _ := strconv.Atoi(employee.UserId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		params = append(params,  _employee.CenterId)

	} else if dataType == "self" {
		sql +=`select v.vid,v.cid,ce.name as centername,r.name rname,c.course_id,co.name,c.teacher_id,e.really_name,v.start_time,v.end_time from
			video v left join class_schedule_detail c on v.schedule_detail_id=c.id left join room r on r.rid=v.rid left join employee e
				on c.teacher_id=e.user_id left join course co on c.course_id=co.cid left join center ce on v.cid=ce.cid
					where c.teacher_id=? and data_rel_status=2`

		params = append(params, employee.UserId)
	}

	center_id := r.FormValue("cid-eq")
	if center_id != "" {
		sql += " and v.cid=?"
		params = append(params, center_id)
		lessgo.Log.Debug(sql)
		lessgo.Log.Debug(params)
	}

	course_id := r.FormValue("courseId-eq")
	if course_id != "" {
		sql += " and c.course_id=?"
		params = append(params, course_id)
		lessgo.Log.Debug(params...)
	}

	name := r.FormValue("name-like")
	if name != "" {
		sql += " and e.really_name like ?"
		params = append(params, "%"+name+"%")
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


	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	lessgo.Log.Debug(sql + " limit ?,?")
	lessgo.Log.Debug(params...)
	lessgo.Log.Debug(params[0], ":",params[1])
	rows, err = db.Query(sql+" limit ?,?", params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}
	colums, _ := rows.Columns()
	column_len := len(colums) - 1
	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < column_len; i++ {
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
