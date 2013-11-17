// Title：客户联系人
//
// Description:
//
// Author:black
//
// Createtime:2013-11-16 17:29
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-16 17:29 black 创建文档
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

func ContractsListAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

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

	sql := " select id,name,phone,is_default from contacts where consumer_id=? "

	params = append(params, consumerId)

	countSql := ""

	countSql = "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, consumerId)

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

	sql += " order by is_default,id desc limit ?,?"

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

		for i := 0; i < 3; i++ {
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

