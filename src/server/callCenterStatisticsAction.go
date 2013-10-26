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
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//CallCenter 统计报表
/*
select c.cid,c.name,count(tc.id) as '全部名单',count(tc.employee_id) as '已分配',count(tc.id)-count(tc.employee_id) as '未分配',count(aa.id) as '未联系',count(bb.id) as '待确认',count(cc.id) as '已废弃',count(dd.id) as '已邀约',count(ee.id) as '确认签到' from center c
left join tmk_consumer tc on c.cid=tc.center_id
left join (select * from consumer where contact_status=1)aa on tc.consumer_id=aa.id and tc.employee_id is not null and tc.employee_id!=0
left join (select * from consumer where contact_status=2)bb on tc.consumer_id=bb.id and tc.employee_id is not null and tc.employee_id!=0
left join (select * from consumer where contact_status=3)cc on tc.consumer_id=cc.id and tc.employee_id is not null and tc.employee_id!=0
left join (select * from consumer where contact_status=4)dd on tc.consumer_id=dd.id and tc.employee_id is not null and tc.employee_id!=0
left join (select * from consumer where contact_status=5)ee on tc.consumer_id=ee.id and tc.employee_id is not null and tc.employee_id!=0
group by c.cid
limit 0,100;
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

	sql := "select c.cid,c.name,count(tc.id) as '全部名单',count(tc.employee_id) as '已分配',count(tc.id)-count(tc.employee_id) as '未分配',count(aa.id) as '未联系',count(bb.id) as '待确认',count(cc.id) as '已废弃',count(dd.id) as '已邀约',count(ee.id) as '确认签到' from center c "
	sql += " left join tmk_consumer tc on c.cid=tc.center_id "
	sql += " left join (select * from consumer where contact_status=1)aa on tc.consumer_id=aa.id and tc.employee_id is not null and tc.employee_id!=0 "
	sql += " left join (select * from consumer where contact_status=2)bb on tc.consumer_id=bb.id and tc.employee_id is not null and tc.employee_id!=0 "
	sql += " left join (select * from consumer where contact_status=3)cc on tc.consumer_id=cc.id and tc.employee_id is not null and tc.employee_id!=0 "
	sql += " left join (select * from consumer where contact_status=4)dd on tc.consumer_id=dd.id and tc.employee_id is not null and tc.employee_id!=0 "
	sql += " left join (select * from consumer where contact_status=5)ee on tc.consumer_id=ee.id and tc.employee_id is not null and tc.employee_id!=0 where c.cid!=9 " //将总部过滤掉

	if dataType == "center" {
		sql += " and c.cid=? "
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
	} else if dataType == "self" {
		// fix
		sql += " and tc.employee_id=" + employee.UserId
	}

	sql += " group by c.cid "

	countSql := "select count(1) from center where cid!=9 "

	if dataType == "center" {
		countSql += " and cid=? "
	}

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
