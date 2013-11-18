// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-12 09:25
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-12 09:25 black 创建文档
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

func ClassScheduleDetailListAction(w http.ResponseWriter, r *http.Request) {

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
	st := r.FormValue("day_date-ge")
	et := r.FormValue("day_date-le")
	courseId := r.FormValue("course_id-eq")

	params := []interface{}{}

	sql := "select csd.id,ce.name as centerName,wy.name as className,teacher.really_name as teacherName,assistant.really_name as assistantName,cour.name as courseName,les.caption,r.name,csd.start_time,csd.end_time,csd.week,csd.capacity,csd.status,wy.class_id,csd.center_id "
	sql += " from class_schedule_detail csd left join center ce on ce.cid=csd.center_id left join employee teacher on teacher.user_id=csd.teacher_id left join employee assistant on assistant.user_id=csd.assistant_id "
	sql += " left join course cour on cour.cid=csd.course_id left join lesson les on les.lid=csd.lesson_id left join room r on r.rid=csd.room_id left join wyclass wy on wy.class_id=csd.class_id where 1=1 "

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
		sql += " and csd.center_id=? "
	}

	if st != "" {
		params = append(params, st)
		sql += " and csd.day_date>=? "
	}

	if et != "" {
		params = append(params, st)
		sql += " and csd.day_date<=? "
	}

	if courseId != "" {
		params = append(params, courseId)
		sql += " and csd.course_id=? "
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

	sql += " order by csd.start_time desc,csd.id desc limit ?,?"

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

		for i := 0; i < 14; i++ {
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

func ClassScheduleDetailSaveAction(w http.ResponseWriter, r *http.Request) {

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
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")
	centerId := r.FormValue("center_id")
	classId := r.FormValue("class_id")
	lessonId := r.FormValue("lesson_id")
	dayDate := r.FormValue("day_date")
	timeId := r.FormValue("time_id")
	roomId := r.FormValue("room_id")
	status := r.FormValue("status")
	teacherId := r.FormValue("teacher_id")
	assistantId := r.FormValue("assistant_id")
	capacity := r.FormValue("capacity")

	db := lessgo.GetMySQL()
	defer db.Close()

	timeSection, err := FindTimeSectionById(timeId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	lesson, err := FindLessonById(lessonId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	startTime := dayDate + strings.Replace(timeSection.StartTime, ":", "", -1)+"00"
	endTime := dayDate + strings.Replace(timeSection.EndTime, ":", "", -1)+"00"

	week := ""

	st, _ := time.ParseInLocation("20060102150405", startTime, time.Local)
	if st.Weekday() == time.Monday {
		week = "1"
	} else if st.Weekday() == time.Tuesday {
		week = "2"
	} else if st.Weekday() == time.Wednesday {
		week = "3"
	} else if st.Weekday() == time.Thursday {
		week = "4"
	} else if st.Weekday() == time.Friday {
		week = "5"
	} else if st.Weekday() == time.Saturday {
		week = "6"
	} else if st.Weekday() == time.Sunday {
		week = "7"
	}

	if id == "" {
		sql := "insert into class_schedule_detail(class_id,teacher_id,assistant_id,course_id,lesson_id,center_id,time_id,room_id,day_date,week,capacity,start_time,end_time,status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

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

		_, err = stmt.Exec(classId, teacherId, assistantId, lesson.CourseId, lessonId, centerId, timeId, roomId, dayDate,week,capacity,startTime,endTime,status)

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
	} else {

			sql := "update class_schedule_detail set class_id=?,teacher_id=?,assistant_id=?,course_id=?,lesson_id=?,center_id=?,time_id=?,room_id=?,day_date=?,week=?,capacity=?,start_time=?,end_time=? where id=? "

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

			_, err = stmt.Exec(classId, teacherId, assistantId, lesson.CourseId, lessonId, centerId, timeId, roomId,dayDate,week,capacity,startTime,endTime,id)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			m["success"] = true
			commonlib.OutputJson(w, m, " ")
	}

}

//客户读取服务
func ClassScheduleDetailLoadAction(w http.ResponseWriter, r *http.Request) {

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

	sql := "select id,class_id,teacher_id,assistant_id,course_id,lesson_id,center_id,time_id,room_id,day_date,capacity,status from class_schedule_detail where id=? "

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

	var classId, teacherId, assistantId, courseId, lessonId, centerId,timeId,roomId,dayDate,capacity,status string

	if rows.Next() {
		err = commonlib.PutRecord(rows, &id,&classId, &teacherId, &assistantId, &courseId, &lessonId, &centerId, &timeId, &roomId, &dayDate, &capacity, &status)

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

	h1 := lessgo.LoadFormObject{"id", id}
	h2 := lessgo.LoadFormObject{"center_id", centerId}
	h3 := lessgo.LoadFormObject{"class_id", centerId+","+classId}
	h4 := lessgo.LoadFormObject{"lesson_id", classId+","+lessonId}
	h5 := lessgo.LoadFormObject{"day_date", dayDate}
	h6 := lessgo.LoadFormObject{"time_id", centerId+","+timeId}
	h7 := lessgo.LoadFormObject{"room_id", centerId+","+roomId}
	h8 := lessgo.LoadFormObject{"status", status}
	h9 := lessgo.LoadFormObject{"teacher_id", centerId+","+teacherId}
	h10 := lessgo.LoadFormObject{"assistant_id", centerId+","+assistantId}
	h11 := lessgo.LoadFormObject{"capacity", capacity}


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
