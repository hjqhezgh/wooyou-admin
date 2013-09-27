// Title：音频相关服务
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
	"net/http"
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
	"strconv"
	"math"
	"math/rand"
	"time"
	"fmt"
	"text/template"
)

//顾问分页数据服务
func ConsultantPhoneListAction(w http.ResponseWriter,r *http.Request ) {

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

	dataType := r.FormValue("dataType")

	if dataType != "all" && dataType != "center" && dataType != "self" {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "非法的url请求"
		commonlib.OutputJson(w, m, " ")
		return
	}

	params := []interface{}{}

	sql := ""
	countSql := ""

	sqlForAllData := "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,phone_count.num b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone"
	sqlForCenterData := "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,phone_count.num b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio where cid=? group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone"
	sqlForEmployeeData := "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,phone_count.num b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone where e.user_id=?"

	if dataType=="all" {
		sql = sqlForAllData
	}else if  dataType=="center"{
		sql = sqlForCenterData
		userId,_ := strconv.Atoi(employee.UserId)
		_employee,err := FindEmployeeById(userId)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
		params = append(params,_employee.CenterId)
	}else if  dataType=="self"{
		sql = sqlForEmployeeData
		params = append(params,employee.UserId)
	}


	countSql = "select count(1) from (" +  sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql,params...)

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

	lessgo.Log.Debug(sql+" limit ?,?")

	params = append(params,(currPageNo-1)*pageSize)
	params = append(params,pageSize)

	rows, err = db.Query(sql+" limit ?,?", params...)

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

		for i:=0;i<8;i++{
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

