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
)

//顾问分页数据服务
func ConsultantPhoneDetailListAction(w http.ResponseWriter, r *http.Request) {

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

	eid := r.FormValue("eid")
	year := r.FormValue("year-eq")
	month := r.FormValue("month-eq")
	week := r.FormValue("week-eq")
	startTime := r.FormValue("start_time-eq")

	st := ""
	et := ""
	flag := true

	if startTime != "" {
		st = startTime + " 00:00:00"
		et = startTime + " 23:59:59"
	} else {
		if week != "" && month != "" && year != "" {
			st, et, flag = lessgo.FindRangeTimeDim("", "", year+month+week)
		} else if month != "" && year != "" {
			st, et, flag = lessgo.FindRangeTimeDim("", year+month, "")
		} else if year != "" {
			st, et, flag = lessgo.FindRangeTimeDim(year, "", "")
		}
	}

	params := []interface{}{}

	sql := "select a.aid,case a.remotephone when c.father_phone then c.father when c.mother_phone then c.mother  else '未知客户' end as c_name,a.remotephone,e.really_name,a.start_time,a.seconds,a.inout,a.is_upload_finish,a.note,c.id,a.filename,a.cid from audio a left join consumer c on (a.remotephone=c.mother_phone and c.mother_phone!='' and c.mother_phone is not null ) or (a.remotephone=c.father_phone and c.father_phone!='' and  c.father_phone is not null) left join employee e on e.phone_in_center=a.localphone where a.remotephone !='' and a.remotephone is not null and e.user_id=? "

	params = append(params, eid)

	if flag {
		if st != "" && et != "" {
			sql += " and a.start_time >= ? and a.start_time<= ?"
			params = append(params, st)
			params = append(params, et)
		}
	} else { //找不到相应的时间区间
		sql += " and a.start_time >= ? and a.start_time<= ?"
		params = append(params, "2000-01-01 00:00:00")
		params = append(params, "2000-01-01 00:00:01")
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
