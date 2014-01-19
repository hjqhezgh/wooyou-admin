// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-15 16:48
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-15 16:48 black 创建文档
package logic

import (
	"database/sql"
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"strconv"
	"strings"
	"time"
)

func AddChildToClass(classId, scheduleId, ids, employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	idList := strings.Split(ids, ",")

	for _, id := range idList {
		childInScheduleFlag, err := checkChildInSchedule(id, scheduleId)

		if err != nil {
			return false, "", err
		}

		if childInScheduleFlag {
			continue
		}

		err = insertScheduleClass(tx, id, scheduleId, classId, employeeId)

		if err != nil {
			return false, "", err
		}

		consumerDataMap, err := getConsumerByChildId(id)

		if err != nil {
			return false, "", err
		}

		consumerId := consumerDataMap["id"]

		if consumerId != "" {
			oldStatus := consumerDataMap["contactStatus"]
			if oldStatus != CONSUMER_STATUS_SIGNIN && oldStatus != CONSUMER_STATUS_NO_SIGNIN { //已经邀约的，或者已经签到的就不需要更改状态了
				updataConsumerMap := make(map[string]interface{})
				updataConsumerMap["contact_status"] = CONSUMER_STATUS_NO_SIGNIN
				err = updateConsumer(tx, updataConsumerMap, consumerId)

				if err != nil {
					return false, "", err
				}

				_, err = insertConsumerStatusLog(tx, consumerId, employeeId, oldStatus, CONSUMER_STATUS_NO_SIGNIN)

				if err != nil {
					return false, "", err
				}
			}

			classDataMap, err := findWyClassById(classId)

			if err != nil {
				return false, "", err
			}

			note := fmt.Sprintf("邀约至%v%v中", classDataMap["start_time"], classDataMap["name"])

			_, err = insertConsumerContactsLog(tx, employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

			if err != nil {
				return false, "", err
			}

		} else {
			return false, "数据关联有问题，请联系IT部门,chilId:" + id, nil
		}
	}

	tx.Commit()

	return true, "", nil
}

func AddChildToClassQuick(classId, scheduleId, consumerId, employeeId string) (flag bool, msg string, err error) {

	childId, err := getChildByConsumerId(consumerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	return AddChildToClass(classId, scheduleId, fmt.Sprint(childId), employeeId)
}

func checkChildInSchedule(childId, scheduleId string) (bool, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select count(1) from schedule_detail_child where child_id=? and schedule_detail_id=? "

	rows, err := db.Query(sql, childId, scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	totalNum := 0

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if totalNum > 0 {
		return true, nil
	}

	return false, nil
}

func insertScheduleClass(tx *sql.Tx, childId, scheduleId, classId, employeeId string) error {

	sql := ""
	if classId != "" {
		sql += "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status,wyclass_id) values(?,?,?,?,?,?) "
	} else {
		sql += "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status) values(?,?,?,?,?) "
	}

	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	if classId != "" {
		_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employeeId, 1, classId)
	} else {
		_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employeeId, 1)
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func updateContractOfScheduleChild(tx *sql.Tx, contractId, childId, scheduleDetailId string) error {
	sql := "update schedule_detail_child set contract_id=? where child_id=? and schedule_detail_id=? "
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(contractId, childId, scheduleDetailId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func updateContractOfScheduleTmpChild(tx *sql.Tx, contractId, childId, scheduleTmpId string) error {
	sql := "update schedule_template_child set contract_id=? where child_id=? and schedule_template_id=?"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(contractId, childId, scheduleTmpId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func getFurtherScheduleIds(timeId, roomId, week, startTime, courseId string) ([]string, error) {

	sql := `
			select id from class_schedule_detail where time_id=? and room_id=? and week=? and start_time>=? and course_id=?`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, timeId, roomId, week, startTime, courseId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	ids := []string{}

	for rows.Next() {
		id := ""
		err := commonlib.PutRecord(rows, &id)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func CreateWeekSchedule(centerId, firstDayOfWeek, employeeId string) (flag bool, msg string, err error) {

	firstDay, err := time.ParseInLocation("20060102150405", firstDayOfWeek, time.Local)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	scheduleTempDatas, err := getScheduleTmpsByCenterId(centerId)

	if err != nil {
		return false, "", err
	}

	for _, scheduleTemp := range scheduleTempDatas {

		scheduleTempId := scheduleTemp["id"]
		week := scheduleTemp["week"]
		roomId := scheduleTemp["room_id"]
		timeId := scheduleTemp["time_id"]
		teacherId := scheduleTemp["teacher_id"]
		assistantId := scheduleTemp["assistant_id"]
		courseId := scheduleTemp["course_id"]

		weekNum, _ := strconv.Atoi(week)

		theDay := firstDay.Add(time.Duration((weekNum-1)*24) * time.Hour)
		date := theDay.Format("20060102")

		scheduleExistFlag, err := checkScheduleExist(centerId, roomId, timeId, date)

		if err != nil {
			return false, "", err
		}

		if !scheduleExistFlag { //相关课表还未生成
			scheduleId, err := insertSchedule(tx, teacherId, assistantId, courseId, centerId, timeId, roomId, date, week, "10")

			if err != nil {
				return false, "", err
			}

			childMapDatas, err := getChildAndContractByScheduleTempId(scheduleTempId)

			for _, childMapData := range childMapDatas {

				contractId := childMapData["contract_id"]

				if contractId == "0" || contractId == "" { //暂时没有合同信息的，直接进行排课操作
					err = insertScheduleChild(tx, fmt.Sprint(scheduleId), childMapData["child_id"], employeeId, "0")

					if err != nil {
						return false, "", err
					}
				}

				contractDataMap, err := getContractById(contractId)
				if err != nil {
					return false, "", err
				}

				validSignInNum, err := getVaildNumOfContract(contractId)
				if err != nil {
					return false, "", err
				}

				totalSignInNum, _ := strconv.Atoi(contractDataMap["left_lesson_num"])

				if totalSignInNum > validSignInNum {
					flag, err := checkContractValid(contractDataMap["expire_date"])
					if err != nil {
						return false, "", err
					}

					if flag {
						err = insertScheduleChild(tx, fmt.Sprint(scheduleId), childMapData["child_id"], employeeId, contractId)

						if err != nil {
							return false, "", err
						}
					}
				}

			}

		}
	}

	tx.Commit()

	return true, "", nil
}

func checkScheduleExist(centerId, roomId, timeId, dayDate string) (flag bool, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select id from class_schedule_detail where center_id=? and room_id=? and time_id=? and day_date=?"
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, centerId, roomId, timeId, dayDate)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	scheduleDetailId := 0

	if rows.Next() {
		err = commonlib.PutRecord(rows, &scheduleDetailId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if scheduleDetailId == 0 {
		return false, nil
	}

	return true, nil
}

func insertSchedule(tx *sql.Tx, teacherId, assistantId, courseId, centerId, timeId, roomId, dayDate, week, capacity string) (id int64, err error) {

	sql := "insert into class_schedule_detail(teacher_id,assistant_id,course_id,center_id,time_id,room_id,day_date,week,start_time,end_time,status,capacity) values(?,?,?,?,?,?,?,?,?,?,?,?)"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	timeSectionDataMap, err := getTimeSectionById(timeId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	startTime := dayDate + strings.Replace(timeSectionDataMap["start_time"], ":", "", -1) + "00"
	endTime := dayDate + strings.Replace(timeSectionDataMap["end_time"], ":", "", -1) + "00"

	res, err := stmt.Exec(teacherId, assistantId, courseId, centerId, timeId, roomId, dayDate, week, startTime, endTime, 1, capacity)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newScheduleDetailId, err := res.LastInsertId()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newScheduleDetailId, nil
}

func insertScheduleChild(tx *sql.Tx, scheduleId, childId, employeeId, contractId string) (err error) {
	sql := "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,contract_id) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employeeId, contractId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func ClassScheduleDetailLeave(childId, scheduleId, employeeId string) (flag bool, msg string, err error) {

	signInExistFlag, err := checkSignInExist(childId, scheduleId)

	if err != nil {
		return false, "", err
	}

	if signInExistFlag {
		return true, "", nil
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	contractId, err := getContractIdByChildIdAndScheduleId(childId, scheduleId)
	if err != nil {
		return false, "", err
	}

	//统一合同号没值的适合，为0
	if contractId == "" {
		contractId = "0"
	}

	err = insertSignIn(tx, scheduleId, childId, SIGN_IN_LEAVE, contractId, "", employeeId, IS_FREE_NO)

	if err != nil {
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}

func ClassScheduleDetailTruant(childId, scheduleId, employeeId string) (flag bool, msg string, err error) {

	signInExistFlag, err := checkSignInExist(childId, scheduleId)

	if err != nil {
		return false, "", err
	}

	if signInExistFlag {
		return true, "", nil
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	contractId, err := getContractIdByChildIdAndScheduleId(childId, scheduleId)
	if err != nil {
		return false, "", err
	}

	//统一合同号没值的适合，为0
	if contractId == "" {
		contractId = "0"
	}

	err = insertSignIn(tx, scheduleId, childId, SIGN_IN_TRUANT, contractId, "", employeeId, IS_FREE_NO)

	if err != nil {
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}
