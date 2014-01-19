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

/*
select a.aid,cons.child,b.remark,cont.name contName,a.remotephone,ce.name centerName,a.start_time,a.seconds,a.inout,a.is_upload_finish,cons.id,a.filename,a.cid,e1.really_name
from
(
	select au.* from audio au
	left join employee e on e.phone_in_center=au.localphone and e.center_id=au.cid
	where e.user_id=37 and au.remotephone != '' and  au.remotephone is not null order by au.aid desc,au.start_time desc  limit 0,20
) a
left join contacts cont on a.remotephone=cont.phone
left join consumer_new cons on cont.consumer_id=cons.id
left join employee e2 on e2.phone_in_center=a.localphone and e2.center_id=a.cid
left join center ce on ce.cid=cons.center_id
left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id  SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) b on b.consumer_id=cons.id
left join  employee e1 on e1.user_id=cons.current_tmk_id
*/
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

	sql := `select a.aid,cons.child,b.remark,cont.name contName,a.remotephone,ce.name centerName,a.start_time,a.seconds,a.inout,a.is_upload_finish,cons.id,a.filename,a.cid,e1.really_name
			from
			(
				select au.* from audio au
				left join employee e on e.phone_in_center=au.localphone and e.center_id=au.cid
				where e.user_id=? and au.remotephone != '' and  au.remotephone is not null %v order by au.aid desc,au.start_time desc  limit ?,?
			) a
			left join contacts cont on a.remotephone=cont.phone
			left join consumer_new cons on cont.consumer_id=cons.id
			left join employee e2 on e2.phone_in_center=a.localphone and e2.center_id=a.cid
			left join center ce on ce.cid=cons.center_id
			%v
			left join  employee e1 on e1.user_id=cons.current_tmk_id`

	params = append(params, eid)

	whereSql := ""

	if flag {
		if st != "" && et != "" {
			whereSql += " and au.start_time >= ? and au.start_time<= ?"
			params = append(params, st)
			params = append(params, et)
		}
	} else { //找不到相应的时间区间
		whereSql += " and au.start_time >= ? and au.start_time<= ?"
		params = append(params, "2000-01-01 00:00:00")
		params = append(params, "2000-01-01 00:00:01")
	}

	countSql := `select count(1) from (
			select au.* from audio au
			left join employee e on e.phone_in_center=au.localphone and e.center_id=au.cid
			where e.user_id=? and au.remotephone != '' and  au.remotephone is not null %v) num`

	lessgo.Log.Debug(fmt.Sprintf(countSql, whereSql))

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(countSql, whereSql), params...)

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

	lessgo.Log.Debug(fmt.Sprintf(sql, whereSql, "left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id  SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) b on b.consumer_id=cons.id"))

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(fmt.Sprintf(sql, whereSql, "left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id  SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) b on b.consumer_id=cons.id"), params...)

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
