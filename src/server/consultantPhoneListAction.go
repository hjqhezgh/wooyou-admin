// Title：顾问通话总览列表
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
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//顾问分页数据服务
func ConsultantPhoneListAction(w http.ResponseWriter, r *http.Request) {

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
		if roleId == "1" || roleId == "3" {
			dataType = "all"
			break
		} else if roleId == "2" {
			dataType = "center"
			break
		} else if roleId == "" {
			dataType = "self"
		}
	}

	cid := r.FormValue("cid-eq")
	name := r.FormValue("name-like")
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

	sql := ""
	countSql := ""

	if dataType == "all" {

		sql += "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,rank.rowNo b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio where remotephone!='' and remotephone is not null "

		if cid != "" {
			sql += " and cid=? "
			params = append(params, cid)
		}

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone left join (select a.*,(@rowNum:=@rowNum+1) as rowNo from (select count(*) num,localphone from audio where remotephone!='' and remotephone is not null "

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone order by num desc) a,(Select (@rowNum :=0) ) b)rank on rank.localphone=phone_count.localphone "

		if name != "" {
			sql += " where e.really_name like ? "
			params = append(params, "%"+name+"%")
		}

		sql += " order by rank.rowNo "

	} else if dataType == "center" {

		sql += "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,rank.rowNo b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio  where cid=? and remotephone!='' and remotephone is not null "

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

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone  left join (select a.*,(@rowNum:=@rowNum+1) as rowNo from (select count(*) num,localphone from audio where remotephone!='' and remotephone is not null "

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone order by num desc) a,(Select (@rowNum :=0) ) b)rank on rank.localphone=phone_count.localphone "

		if name != "" {
			sql += " where e.really_name like ? "
			params = append(params, "%"+name+"%")
		}

		sql += " order by rank.rowNo "

	} else if dataType == "self" {

		sql += "select c.name,c.cid,e.user_id,e.really_name,phone_count.num a,rank.rowNo b,phone_count.num c,phone_count.num d from (select count(*) num,localphone,cid from audio where remotephone!='' and remotephone is not null "

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone) phone_count left join center c on c.cid=phone_count.cid left join employee e on e.phone_in_center=phone_count.localphone  left join (select a.*,(@rowNum:=@rowNum+1) as rowNo from (select count(*) num,localphone from audio where remotephone!='' and remotephone is not null "

		if flag {
			if st != "" && et != "" {
				sql += " and start_time >= ? and start_time<= ?"
				params = append(params, st)
				params = append(params, et)
			}
		} else { //找不到相应的时间区间
			sql += " and start_time >= ? and start_time<= ?"
			params = append(params, "2000-01-01 00:00:00")
			params = append(params, "2000-01-01 00:00:01")
		}

		sql += " group by  localphone order by num desc) a,(Select (@rowNum :=0) ) b)rank on rank.localphone=phone_count.localphone "

		sql += " where e.user_id=? "

		params = append(params, employee.UserId)

		sql += " order by rank.rowNo "
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

	lessgo.Log.Debug(sql + " limit ?,?")

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

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

		for i := 0; i < 8; i++ {
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
