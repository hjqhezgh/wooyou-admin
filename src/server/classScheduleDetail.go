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
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RoomOrTime struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Weekdate struct {
	Date string `json:"date"`
	Week string `json:"week"`
}

type Schedule struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	PersonNum           int    `json:"personNum"`
	CurrentTMKPersonNum int    `json:"currentTMKPersonNum"`
	SignNum             int    `json:"signNum"`
	Teacher             string `json:"teacher"`
	Assistant           string `json:"assistant"`
	IsNormal            int    `json:"isNormal"`
	RoomId              int    `json:"roomId"`
	TimeId              int    `json:"timeId"`
	Week                int    `json:"week"`
	ClassId             int    `json:"classId"`
	CenterId            int    `json:"centerId"`
	Code                string `json:"code"`
}

/*
课表查询sql
select csd.id,csd.room_id,csd.time_id,csd.course_id,teacher.really_name teaName,assistant.really_name assName,personNum.num perNum,signNum.num sigNum,wc.name,cour.name from class_schedule_detail csd
left join employee teacher on teacher.user_id=csd.teacher_id
left join employee assistant on assistant.user_id=csd.assistant_id
left join (select count(1) num,schedule_detail_id from schedule_detail_child group by schedule_detail_id) personNum on personNum.schedule_detail_id = csd.id
left join time_section ts on ts.id=csd.time_id
left join (select count(1) num,schedule_detail_id from sign_in group by schedule_detail_id) signNum on signNum.schedule_detail_id = csd.id
left join wyclass wc on wc.class_id=csd.class_id
left join course cour on cour.cid=csd.course_id
*/
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

	err := r.ParseForm()

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	date := r.FormValue("date")
	dateType := r.FormValue("type")

	searchDate := time.Now()

	if date != "" {
		searchDate, _ = time.ParseInLocation("20060102150405", date, time.Local)
	}

	if dateType == "pre" {
		searchDate = searchDate.Add(time.Duration(-7*24) * time.Hour)
	} else if dateType == "next" {
		searchDate = searchDate.Add(time.Duration(7*24) * time.Hour)
	}

	week := 0

	if searchDate.Weekday() == time.Monday {
		week = 1
	} else if searchDate.Weekday() == time.Tuesday {
		week = 2
	} else if searchDate.Weekday() == time.Wednesday {
		week = 3
	} else if searchDate.Weekday() == time.Thursday {
		week = 4
	} else if searchDate.Weekday() == time.Friday {
		week = 5
	} else if searchDate.Weekday() == time.Saturday {
		week = 6
	} else if searchDate.Weekday() == time.Sunday {
		week = 7
	}

	monday := searchDate.Add(-1 * time.Duration(week-1) * 24 * time.Hour)
	sunday := searchDate.Add(time.Duration(7-week) * 24 * time.Hour)
	st := monday.Format("20060102") + "000000"
	et := sunday.Format("20060102") + "235959"

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()
	rooms := []RoomOrTime{}
	times := []RoomOrTime{}
	weekdates := []Weekdate{}
	schedules := []Schedule{}

	getRoomInfoSql := "select rid,name from room where cid=? "
	lessgo.Log.Debug(getRoomInfoSql)

	rows, err := db.Query(getRoomInfoSql, _employee.CenterId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		room := RoomOrTime{}
		err = commonlib.PutRecord(rows, &room.Id, &room.Name)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		rooms = append(rooms, room)
	}

	getTimeSql := "select id,start_time,end_time from time_section where center_id=? order by lesson_no"
	lessgo.Log.Debug(getRoomInfoSql)

	rows, err = db.Query(getTimeSql, _employee.CenterId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		timeObject := RoomOrTime{}
		startTime := ""
		endTime := ""
		err = commonlib.PutRecord(rows, &timeObject.Id, &startTime, &endTime)
		timeObject.Name = startTime + "~" + endTime

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		times = append(times, timeObject)
	}

	for i := 0; i < 7; i++ {
		var weekdate = Weekdate{}
		if i == 0 {
			weekdate.Week = "星期一"
		} else if i == 1 {
			weekdate.Week = "星期二"
		} else if i == 2 {
			weekdate.Week = "星期三"
		} else if i == 3 {
			weekdate.Week = "星期四"
		} else if i == 4 {
			weekdate.Week = "星期五"
		} else if i == 5 {
			weekdate.Week = "星期六"
		} else if i == 6 {
			weekdate.Week = "星期日"
		}

		theday := monday.Add(time.Duration(i*24) * time.Hour)
		weekdate.Date = theday.Format("2006-01-02")

		weekdates = append(weekdates, weekdate)
	}

	getScheduleInfoSql := "select csd.id,csd.room_id,csd.time_id,csd.course_id,teacher.really_name teaName,assistant.really_name assName,personNum.num perNum,signNum.num sigNum,wc.name,cour.name,csd.week,csd.class_id,csd.center_id,wc.code from class_schedule_detail csd "
	getScheduleInfoSql += " left join employee teacher on teacher.user_id=csd.teacher_id  "
	getScheduleInfoSql += " left join employee assistant on assistant.user_id=csd.assistant_id "
	getScheduleInfoSql += " left join (select count(1) num,schedule_detail_id from schedule_detail_child group by schedule_detail_id) personNum on personNum.schedule_detail_id = csd.id "
	getScheduleInfoSql += " left join time_section ts on ts.id=csd.time_id "
	getScheduleInfoSql += " left join (select count(1) num,schedule_detail_id from sign_in group by schedule_detail_id) signNum on signNum.schedule_detail_id = csd.id "
	getScheduleInfoSql += " left join wyclass wc on wc.class_id=csd.class_id "
	getScheduleInfoSql += " left join course cour on cour.cid=csd.course_id "
	getScheduleInfoSql += " where csd.start_time>=? and csd.start_time<=? and csd.center_id=? "
	lessgo.Log.Debug(getScheduleInfoSql)

	rows, err = db.Query(getScheduleInfoSql, st, et, _employee.CenterId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		schedule := Schedule{}
		courseId := 0
		className := ""
		courseName := ""

		err = commonlib.PutRecord(rows, &schedule.Id, &schedule.RoomId, &schedule.TimeId, &courseId, &schedule.Teacher, &schedule.Assistant, &schedule.PersonNum, &schedule.SignNum, &className, &courseName, &schedule.Week, &schedule.ClassId, &schedule.CenterId, &schedule.Code)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if courseId != 0 {
			schedule.IsNormal = 1
			schedule.Name = courseName
		} else {
			schedule.IsNormal = 2
			schedule.Name = className
		}

		schedules = append(schedules, schedule)
	}

	m["success"] = true
	m["code"] = 200
	m["msg"] = "成功"
	m["times"] = times
	m["rooms"] = rooms
	m["weekdates"] = weekdates
	m["schedules"] = schedules
	m["firstDayOfWeek"] = st

	commonlib.OutputJson(w, m, "")
}

func ClassScheduleDetailListQuickAction(w http.ResponseWriter, r *http.Request) {

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

	date := r.FormValue("date")
	centerId := r.FormValue("centerId")
	dateType := r.FormValue("type")

	searchDate := time.Now()

	if date != "" {
		searchDate, _ = time.ParseInLocation("20060102150405", date, time.Local)
	}

	if dateType == "pre" {
		searchDate = searchDate.Add(time.Duration(-7*24) * time.Hour)
	} else if dateType == "next" {
		searchDate = searchDate.Add(time.Duration(7*24) * time.Hour)
	}

	week := 0

	if searchDate.Weekday() == time.Monday {
		week = 1
	} else if searchDate.Weekday() == time.Tuesday {
		week = 2
	} else if searchDate.Weekday() == time.Wednesday {
		week = 3
	} else if searchDate.Weekday() == time.Thursday {
		week = 4
	} else if searchDate.Weekday() == time.Friday {
		week = 5
	} else if searchDate.Weekday() == time.Saturday {
		week = 6
	} else if searchDate.Weekday() == time.Sunday {
		week = 7
	}

	monday := searchDate.Add(-1 * time.Duration(week-1) * 24 * time.Hour)
	sunday := searchDate.Add(time.Duration(7-week) * 24 * time.Hour)
	st := monday.Format("20060102") + "000000"
	et := sunday.Format("20060102") + "235959"

	db := lessgo.GetMySQL()
	defer db.Close()
	rooms := []RoomOrTime{}
	times := []RoomOrTime{}
	weekdates := []Weekdate{}
	schedules := []Schedule{}

	getRoomInfoSql := "select rid,name from room where cid=? "
	lessgo.Log.Debug(getRoomInfoSql)

	rows, err := db.Query(getRoomInfoSql, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		room := RoomOrTime{}
		err = commonlib.PutRecord(rows, &room.Id, &room.Name)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		rooms = append(rooms, room)
	}

	getTimeSql := "select id,start_time,end_time from time_section where center_id=? order by lesson_no"
	lessgo.Log.Debug(getRoomInfoSql)

	rows, err = db.Query(getTimeSql, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		timeObject := RoomOrTime{}
		startTime := ""
		endTime := ""
		err = commonlib.PutRecord(rows, &timeObject.Id, &startTime, &endTime)
		timeObject.Name = startTime + "~" + endTime

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		times = append(times, timeObject)
	}

	for i := 0; i < 7; i++ {
		var weekdate = Weekdate{}
		if i == 0 {
			weekdate.Week = "星期一"
		} else if i == 1 {
			weekdate.Week = "星期二"
		} else if i == 2 {
			weekdate.Week = "星期三"
		} else if i == 3 {
			weekdate.Week = "星期四"
		} else if i == 4 {
			weekdate.Week = "星期五"
		} else if i == 5 {
			weekdate.Week = "星期六"
		} else if i == 6 {
			weekdate.Week = "星期日"
		}

		theday := monday.Add(time.Duration(i*24) * time.Hour)
		weekdate.Date = theday.Format("2006-01-02")

		weekdates = append(weekdates, weekdate)
	}

	getScheduleInfoSql := "select csd.id,csd.room_id,csd.time_id,csd.course_id,teacher.really_name teaName,assistant.really_name assName,personNum.num perNum,signNum.num sigNum,wc.name,cour.name,csd.week,csd.class_id,csd.center_id,wc.code,currTMKPersonNum.num from class_schedule_detail csd "
	getScheduleInfoSql += " left join employee teacher on teacher.user_id=csd.teacher_id  "
	getScheduleInfoSql += " left join employee assistant on assistant.user_id=csd.assistant_id "
	getScheduleInfoSql += " left join (select count(1) num,schedule_detail_id from schedule_detail_child group by schedule_detail_id) personNum on personNum.schedule_detail_id = csd.id "
	getScheduleInfoSql += " left join (select count(1) num,schedule_detail_id from schedule_detail_child where create_user=? group by schedule_detail_id) currTMKPersonNum on currTMKPersonNum.schedule_detail_id = csd.id "
	getScheduleInfoSql += " left join time_section ts on ts.id=csd.time_id "
	getScheduleInfoSql += " left join (select count(1) num,schedule_detail_id from sign_in group by schedule_detail_id) signNum on signNum.schedule_detail_id = csd.id "
	getScheduleInfoSql += " left join wyclass wc on wc.class_id=csd.class_id "
	getScheduleInfoSql += " left join course cour on cour.cid=csd.course_id "
	getScheduleInfoSql += " where csd.start_time>=? and csd.start_time<=? and csd.center_id=? and csd.course_id is null"
	lessgo.Log.Debug(getScheduleInfoSql)

	rows, err = db.Query(getScheduleInfoSql, employee.UserId,st, et, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		schedule := Schedule{}
		courseId := 0
		className := ""
		courseName := ""

		err = commonlib.PutRecord(rows, &schedule.Id, &schedule.RoomId, &schedule.TimeId, &courseId, &schedule.Teacher, &schedule.Assistant, &schedule.PersonNum, &schedule.SignNum, &className, &courseName, &schedule.Week, &schedule.ClassId, &schedule.CenterId, &schedule.Code,&schedule.CurrentTMKPersonNum)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		if courseId != 0 {
			schedule.IsNormal = 1
			schedule.Name = courseName
		} else {
			schedule.IsNormal = 2
			schedule.Name = className
		}

		schedules = append(schedules, schedule)
	}

	m["success"] = true
	m["code"] = 200
	m["msg"] = "成功"
	m["times"] = times
	m["rooms"] = rooms
	m["weekdates"] = weekdates
	m["schedules"] = schedules
	m["firstDayOfWeek"] = st

	commonlib.OutputJson(w, m, "")
}

func ClassScheduleDetailAddAction(w http.ResponseWriter, r *http.Request) {

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

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	roomId := r.FormValue("roomId")
	timeId := r.FormValue("timeId")
	courseId := r.FormValue("courseId")
	teacherId := r.FormValue("teacherId")
	assistantId := r.FormValue("assistantId")
	week := r.FormValue("week")
	date := r.FormValue("date")
	scheduleType := r.FormValue("type")
	isTmp := r.FormValue("isTmp")
	capacity := r.FormValue("capacity")
	name := r.FormValue("name")
	code := r.FormValue("code")

	date = strings.Replace(date, "-", "", -1)

	timeSection, err := FindTimeSectionById(timeId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	startTime := date + strings.Replace(timeSection.StartTime, ":", "", -1) + "00"
	endTime := date + strings.Replace(timeSection.EndTime, ":", "", -1) + "00"

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

	if scheduleType == "normal" {

		insertScheduleSql := "insert into class_schedule_detail(teacher_id,assistant_id,course_id,center_id,time_id,room_id,day_date,week,start_time,end_time,status,capacity) values(?,?,?,?,?,?,?,?,?,?,?,?)"
		lessgo.Log.Debug(insertScheduleSql)

		stmt, err := tx.Prepare(insertScheduleSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(teacherId, assistantId, courseId, _employee.CenterId, timeId, roomId, date, week, startTime, endTime, 1, capacity)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		if isTmp == "1" {
			selectTmpSql := "select id from schedule_template where center_id=? and room_id=? and time_id=? "
			lessgo.Log.Debug(selectTmpSql)

			rows, err := db.Query(selectTmpSql, _employee.CenterId, roomId, timeId)

			scheduleTemplateId := 0

			if rows.Next() {
				err := commonlib.PutRecord(rows, &scheduleTemplateId)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门"
					commonlib.OutputJson(w, m, " ")
					return
				}
			}

			if scheduleTemplateId == 0 {
				insertScheduleTmpSql := "insert into schedule_template(center_id,room_id,teacher_id,assistant_id,time_id,week,course_id) values(?,?,?,?,?,?,?)"
				lessgo.Log.Debug(insertScheduleTmpSql)

				stmt, err = tx.Prepare(insertScheduleTmpSql)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = stmt.Exec(_employee.CenterId, roomId, teacherId, assistantId, timeId, week, courseId)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}
			} else {
				updateScheduleTmpSql := "update schedule_template set teacher_id=?, assistant_id=?,course_id=? where id=? "
				lessgo.Log.Debug(updateScheduleTmpSql)

				stmt, err = tx.Prepare(updateScheduleTmpSql)
				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}
				_, err = stmt.Exec(teacherId, assistantId, courseId, scheduleTemplateId)

				if err != nil {
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
		commonlib.OutputJson(w, m, " ")
	} else {
		insertClassSql := "insert into wyclass(name,create_time,center_id,start_time,code) values(?,?,?,?,?)"
		lessgo.Log.Debug(insertClassSql)

		stmt, err := tx.Prepare(insertClassSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		res, err := stmt.Exec(name, time.Now().Format("20060102150405"), _employee.CenterId, startTime, code)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		classId, err := res.LastInsertId()

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		insertScheduleSql := "insert into class_schedule_detail(class_id,center_id,time_id,room_id,day_date,week,start_time,end_time,status,capacity) values(?,?,?,?,?,?,?,?,?,?)"
		lessgo.Log.Debug(insertScheduleSql)

		stmt, err = tx.Prepare(insertScheduleSql)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(classId, _employee.CenterId, timeId, roomId, date, week, startTime, endTime, 1, capacity)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		tx.Commit()

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	}
}

func ClassScheduleDetailModifyAction(w http.ResponseWriter, r *http.Request) {

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

	scheduleId := r.FormValue("scheduleId")
	name := r.FormValue("name")
	code := r.FormValue("code")
	teacherId := r.FormValue("teacherId")
	assistantId := r.FormValue("assistantId")
	scheduleType := r.FormValue("type")
	classId := r.FormValue("classId")

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

	if scheduleType == "normal" {

		updateSql := "update class_schedule_detail set teacher_id=?,assistant_id=? where id=?"
		lessgo.Log.Debug(updateSql)

		stmt, err := tx.Prepare(updateSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(teacherId, assistantId, scheduleId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		tx.Commit()

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	} else {
		updateSql := "update wyclass set name=? , code=? where class_id=?"
		lessgo.Log.Debug(updateSql)

		stmt, err := tx.Prepare(updateSql)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = stmt.Exec(name, code, classId)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		tx.Commit()

		m["success"] = true
		commonlib.OutputJson(w, m, " ")
	}
}

//课表读取
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

	id := r.FormValue("scheduleId")
	roomId := r.FormValue("roomId")
	timeId := r.FormValue("timeId")
	week := r.FormValue("week")
	date := r.FormValue("date")
	//	schedule := r.FormValue("type")

	db := lessgo.GetMySQL()
	defer db.Close()

	loadFormObjects := []lessgo.LoadFormObject{}

	if id != "" {
		getScheduleDetailSql := "select csd.teacher_id,csd.assistant_id,wc.name,wc.code,wc.class_id from class_schedule_detail csd left join wyclass wc on wc.class_id=csd.class_id where csd.id=?"
		lessgo.Log.Debug(getScheduleDetailSql)

		rows, err := db.Query(getScheduleDetailSql, id)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		var teacherId, assistantId, className, classCode, classId string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &teacherId, &assistantId, &className, &classCode, &classId)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		m["success"] = true

		h1 := lessgo.LoadFormObject{"scheduleId", id}
		h2 := lessgo.LoadFormObject{"name", className}
		h3 := lessgo.LoadFormObject{"code", classCode}
		h4 := lessgo.LoadFormObject{"teacherId", teacherId}
		h5 := lessgo.LoadFormObject{"assistantId", assistantId}
		h6 := lessgo.LoadFormObject{"classId", classId}

		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
		loadFormObjects = append(loadFormObjects, h3)
		loadFormObjects = append(loadFormObjects, h4)
		loadFormObjects = append(loadFormObjects, h5)
		loadFormObjects = append(loadFormObjects, h6)
	} else {
		roomInfoSql := "select name from room where rid=? "
		lessgo.Log.Debug(roomInfoSql)

		rows, err := db.Query(roomInfoSql, roomId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		var roomName string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &roomName)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		timeInfoSql := "select start_time,end_time from time_section where id=?"
		lessgo.Log.Debug(timeInfoSql)

		rows, err = db.Query(timeInfoSql, timeId)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		var startTime, endTime string

		if rows.Next() {
			err = commonlib.PutRecord(rows, &startTime, &endTime)

			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		weekDesc := ""
		if week == "1" {
			weekDesc = "星期一"
		} else if week == "2" {
			weekDesc = "星期二"
		} else if week == "3" {
			weekDesc = "星期三"
		} else if week == "4" {
			weekDesc = "星期四"
		} else if week == "5" {
			weekDesc = "星期五"
		} else if week == "6" {
			weekDesc = "星期六"
		} else if week == "7" {
			weekDesc = "星期日"
		}

		m["success"] = true

		h1 := lessgo.LoadFormObject{"id", id}
		h2 := lessgo.LoadFormObject{"roomId", roomId}
		h3 := lessgo.LoadFormObject{"timeId", timeId}
		h4 := lessgo.LoadFormObject{"roomName", roomName}
		h5 := lessgo.LoadFormObject{"timeDesc", startTime + "~" + endTime}
		h6 := lessgo.LoadFormObject{"weekDesc", date + "(" + weekDesc + ")"}
		h7 := lessgo.LoadFormObject{"week", week}
		h8 := lessgo.LoadFormObject{"date", date}
		h9 := lessgo.LoadFormObject{"capacity", "10"}

		loadFormObjects = append(loadFormObjects, h1)
		loadFormObjects = append(loadFormObjects, h2)
		loadFormObjects = append(loadFormObjects, h3)
		loadFormObjects = append(loadFormObjects, h4)
		loadFormObjects = append(loadFormObjects, h5)
		loadFormObjects = append(loadFormObjects, h6)
		loadFormObjects = append(loadFormObjects, h7)
		loadFormObjects = append(loadFormObjects, h8)
		loadFormObjects = append(loadFormObjects, h9)
	}

	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
}

func CreateWeekScheduleAction(w http.ResponseWriter, r *http.Request) {

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

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	err = r.ParseForm()

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	firstDayOfWeek := r.FormValue("firstDayOfWeek")

	db := lessgo.GetMySQL()
	defer db.Close()

	getScheduleTmpSql := "select id,room_id,teacher_id,assistant_id,time_id,week,course_id from schedule_template where center_id=? "
	lessgo.Log.Debug(getScheduleTmpSql)

	rows, err := db.Query(getScheduleTmpSql, _employee.CenterId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	firstDay, _ := time.ParseInLocation("20060102150405", firstDayOfWeek, time.Local)

	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	for rows.Next() {
		var scheduleTmpId, roomId, teacherId, assistantId, timeId, courseId string
		week := 0

		err = commonlib.PutRecord(rows, &scheduleTmpId, &roomId, &teacherId, &assistantId, &timeId, &week, &courseId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		theDay := firstDay.Add(time.Duration(week-1) * 24)
		date := theDay.Format("20060102")

		getScheduleDetailSql := "select id from class_schedule_detail where center_id=? and room_id=? and time_id=? and day_date=?"
		lessgo.Log.Debug(getScheduleDetailSql)

		scheduleDetailRows, err := db.Query(getScheduleDetailSql, _employee.CenterId, roomId, timeId, date)
		if err != nil {
			lessgo.Log.Error(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		scheduleDetailId := 0

		if scheduleDetailRows.Next() {
			err = commonlib.PutRecord(scheduleDetailRows, &scheduleDetailId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}
		}

		if scheduleDetailId == 0 {

			insertScheduleSql := "insert into class_schedule_detail(teacher_id,assistant_id,course_id,center_id,time_id,room_id,day_date,week,start_time,end_time,status,capacity) values(?,?,?,?,?,?,?,?,?,?,?,?)"
			lessgo.Log.Debug(insertScheduleSql)

			stmt, err := tx.Prepare(insertScheduleSql)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			timeSection, err := FindTimeSectionById(timeId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			startTime := date + strings.Replace(timeSection.StartTime, ":", "", -1) + "00"
			endTime := date + strings.Replace(timeSection.EndTime, ":", "", -1) + "00"

			res, err := stmt.Exec(teacherId, assistantId, courseId, _employee.CenterId, timeId, roomId, date, week, startTime, endTime, 1, 10)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			newScheduleDetailId, err := res.LastInsertId()

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			getChildSql := "select child_id,contract_id from schedule_template_child where schedule_template_id=? "
			lessgo.Log.Debug(getChildSql)

			childRows, err := db.Query(getChildSql, scheduleTmpId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			for childRows.Next() {
				var childId string
				var contractId int
				err = commonlib.PutRecord(childRows, &childId, &contractId)

				if err != nil {
					lessgo.Log.Error(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				insertScheduleChildSql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,contract_id) values(?,?,?,?,?)"
				lessgo.Log.Debug(insertScheduleChildSql)

				stmt, err = tx.Prepare(insertScheduleChildSql)

				if err != nil {
					lessgo.Log.Error(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = stmt.Exec(newScheduleDetailId, childId, time.Now().Format("20060102150405"), employee.UserId, contractId)

				if err != nil {
					lessgo.Log.Error(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					commonlib.OutputJson(w, m, " ")
					return
				}
			}

		}
	}

	tx.Commit()

	m["success"] = true
	commonlib.OutputJson(w, m, " ")
}

//仅删除课表，不删除模板
func DeleteSingleScheduleAction(w http.ResponseWriter, r *http.Request) {

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

	db := lessgo.GetMySQL()
	defer db.Close()

	checkExistSql := "select count(1) from schedule_detail_child where schedule_detail_id=?"

	lessgo.Log.Debug(checkExistSql)

	rows, err := db.Query(checkExistSql, scheduleId)

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
		m["msg"] = "有人上课，不能删除"
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

	deteleSql := "delete from class_schedule_detail where id=?"

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

	_, err = stmt.Exec(scheduleId)

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
