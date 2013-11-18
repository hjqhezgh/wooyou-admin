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
		if roleCode == "admin" || roleCode == "yyzj" || roleCode == "zjl" || roleCode == "yyzy" || roleCode == "tmk"{
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

	sql := "select  c.class_id,c.name,ce.name as cename,cou.name as courseName,c.end_time,c.deadline,c.max_child_num,tea.really_name as teacherName,ass.really_name as assName,c.center_id from wyclass c left join center ce on ce.cid=c.center_id left join employee tea on tea.user_id=c.teacher_id left join employee ass on ass.user_id=c.assistant_id left join course cou on c.course_id=cou.cid where 1=1 "

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

	sql := "select class_id,name,center_id,assistant_id,teacher_id,max_child_num,deadline,end_time,course_id from wyclass where class_id=? "

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

	var classId, name, centerId, assistantId, teacherId, maxChildNum, deadline, endTime, courseId string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &classId, &name, &centerId, &assistantId, &teacherId, &maxChildNum, &deadline, &endTime, &courseId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	m["success"] = true

	loadFormObjects := []lessgo.LoadFormObject{}

	h1 := lessgo.LoadFormObject{"name", name}
	h2 := lessgo.LoadFormObject{"class_id", classId}
	h3 := lessgo.LoadFormObject{"center_id", centerId}
	h4 := lessgo.LoadFormObject{"course_id", centerId + "," + courseId}
	h5 := lessgo.LoadFormObject{"end_time", endTime}
	h6 := lessgo.LoadFormObject{"deadline", deadline}
	h7 := lessgo.LoadFormObject{"max_child_num", maxChildNum}
	h8 := lessgo.LoadFormObject{"teacher_id", centerId + "," + teacherId}
	h9 := lessgo.LoadFormObject{"assistant_id", centerId + "," + assistantId}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)
	loadFormObjects = append(loadFormObjects, h8)
	loadFormObjects = append(loadFormObjects, h9)

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}

func WyClassUpdateAction(w http.ResponseWriter, r *http.Request) {

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

	class_id := r.FormValue("class_id")
	name := r.FormValue("name")
	center_id := r.FormValue("center_id")
	course_id := r.FormValue("course_id")
	end_time := r.FormValue("end_time")
	deadline := r.FormValue("deadline")
	max_child_num := r.FormValue("max_child_num")
	teacher_id := r.FormValue("teacher_id")
	assistant_id := r.FormValue("assistant_id")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "update wyclass set name=?,center_id=?,course_id=?,end_time=?,deadline=?,max_child_num=?,teacher_id=?,assistant_id=? where class_id=? "

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

	_, err = stmt.Exec(name, center_id, course_id, end_time, deadline, max_child_num, teacher_id, assistant_id, class_id)

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

}

func WyClassInsertAction(w http.ResponseWriter, r *http.Request) {

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

	name := r.FormValue("name")
	center_id := r.FormValue("center_id")
	course_id := r.FormValue("course_id")
	end_time := r.FormValue("end_time")
	deadline := r.FormValue("deadline")
	max_child_num := r.FormValue("max_child_num")
	child_num := r.FormValue("child_num")
	teacher_id := r.FormValue("teacher_id")
	assistant_id := r.FormValue("assistant_id")

	db := lessgo.GetMySQL()
	defer db.Close()

	checkClassExistSql := "select count(1) from wyclass where center_id=? and name=? "
	lessgo.Log.Debug(checkClassExistSql)

	rows, err := db.Query(checkClassExistSql, center_id,name)

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
		m["msg"] = "班级名重复，请更换一个班级名重复"
		commonlib.OutputJson(w, m, " ")
		return
	}


	sql := "insert into wyclass(assistant_id,name,create_time,course_id,center_id,child_num,end_time,deadline,max_child_num,teacher_id) values(?,?,?,?,?,?,?,?,?,?)"

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

	_, err = stmt.Exec(assistant_id,name,time.Now().Format("20060102150405"), course_id, center_id,child_num, end_time, deadline, max_child_num, teacher_id)

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
