// Title：Call Center跟踪统计
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
// 1.0 2013-10-15 10:46 black 创建文档
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

//CallCenter 统计报表
/*
select ce.cid,ce.name,a.num as '全部名单', b.num as '未联系', c.num as '待确认', d.num as '已废弃', e.num as '已邀约',f.num as '已签到'  from center ce
left join
(select count(1) num,center_id from consumer_new group by center_id )a on a.center_id=ce.cid
left join
(select count(1) num,center_id from consumer_new where contact_status=1 group by center_id )b on a.center_id=ce.cid
left join
(select count(1) num,center_id from consumer_new where contact_status=2 group by center_id )c on b.center_id=ce.cid
left join
(select count(1) num,center_id from consumer_new where contact_status=3 group by center_id )d on c.center_id=ce.cid
left join
(select count(1) num,center_id from consumer_new where contact_status=4 group by center_id )e on d.center_id=ce.cid
left join
(select count(1) num,center_id from consumer_new where contact_status=5 group by center_id )f on e.center_id=ce.cid
*/
func CallCenterStatisticsAction(w http.ResponseWriter, r *http.Request) {

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
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" {
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
		}
	}

	startTime := r.FormValue("startTime-ge")
	endTime := r.FormValue("endTime-le")

	params := []interface{}{}

	sql := " select ce.cid,ce.name,a.num as '全部名单', b.num as '未联系', c.num as '拨打电话数', d.num as '邀约数', e.num as '签到数',f.num as '定金',g.num as '全额',1,2,3  from center ce "
	sql += " left join "
	sql += " (select count(1) num,center_id from consumer_new group by center_id )a on a.center_id=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,center_id from consumer_new where contact_status=1 group by center_id )b on b.center_id=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,cid from audio where start_time >=? and start_time<=? group by cid )c on c.cid=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,ch.center_id from schedule_detail_child sdc left join child ch on sdc.child_id=ch.cid where wyclass_id is not null and sdc.create_time>=? and sdc.create_time <=? group by ch.center_id)d on d.center_id=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,ch.center_id from sign_in si left join child ch on si.child_id=ch.cid where si.type=1 and  (wyclass_id is not null or (wyclass_id is null and schedule_detail_id is null)) and si.sign_time>=? and si.sign_time<=? group by ch.center_id )e on e.center_id=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,center_id from consumer_new where pay_status=1 and pay_time>=? and pay_time<=? group by center_id )f on f.center_id=ce.cid "
	sql += " left join "
	sql += " (select count(1) num,center_id from consumer_new where pay_status=2 and pay_time>=? and pay_time<=? group by center_id )g on g.center_id=ce.cid where ce.cid!=9 " //屏蔽总部数据

	defaultStartTime := "2000-01-01 00:0:000"
	defaultEndTime := "2999-12-31 00:00:00"

	if startTime != "" {
		defaultStartTime = startTime
	}

	if endTime != "" {
		defaultEndTime = endTime
	}
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)

	defaultStartTime = strings.Replace(defaultStartTime, "-", "", -1)

	defaultEndTime = strings.Replace(defaultEndTime, "-", "", -1)

	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)
	params = append(params, defaultStartTime)
	params = append(params, defaultEndTime)

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)

	if dataType == "center" {
		sql += " and ce.cid=? "

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		params = append(params, _employee.CenterId)
	}

	countSql := "select count(1) from center where cid!=9 "

	if dataType == "center" {
		countSql += " and cid=" + _employee.CenterId
	}

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql)

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

	sql += " order by ce.cid limit ?,?"

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
