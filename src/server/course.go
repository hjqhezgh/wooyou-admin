// Title：课程相关服务
//
// Description:
//
// Author:black
//
// Createtime:2013-11-11 13:41
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-11 13:41 black 创建文档
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
	"time"
	"os"
	"io"
)

func CourseListAction(w http.ResponseWriter, r *http.Request) {

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

	dataType := ""

	roleIds := strings.Split(employee.RoleId, ",")

	for _, roleId := range roleIds {
		if roleId == "1" || roleId == "3" || roleId == "6" || roleId == "10" {
			dataType = "all"
			break
		} else {
			dataType = "center"
		}
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

	centerId := r.FormValue("cid-eq")

	params := []interface{}{}

	sql := "select  c.cid,c.name,ce.name as cename,c.price,t.name as tname,c.is_probation,c.begin_age,c.end_age,c.intro,c.lesson_num from course c left join center ce on ce.cid=c.center_id left join course_type t on t.tid=c.type where 1=1 "

	if dataType == "center" {
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
		sql += " and ce.cid=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		sql += " and c.center_id=? "
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

	sql += " order by c.cid desc limit ?,?"

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

		for i := 0; i < 9; i++ {
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

//根据中心id获取课程下拉数据
func CourseByCenterIdListAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	centerId := r.FormValue("id")

	sql := "select cid,name from course where center_id=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, centerId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	type Result struct {
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}

	objects := []*Result{}

	for rows.Next() {

		model := new(Result)

		err = commonlib.PutRecord(rows, &model.Value, &model.Desc)

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
	m["success"] = true
	m["datas"] = objects

	commonlib.OutputJson(w, m, " ")
	return
}

func CourseSaveAction(w http.ResponseWriter, r *http.Request) {

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
	name := r.FormValue("name")
	center_id := r.FormValue("center_id")
	price := r.FormValue("price")
	is_probation := r.FormValue("is_probation")
	typeString := r.FormValue("type")
	begin_age := r.FormValue("begin_age")
	end_age := r.FormValue("end_age")
	intro := r.FormValue("intro")
	app_display_level := r.FormValue("app_display_level")
	lesson_num := r.FormValue("lesson_num")
	courseTmpImg := r.FormValue("courseImg")

	db := lessgo.GetMySQL()
	defer db.Close()

	if id == "" {
		checkCourseExistSql := "select count(1) from course where center_id=? and name=? "
		lessgo.Log.Debug(checkCourseExistSql)

		rows, err := db.Query(checkCourseExistSql, center_id,name)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		totalNum := 0

		for rows.Next() {
			err = commonlib.PutRecord(rows, &totalNum)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if totalNum > 0 {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "课程名重复，请更换一个课程名重复"
			commonlib.OutputJson(w, m, " ")
			return
		}

		sql := "insert into course(name,center_id,price,is_probation,type,begin_age,end_age,intro,app_display_level,create_time,lesson_num) values(?,?,?,?,?,?,?,?,?,?,?)"

		lessgo.Log.Debug(sql)

		stmt, err := db.Prepare(sql)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		res, err := stmt.Exec(name, center_id, price, is_probation, typeString, begin_age, end_age, intro,app_display_level,time.Now().Format("20060102150405"),lesson_num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		courseId ,err := res.LastInsertId()

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		if courseTmpImg!="" {
			tmpFile, err := os.OpenFile(".."+courseTmpImg, os.O_RDWR, 0777)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			courseImgDir ,err := lessgo.Config.GetValue("wooyou", "courseImgDir")

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			err = os.MkdirAll(fmt.Sprint(courseImgDir+"/",courseId), 0777)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			disFile, err := os.Create(fmt.Sprint(courseImgDir+"/",courseId,"/480_230.png"))

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			io.Copy(disFile, tmpFile)

			os.Remove(".."+courseTmpImg)
		}


		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	} else {

		sql := "update course set name=?,center_id=?,price=?,is_probation=?,type=?,begin_age=?,end_age=?,intro=?,lesson_num=? where cid=? "

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

		_, err = stmt.Exec(name, center_id, price, is_probation, typeString, begin_age, end_age, intro,lesson_num,id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}


		if courseTmpImg != "" {
			tmpFile, err := os.OpenFile(".."+courseTmpImg, os.O_RDWR, 0777)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			courseImgDir ,err := lessgo.Config.GetValue("wooyou", "courseImgDir")

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			err = os.MkdirAll(fmt.Sprint(courseImgDir+"/",id), 0777)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			disFile, err := os.Create(fmt.Sprint(courseImgDir+"/",id,"/480_230.png"))

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			io.Copy(disFile, tmpFile)

			os.Remove(".."+courseTmpImg)
		}


		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	}

}


func CourseLoadAction(w http.ResponseWriter, r *http.Request) {

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

	sql := "select name,center_id,price,is_probation,type,begin_age,end_age,intro,lesson_num from course where cid=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	var name, centerId, price, isProbation, typeString, beginAge, endAge,intro,lessonNum string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &name, &centerId, &price, &isProbation, &typeString, &beginAge, &endAge, &intro, &lessonNum)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	m["success"] = true

	loadFormObjects := []lessgo.LoadFormObject{}

	h1 := lessgo.LoadFormObject{"id", id}
	h2 := lessgo.LoadFormObject{"name", name}
	h3 := lessgo.LoadFormObject{"center_id", centerId}
	h4 := lessgo.LoadFormObject{"price", price}
	h5 := lessgo.LoadFormObject{"is_probation", isProbation}
	h6 := lessgo.LoadFormObject{"type", typeString}
	h7 := lessgo.LoadFormObject{"begin_age", beginAge}
	h8 := lessgo.LoadFormObject{"end_age", endAge}
	h9 := lessgo.LoadFormObject{"intro", intro}
	h10 := lessgo.LoadFormObject{"lesson_num", lessonNum}
	h11 := lessgo.LoadFormObject{"courseImg", fmt.Sprint("http://app.wooyou.com.cn:9100/pic/course/",id,"/480_230.png")}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)
	loadFormObjects = append(loadFormObjects, h9)
	loadFormObjects = append(loadFormObjects, h10)
	loadFormObjects = append(loadFormObjects, h11)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}
