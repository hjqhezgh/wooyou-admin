// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-11 14:23
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-11 14:23 black 创建文档
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
)

type WyClass struct {
	ClassId     int `json:"classId"`
	AssistantId int
	Name        string `json:"name"`
	CourseId    int
	CenterId    int
	ChildNum    int
	EndTime     string
	DeadLine    string
	MaxChildNum int
	TeacherId   int
	IsProbation string
}

func WyClassListAction(w http.ResponseWriter, r *http.Request) {

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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "tmk" {
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
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

	sql := "select  c.class_id,c.name,ce.name as cename,cou.name as courseName,c.end_time,c.deadline,c.max_child_num,tea.really_name as teacherName,ass.really_name as assName,c.center_id from wyclass c left join center ce on ce.cid=c.center_id left join employee tea on tea.user_id=c.teacher_id left join employee ass on ass.user_id=c.assistant_id left join course cou on c.course_id=cou.cid where 1=1 and (c.start_time is null or c.start_time = '') "

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

	sql += " order by c.class_id desc limit ?,?"

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

func WyClassLoadAction(w http.ResponseWriter, r *http.Request) {

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

	loadFormObjects := []lessgo.LoadFormObject{}

	if id != "" {
		sql := "select class_id,name,code,start_time from wyclass where class_id=? "

		lessgo.Log.Debug(sql)

		db := lessgo.GetMySQL()
		defer db.Close()

		rows, err := db.Query(sql, id)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		var classId, name, code, startTime string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &classId, &name, &code, &startTime)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		m["success"] = true

		h1 := lessgo.LoadFormObject{"name", name}
		h2 := lessgo.LoadFormObject{"class_id", classId}
		h3 := lessgo.LoadFormObject{"code", code}
		h4 := lessgo.LoadFormObject{"start_time", startTime}

		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
		loadFormObjects = append(loadFormObjects, h3)
		loadFormObjects = append(loadFormObjects, h4)

		m["datas"] = loadFormObjects
		commonlib.OutputJson(w, m, " ")
		return
	} else {
		userId, _ := strconv.Atoi(employee.UserId)
		_employee, err := FindEmployeeById(userId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		h1 := lessgo.LoadFormObject{"center_id", _employee.CenterId}
		loadFormObjects = append(loadFormObjects, h1)
		m["datas"] = loadFormObjects
		m["success"] = true
		commonlib.OutputJson(w, m, " ")
		return
	}

}

func WyClassSaveAction(w http.ResponseWriter, r *http.Request) {

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

	id := r.FormValue("class_id")
	name := r.FormValue("name")
	center_id := r.FormValue("center_id")
	code := r.FormValue("code")
	startTime := r.FormValue("start_time")

	db := lessgo.GetMySQL()
	defer db.Close()

	if id == "" {
		sql := "insert into wyclass(name,create_time,center_id,code,start_time) values(?,?,?,?,?)"

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

		_, err = stmt.Exec(name, time.Now().Format("20060102150405"), center_id, code, startTime)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
		return
	} else {
		sql := "update wyclass set name=?,code=? where class_id=? "

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

		_, err = stmt.Exec(name, code, id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
		return
	}

}

func LoadChildInClassAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")

	sql := "select c.cid,c.name from class_child cc left join child c on cc.child_id=c.cid where cc.class_id=? "

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

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

func SaveChildToClassAction(w http.ResponseWriter, r *http.Request) {

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
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	classId := r.FormValue("id")
	studIdsString := r.FormValue("cids")
	noidsString := r.FormValue("noids")

	db := lessgo.GetMySQL()
	defer db.Close()

	class, err := FindClassById(classId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	studIds := strings.Split(studIdsString, ",")
	if len(studIds) > class.MaxChildNum {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "班级人数超出可以容纳的人数"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if class.IsProbation == "1" { //免费课
		for _, studId := range studIds {
			findDataInMiddleTableSql := "select count(1) from class_child where class_id=? and child_id=? "
			lessgo.Log.Debug(findDataInMiddleTableSql)
			findDataInMiddleTableRows, err := db.Query(findDataInMiddleTableSql, classId, studId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			middleTableHaveData := 0

			if findDataInMiddleTableRows.Next() {
				err = commonlib.PutRecord(findDataInMiddleTableRows, &middleTableHaveData)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}
			}

			if middleTableHaveData == 0 { //中间表没数据
				insertContractSql := "insert into contract(child_id,apply_time,parent_id,price,employee_id,center_id,course_id,left_lesson_num,type,status) values(?,?,?,?,?,?,?,?,?,?)"
				lessgo.Log.Debug(insertContractSql)

				child, err := FindChildById(studId)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				insertStmt, err := tx.Prepare(insertContractSql)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				res, err := insertStmt.Exec(studId, time.Now().Format("20060102150405"), child.Pid, 0, employee.UserId, child.CenterId, class.CourseId, 1, "1", "1")

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				contractId, err := res.LastInsertId()

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				insertMiddleSql := "insert into class_child values(?,?,?)"
				lessgo.Log.Debug(insertContractSql)

				insertMiddleStmt, err := tx.Prepare(insertMiddleSql)
				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = insertMiddleStmt.Exec(studId, class.ClassId, contractId)

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				getConsumerIdSql := "select cons.id from consumer_new cons left join parent p on p.pid=cons.parent_id left join child c on c.pid=p.pid where c.cid=? "
				lessgo.Log.Debug(getConsumerIdSql)

				getConsumerRow, err := db.Query(getConsumerIdSql, child.Cid)

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				consumerId := 0

				if getConsumerRow.Next() {
					err = commonlib.PutRecord(getConsumerRow, &consumerId)

					if err != nil {
						tx.Rollback()
						lessgo.Log.Warn(err.Error())
						m["success"] = false
						m["code"] = 100
						m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
						commonlib.OutputJson(w, m, " ")
						return
					}
				}

				updateConsumerStatusSql := "update consumer_new set contact_status=4 where id=? "
				lessgo.Log.Debug(insertContractSql)

				updateConsumerStatusSqlStmt, err := tx.Prepare(updateConsumerStatusSql)
				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = updateConsumerStatusSqlStmt.Exec(consumerId)

				if err != nil {
					tx.Rollback()
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

			} else {
				continue
			}
		}

		noIds := strings.Split(noidsString, ",")

		for _, noId := range noIds {

			// todo 免费课的合同。如果没上过的，就把这个合同删除

			deleteMiddleSql := "delete from class_child where class_id=? and child_id=?"
			lessgo.Log.Debug(deleteMiddleSql)

			stmt, err := tx.Prepare(deleteMiddleSql)
			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(class.ClassId, noId)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		tx.Commit()

		m["success"] = true
		m["code"] = 200
		m["msg"] = "小孩子分配成功"
		commonlib.OutputJson(w, m, " ")
		return

	} else {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "暂时不支持开设常规班级，敬请期待"
		commonlib.OutputJson(w, m, " ")
		return
	}

}

func FindClassById(id string) (WyClass, error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select wc.class_id,wc.assistant_id,wc.name,wc.course_id,wc.center_id,wc.child_num,wc.end_time,wc.deadline,wc.max_child_num,wc.teacher_id,c.is_probation from wyclass wc left join course c on wc.course_id=c.cid where wc.class_id=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		return WyClass{}, err
	}

	class := WyClass{}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &class.ClassId, &class.AssistantId, &class.Name, &class.CourseId, &class.CenterId, &class.ChildNum, &class.EndTime, &class.DeadLine, &class.MaxChildNum, &class.TeacherId, &class.IsProbation)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return WyClass{}, err
		}
	}

	return class, nil
}

func ClassByCenterIdAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	id := r.FormValue("id")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select class_id,name from wyclass where center_id=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	classes := []WyClass{}

	for rows.Next() {
		wyClass := WyClass{}

		err := commonlib.PutRecord(rows, &wyClass.ClassId, &wyClass.Name)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		classes = append(classes, wyClass)
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = classes

	commonlib.OutputJson(w, m, " ")
}

//免费课列表
/*
select class.class_id,ce.name,class.code,class.name,class.start_time,childNum.num,signInNum.num
from wy_class class
left join (select count(1) num,wyclass_id from wyclass_free_child group by wyclass_id) childNum on class.class_id=childNum.wyclass_id
left join (select count(1) num,wyclass_id from wyclass_free_sign_in group by wyclass_id)signInNUm on class.class_id=signInNum.wyclass_id
left join center ce on ce.cid=class.center_id
*/
func WyClassFreeListAction(w http.ResponseWriter, r *http.Request) {

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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "tmk" {
			dataType = "all"
			break
		} else {
			dataType = "center"
			break
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
	startTime := r.FormValue("startTime-ge")
	endTime := r.FormValue("endTime-le")
	kw := r.FormValue("kw-eq")

	params := []interface{}{}

	sql := "select csd.id,ce.name as centerName,class.code,class.name,csd.start_time,csd.class_id,csd.center_id "
	sql += " from class_schedule_detail csd "
	sql += " left join wyclass class on csd.class_id=class.class_id "
	sql += " left join center ce on ce.cid=csd.center_id where csd.class_id is not null  and csd.start_time is not null and csd.start_time != '' "

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
		sql += " and csd.center_id=? "
	}

	if centerId != "" && dataType == "all" {
		params = append(params, centerId)
		sql += " and csd.center_id=? "
	}

	if startTime != "" {
		params = append(params, startTime)
		sql += " and csd.start_time>=? "
	}

	if endTime != "" {
		params = append(params, endTime)
		sql += " and csd.start_time<=? "
	}

	if kw != "" {
		sql += " and csd.class_id in (select sdc.wyclass_id from child ch left join schedule_detail_child sdc "
		sql += " on sdc.child_id=ch.cid "
		sql += " left join parent p on p.pid=ch.pid "
		sql += " where ch.name=? or p.telephone=?) "
		params = append(params, kw)
		params = append(params, kw)
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

	sql += " order by csd.id desc limit ?,?"

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

		for i := 0; i < 6; i++ {
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

//学员手动签到
func ChildSignInAction(w http.ResponseWriter, r *http.Request) {

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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "tmk" {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "对不起，您没有权限进行签到"
			commonlib.OutputJson(w, m, " ")
			return
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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	idList := strings.Split(ids, ",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for _, id := range idList {

		checkExistSql := "select count(1) from sign_in where child_id=? and wyclass_id=? and schedule_detail_id=?"

		lessgo.Log.Debug(checkExistSql)

		rows, err := db.Query(checkExistSql, id, classId, scheduleId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		num := 0

		if rows.Next() {

			err = commonlib.PutRecord(rows, &num)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if num > 0 {
			continue
		}

		insertWFSSql := "insert into sign_in(child_id,sign_time,wyclass_id,employee_id,schedule_detail_id,type) values(?,?,?,?,?,?)"

		lessgo.Log.Debug(insertWFSSql)

		insertWFSStmt, err := tx.Prepare(insertWFSSql)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = insertWFSStmt.Exec(id, time.Now().Format("20060102150405"), classId, employee.UserId, scheduleId, 1)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		getConsumerIdSql := "select cons.id from consumer_new cons left join parent p on p.pid=cons.parent_id left join child ch on ch.pid=p.pid where ch.cid=?"
		lessgo.Log.Debug(getConsumerIdSql)

		rows, err = db.Query(getConsumerIdSql, id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		consumerId := 0

		if rows.Next() {
			err = commonlib.PutRecord(rows, &consumerId)
			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if consumerId != 0 {
			updateConsuemrStatusSql := "update consumer_new set contact_status=?,sign_in_time=? where id=? "

			lessgo.Log.Debug(updateConsuemrStatusSql)

			updateConsuemrStatusStmt, err := tx.Prepare(updateConsuemrStatusSql)

			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = updateConsuemrStatusStmt.Exec(CONSUMER_STATUS_SIGNIN, time.Now().Format("20060102150405"), consumerId)
			if err != nil {
				tx.Rollback()

				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			insertLogSql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"

			stmt, err := tx.Prepare(insertLogSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(consumerId, employee.UserId, time.Now().Format("20060102150405"), CONSUMER_STATUS_NO_SIGNIN, CONSUMER_STATUS_SIGNIN)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											select ?,?,concat('签到',wc.start_time,wc.name),?,3 from wyclass wc where wc.class_id=?`
			lessgo.Log.Debug(insertConsumerContactLogSql)

			stmt, err = tx.Prepare(insertConsumerContactLogSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), consumerId, classId)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func ChildSignInWithoutAction(w http.ResponseWriter, r *http.Request) {

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

	roleCodes := strings.Split(employee.RoleCode, ",")

	for _, roleCode := range roleCodes {
		if roleCode == "tmk" {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "对不起，您没有权限进行签到"
			commonlib.OutputJson(w, m, " ")
			return
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

	ids := r.FormValue("ids")

	idList := strings.Split(ids, ",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for _, id := range idList {

		getChildIdSql := "select ch.cid from consumer_new cons left join parent p on p.pid=cons.parent_id left join child ch on ch.pid=p.pid where cons.id=?"
		lessgo.Log.Debug(getChildIdSql)

		rows, err := db.Query(getChildIdSql, id)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		chilId := 0;
		if rows.Next() {

			err = commonlib.PutRecord(rows, &chilId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if chilId==0{
			m["success"] = false
			m["code"] = 100
			m["msg"] = "该客户无小孩子关联，无法无班签到"
			commonlib.OutputJson(w, m, " ")
			return
		}

		checkExistSql := "select contact_status from consumer_new where id=? "

		lessgo.Log.Debug(checkExistSql)

		rows, err = db.Query(checkExistSql, id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		nowStatus := "1"

		if rows.Next() {

			err = commonlib.PutRecord(rows, &nowStatus)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if nowStatus == CONSUMER_STATUS_SIGNIN {
			continue
		}

		insertWFSSql := "insert into sign_in(child_id,sign_time,employee_id,type) values(?,?,?,?)"

		lessgo.Log.Debug(insertWFSSql)

		insertWFSStmt, err := tx.Prepare(insertWFSSql)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = insertWFSStmt.Exec(chilId, time.Now().Format("20060102150405"), employee.UserId, 1)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		updateConsuemrStatusSql := "update consumer_new set contact_status=?,sign_in_time=? where id=? "

		lessgo.Log.Debug(updateConsuemrStatusSql)

		updateConsuemrStatusStmt, err := tx.Prepare(updateConsuemrStatusSql)

		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = updateConsuemrStatusStmt.Exec(CONSUMER_STATUS_SIGNIN, time.Now().Format("20060102150405"), id)
		if err != nil {
			tx.Rollback()

			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertLogSql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"

		stmt, err := tx.Prepare(insertLogSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(id, employee.UserId, time.Now().Format("20060102150405"), CONSUMER_STATUS_NO_CONTACT, CONSUMER_STATUS_SIGNIN)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}


		insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											values(?,?,'无班签到',?,3)`
		lessgo.Log.Debug(insertConsumerContactLogSql)

		stmt, err = tx.Prepare(insertConsumerContactLogSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), id)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

//从班级中移除
func RemoveChildFromClassAction(w http.ResponseWriter, r *http.Request) {

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

	roleCodes := strings.Split(employee.RoleCode, ",")

	dataType := ""

	for _, roleCode := range roleCodes {
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "cd" {
			dataType = "all"
			break
		} else {
			dataType = "self"
			break
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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	if strings.Contains(ids, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该操作只能操作一个学员"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	checkExistSql := "select count(1) from sign_in where child_id=? and wyclass_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(checkExistSql)

	rows, err := db.Query(checkExistSql, ids, classId, scheduleId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if num > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已签到的人员不能移除"
		commonlib.OutputJson(w, m, " ")
		return
	}

	if dataType != "all" { //主管级别的是可以直接从班级中剔除任何人的，其他人只能剔除自己邀约的
		selectInviteTmk := "select create_user from schedule_detail_child where child_id=? and wyclass_id=? and schedule_detail_id=?"

		lessgo.Log.Debug(selectInviteTmk)

		rows, err = db.Query(selectInviteTmk, ids, classId, scheduleId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		employeeId := ""

		if rows.Next() {

			err = commonlib.PutRecord(rows, &employeeId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if employeeId != employee.UserId {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "您不能剔除其他人邀请的学员"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	deteleSql := "delete from schedule_detail_child where child_id=? and wyclass_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(deteleSql)

	stmt, err := tx.Prepare(deteleSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(ids, classId, scheduleId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	getConsumerIdSql := "select cons.id from consumer_new cons left join parent p on p.pid=cons.parent_id left join child ch on ch.pid=p.pid where ch.cid=?"
	lessgo.Log.Debug(getConsumerIdSql)

	rows, err = db.Query(getConsumerIdSql, ids)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	consumerId := 0

	if rows.Next() {
		err = commonlib.PutRecord(rows, &consumerId)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if consumerId != 0 {
		updateConsumerSql := "update consumer_new set contact_status=? where id=? "
		lessgo.Log.Debug(updateConsumerSql)
		stmt, err = tx.Prepare(updateConsumerSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(CONSUMER_STATUS_WAIT, consumerId)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertLogSql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertLogSql)
		stmt, err = tx.Prepare(insertLogSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(consumerId, employee.UserId, time.Now().Format("20060102150405"), CONSUMER_STATUS_NO_SIGNIN, CONSUMER_STATUS_WAIT)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											select ?,?,concat('从',wc.start_time,wc.name,'中剔除'),?,3 from wyclass wc where wc.class_id=?`
		lessgo.Log.Debug(insertConsumerContactLogSql)

		stmt, err = tx.Prepare(insertConsumerContactLogSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), consumerId, classId)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

//从班级中移除
func RemoveChildForNormalOnceAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	if strings.Contains(ids, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该操作只能操作一个学生"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	checkExistSql := "select count(1) from sign_in where child_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(checkExistSql)

	rows, err := db.Query(checkExistSql, ids, scheduleId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if num > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已签到的人员不能移除"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	deteleSql := "delete from schedule_detail_child where child_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(deteleSql)

	stmt, err := tx.Prepare(deteleSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(ids, scheduleId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func RemoveChildForNormalAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	if strings.Contains(ids, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该操作只能操作一个学生"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	checkExistSql := "select count(1) from sign_in where child_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(checkExistSql)

	rows, err := db.Query(checkExistSql, ids, scheduleId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if num > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已签到的人员不能移除"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	deteleSql := "delete from schedule_detail_child where child_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(deteleSql)

	stmt, err := tx.Prepare(deteleSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(ids, scheduleId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	getTmpIdSql := "select st.id from class_schedule_detail csd left join schedule_template st on csd.center_id=st.center_id and csd.week=st.week and csd.room_id=st.room_id and csd.time_id=st.time_id where csd.id=?"
	lessgo.Log.Debug(getTmpIdSql)

	rows, err = db.Query(getTmpIdSql, scheduleId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	tmpId := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &tmpId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if tmpId != 0 {
		deleteTmpChildSql := "delete from schedule_template_child where schedule_template_id=? and child_id=? "
		lessgo.Log.Debug(deleteTmpChildSql)

		stmt, err = tx.Prepare(deleteTmpChildSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(tmpId, ids)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func LeaveForChildAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	if strings.Contains(ids, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该操作只能操作一个学生"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	checkExistSql := "select count(1) from sign_in where child_id=? and schedule_detail_id=?"

	lessgo.Log.Debug(checkExistSql)

	rows, err := db.Query(checkExistSql, ids, scheduleId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if num > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已签到的学生不能请假"
		commonlib.OutputJson(w, m, " ")
		return
	}

	getContractId := "select contract_id from schedule_detail_child where child_id=? and schedule_detail_id=? "
	lessgo.Log.Debug(checkExistSql)

	rows, err = db.Query(getContractId, ids, scheduleId)
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	contractId := ""

	if rows.Next() {

		err = commonlib.PutRecord(rows, &contractId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if contractId == "" {
		contractId = "0"
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertSql := "insert into sign_in(child_id,sign_time,schedule_detail_id,type,contract_id,employee_id) values(?,?,?,?,?,?)"

	lessgo.Log.Debug(insertSql)

	stmt, err := tx.Prepare(insertSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(ids, time.Now().Format("20060102150405"), scheduleId, 2, contractId, employee.UserId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

//报班
func AddChildToClassAction(w http.ResponseWriter, r *http.Request) {

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

	classId := r.FormValue("classId")
	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	idList := strings.Split(ids, ",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	for _, id := range idList {

		checkExistSql := "select count(1) from schedule_detail_child where child_id=? and wyclass_id=? and schedule_detail_id=? "

		rows, err := db.Query(checkExistSql, id, classId, scheduleId)

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

		if totalNum > 0 {
			continue
		}

		insertToClassSql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status,wyclass_id) values(?,?,?,?,?,?) "
		lessgo.Log.Debug(insertToClassSql)
		stmt, err := tx.Prepare(insertToClassSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(scheduleId, id, time.Now().Format("20060102150405"), employee.UserId, 1, classId)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		getConsumerIdSql := "select cons.id from consumer_new cons left join parent p on p.pid=cons.parent_id left join child ch on ch.pid=p.pid where ch.cid=?"
		lessgo.Log.Debug(getConsumerIdSql)

		rows, err = db.Query(getConsumerIdSql, id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		consumerId := 0

		if rows.Next() {
			err = commonlib.PutRecord(rows, &consumerId)
			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if consumerId != 0 {
			selectContatusSql := "select contact_status from consumer_new where id=? "
			lessgo.Log.Debug(selectContatusSql)

			rows, err = db.Query(selectContatusSql, consumerId)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			oldStatus := ""

			if rows.Next() {
				err = commonlib.PutRecord(rows, &oldStatus)
				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门"
					commonlib.OutputJson(w, m, " ")
					return
				}
			}

			updateConsumerStatusSql := "update consumer_new set contact_status=? where id=? "
			lessgo.Log.Debug(updateConsumerStatusSql)

			stmt, err = tx.Prepare(updateConsumerStatusSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(CONSUMER_STATUS_NO_SIGNIN, consumerId)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			insertLogSql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"
			lessgo.Log.Debug(insertLogSql)
			stmt, err = tx.Prepare(insertLogSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(consumerId, employee.UserId, time.Now().Format("20060102150405"), oldStatus, CONSUMER_STATUS_NO_SIGNIN)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											select ?,?,concat('邀约至',wc.start_time,wc.name,'中'),?,3 from wyclass wc where wc.class_id=?`
			lessgo.Log.Debug(insertConsumerContactLogSql)

			stmt, err = tx.Prepare(insertConsumerContactLogSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), consumerId, classId)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

		}

	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

//报班
func AddChildToClassQuickAction(w http.ResponseWriter, r *http.Request) {

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

	consumerId := r.FormValue("consumerId")
	scheduleId := r.FormValue("scheduleId")
	classId := r.FormValue("classId")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	getChildIdSql := "select ch.cid from consumer_new cons left join parent p on p.pid=cons.parent_id left join child ch on ch.pid=p.pid where cons.id=?"
	lessgo.Log.Debug(getChildIdSql)

	childId := 0

	rows, err := db.Query(getChildIdSql, consumerId)
	if rows.Next() {
		err := rows.Scan(&childId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if childId == 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该客户没有关联小朋友无法排课"
		commonlib.OutputJson(w, m, " ")
		return
	}

	checkExistSql := "select count(1) from schedule_detail_child where child_id=? and schedule_detail_id=? "

	rows, err = db.Query(checkExistSql, childId, scheduleId)

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

	if totalNum > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该学员已经在班级中，无需重复分配"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertToClassSql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status,wyclass_id) values(?,?,?,?,?,?) "
	lessgo.Log.Debug(insertToClassSql)
	stmt, err := tx.Prepare(insertToClassSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employee.UserId, 1, classId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	selectContatusSql := "select contact_status from consumer_new where id=? "
	lessgo.Log.Debug(selectContatusSql)

	rows, err = db.Query(selectContatusSql, consumerId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	oldStatus := ""

	if rows.Next() {
		err = commonlib.PutRecord(rows, &oldStatus)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	updateConsumerStatusSql := "update consumer_new set contact_status=? where id=? "
	lessgo.Log.Debug(updateConsumerStatusSql)

	stmt, err = tx.Prepare(updateConsumerStatusSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(CONSUMER_STATUS_NO_SIGNIN, consumerId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertLogSql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"
	lessgo.Log.Debug(insertLogSql)
	stmt, err = tx.Prepare(insertLogSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(consumerId, employee.UserId, time.Now().Format("20060102150405"), oldStatus, CONSUMER_STATUS_NO_SIGNIN)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											select ?,?,concat('邀约至',wc.start_time,wc.name,'中'),?,3 from wyclass wc where wc.class_id=?`
	lessgo.Log.Debug(insertConsumerContactLogSql)

	stmt, err = tx.Prepare(insertConsumerContactLogSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), consumerId, classId)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func WyClassChangeClassAction(w http.ResponseWriter, r *http.Request) {

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
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("childId")
	scheduleId := r.FormValue("scheduleId")
	oldClassId := r.FormValue("oldClassId")
	newClassId := r.FormValue("newClassId")

	db := lessgo.GetMySQL()
	defer db.Close()

	checkIsSignSql := "select count(1) from sign_in where child_id=? and wyclass_id=?"
	lessgo.Log.Debug(checkIsSignSql)

	rows, err := db.Query(checkIsSignSql, id,oldClassId)
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	isSign := 0

	if rows.Next() {
		err = commonlib.PutRecord(rows, &isSign)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if isSign > 0  {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "已签到的学生无法调班"
		commonlib.OutputJson(w, m, " ")
		return
	}

	checkExistSql := "select count(1) from schedule_detail_child where wyclass_id=? and schedule_detail_id=? and child_id=? "
	lessgo.Log.Debug(checkExistSql)

	rows, err = db.Query(checkExistSql, newClassId,scheduleId, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	totalNum := 0

	if rows.Next() {
		err = commonlib.PutRecord(rows, &totalNum)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if totalNum > 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该学生已经在此班级中，无需重复排班"
		commonlib.OutputJson(w, m, " ")
		return
	}

	getInviteEmployeeIdSql := "select create_user from schedule_detail_child where wyclass_id=? and child_id=? "
	lessgo.Log.Debug(getInviteEmployeeIdSql)

	rows, err = db.Query(getInviteEmployeeIdSql, oldClassId,id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	inviteEmployeeId := ""

	if rows.Next() {
		err = commonlib.PutRecord(rows, &inviteEmployeeId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	deleteOldClassSql := "delete from schedule_detail_child where wyclass_id=? and child_id=? "
	lessgo.Log.Debug(deleteOldClassSql)

	stmt, err := tx.Prepare(deleteOldClassSql)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(oldClassId, id)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertNewClassSql := "insert into schedule_detail_child(create_time,create_user,child_id,wyclass_id,schedule_detail_id) values(?,?,?,?,?)"
	lessgo.Log.Debug(insertNewClassSql)
	stmt, err = tx.Prepare(insertNewClassSql)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(time.Now().Format("20060102150405"), inviteEmployeeId, id, newClassId,scheduleId)

	if err != nil {
		tx.Rollback()

		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	insertConsumerContactLogSql := `insert into consumer_contact_log(create_user,create_time,note,consumer_id,type)
											select ?,?,concat('从',wc1.start_time,wc1.name,'中调至',wc2.start_time,wc2.name,'中'),cons.id,3
											from wyclass wc1,wyclass wc2,consumer_new cons,child ch
											where wc1.class_id=? and wc2.class_id=? and ch.cid=? and ch.pid=cons.parent_id`
	lessgo.Log.Debug(insertConsumerContactLogSql)

	stmt, err = tx.Prepare(insertConsumerContactLogSql)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	_, err = stmt.Exec(employee.UserId, time.Now().Format("20060102150405"), oldClassId,newClassId,id)

	if err != nil {
		tx.Rollback()
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}


	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func AddChildForNormalOnceAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	ids := r.FormValue("ids")

	idList := strings.Split(ids, ",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	for _, id := range idList {

		checkExistSql := "select count(1) from schedule_detail_child where child_id=? and  schedule_detail_id=? "

		rows, err := db.Query(checkExistSql, id, scheduleId)

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

		if totalNum > 0 {
			continue
		}

		insertToClassSql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user) values(?,?,?,?) "
		lessgo.Log.Debug(insertToClassSql)
		stmt, err := tx.Prepare(insertToClassSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(scheduleId, id, time.Now().Format("20060102150405"), employee.UserId)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

	}

	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}

func AddChildForNormalTempelateAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	id := r.FormValue("ids")

	if strings.Contains(id, ",") {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "只能操作一个学生"
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	getScheduleTempId := "select st.id,st.room_id,st.time_id,st.week,csd.start_time from class_schedule_detail csd left join schedule_template st on csd.center_id=st.center_id and csd.room_id=st.room_id and csd.time_id=st.time_id and csd.week=st.week where csd.id=? "
	lessgo.Log.Debug(getScheduleTempId)
	rows, err := db.Query(getScheduleTempId, scheduleId)

	scheduleTempId := 0
	stRoomId := 0
	stTimeId := 0
	stWeek := 0
	scdStartTime := ""

	if rows.Next() {
		err := commonlib.PutRecord(rows, &scheduleTempId,&stRoomId,&stTimeId,&stWeek,&scdStartTime)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if scheduleTempId == 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "该课表不是模板课表，无法跟班"
		commonlib.OutputJson(w, m, " ")
		return
	}

	getFurtherScheduleSql := "select id from class_schedule_detail where time_id=? and room_id=? and week=? and start_time>=? "
	lessgo.Log.Debug(getFurtherScheduleSql)
	scheduleRows, err := db.Query(getFurtherScheduleSql, stTimeId,stRoomId,stWeek,scdStartTime)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for scheduleRows.Next() {

		scheduleId := 0

		err := commonlib.PutRecord(scheduleRows, &scheduleId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		checkExistSql := "select count(1) from schedule_detail_child where child_id=? and schedule_detail_id=? "
		lessgo.Log.Debug(checkExistSql)
		rows, err = db.Query(checkExistSql, id, scheduleId)

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

		if totalNum == 0 {
			insertToClassSql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user) values(?,?,?,?) "
			lessgo.Log.Debug(insertToClassSql)
			stmt, err := tx.Prepare(insertToClassSql)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(scheduleId, id, time.Now().Format("20060102150405"), employee.UserId)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}
	}

	checkTempExistSql := "select count(1) from schedule_template_child where child_id=? and  schedule_template_id=? "
	lessgo.Log.Debug(checkTempExistSql)
	rows, err = db.Query(checkTempExistSql, id, scheduleTempId)

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

	if totalNum == 0 {
		insertToTempSql := "insert into schedule_template_child(schedule_template_id,child_id) values(?,?) "
		lessgo.Log.Debug(insertToTempSql)

		stmt, err := tx.Prepare(insertToTempSql)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(scheduleTempId, id)

		if err != nil {
			tx.Rollback()
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}
	tx.Commit()

	m["success"] = true
	m["code"] = 200
	m["msg"] = "操作成功"
	commonlib.OutputJson(w, m, " ")
	return
}
