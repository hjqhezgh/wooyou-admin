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

		err = insertScheduleChild(tx, id, scheduleId, classId,"0", employeeId,IS_FREE_YES)

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

			note := ""

			if classId == ""{//常规课程
				scheduleDataMap,err := getScheduleDetailId(scheduleId)

				if err != nil {
					return false, "", err
				}

				note = fmt.Sprintf("邀约至[时间]%v[教室]%v[课程]%v中", scheduleDataMap["startTime"], scheduleDataMap["roomName"], scheduleDataMap["courseName"])
			}else{
				classDataMap, err := getWyClassById(classId)

				if err != nil {
					return false, "", err
				}

				note = fmt.Sprintf("邀约至%v%v中", classDataMap["start_time"], classDataMap["name"])
			}

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

func insertScheduleChild(tx *sql.Tx, childId, scheduleId, classId, employeeId,contractId,isFree string) error {

	sql := ""
	if classId != "" {
		sql += "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status,wyclass_id,is_free) values(?,?,?,?,?,?,?) "
	} else {
		sql += "insert into schedule_detail_child(schedule_detail_id,child_id,create_time,create_user,sms_status,is_free) values(?,?,?,?,?,?) "
	}

	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	if classId != "" {
		_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employeeId, 1, classId,isFree)
	} else {
		_, err = stmt.Exec(scheduleId, childId, time.Now().Format("20060102150405"), employeeId, 1,isFree)
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
					err = insertScheduleChild(tx, fmt.Sprint(scheduleId),"", childMapData["child_id"], employeeId, "0",IS_FREE_NO)

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
						err = insertScheduleChild(tx, fmt.Sprint(scheduleId),"", childMapData["child_id"], employeeId, contractId,IS_FREE_NO)

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

	contractId,_, err := getContractIdByChildIdAndScheduleId(childId, scheduleId)
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

	contractId,_, err := getContractIdByChildIdAndScheduleId(childId, scheduleId)
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

func ClassScheduleDetailSignIn(childIds, scheduleId,classId, employeeId string) (flag bool, msg string, err error) {

	idList := strings.Split(childIds, ",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	for _, childId := range idList {
		signInExistFlag, err := checkSignInExist(childId, scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if signInExistFlag {
			continue
		}

		contractId,isFree, err := getContractIdByChildIdAndScheduleId(childId, scheduleId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		//统一合同号没值的适合，为0
		if contractId == "" {
			contractId = "0"
		}

		childDataMap,err := getChildById(childId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		err = insertSignIn(tx, scheduleId, childId, SIGN_IN_SUCCESS, contractId, childDataMap["cardId"], employeeId, isFree)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if isFree == IS_FREE_YES {//免费课有额外逻辑
			consumerDataMap,err := getConsumerByChildId(childId)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			consumerId := consumerDataMap["id"]

			if consumerId == ""{
				return false, "数据关联出错，请联系it部门，childId：" + childId, nil
			}

			consumerUpdateMap := make(map[string]interface{})
			consumerUpdateMap["contact_status"] = CONSUMER_STATUS_SIGNIN
			consumerUpdateMap["sign_in_time"] = time.Now().Format("20060102150405")

			err = updateConsumer(tx,consumerUpdateMap,consumerId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			_,err = insertConsumerStatusLog(tx,consumerId, employeeId, consumerDataMap["contactStatus"], CONSUMER_STATUS_SIGNIN)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			note := ""

			if classId == ""{//常规课
				scheduleDataMap,err := getScheduleDetailId(scheduleId)

				if err != nil {
					return false, "", err
				}
				note = fmt.Sprintf("签到常规课[时间]%v[教室]%v[课程]%v中", scheduleDataMap["startTime"], scheduleDataMap["roomName"], scheduleDataMap["courseName"])
			}else{
				classDataMap,err := getWyClassById(classId)

				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}

				note = "签到"+classDataMap["start_time"]+classDataMap["name"]
			}

			insertConsumerContactsLog(tx,employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}
		}
	}

	tx.Commit()

	return true, "", nil
}

func ChildSignInWithoutClass(consumerIds,employeeId string) (flag bool, msg string, err error){

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	idList := strings.Split(consumerIds, ",")

	for _, consumerId := range idList {

		consumerDataMap,err := getConsumerById(consumerId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if consumerDataMap["contactStatus"] == CONSUMER_STATUS_SIGNIN {
			continue
		}

		childId,err := getChildByConsumerId(consumerId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		signInFlag,err := checkSignInExist(fmt.Sprint(childId),"")

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if signInFlag{
			continue
		}

		err = insertSignIn(tx,"",fmt.Sprint(childId),SIGN_IN_SUCCESS,"0","",employeeId,IS_FREE_YES)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		consumerUpdateMap := make(map[string]interface{})
		consumerUpdateMap["contact_status"] = CONSUMER_STATUS_SIGNIN
		consumerUpdateMap["sign_in_time"] = time.Now().Format("20060102150405")

		err = updateConsumer(tx,consumerUpdateMap,consumerId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		_,err = insertConsumerStatusLog(tx,consumerId, employeeId, consumerDataMap["contactStatus"], CONSUMER_STATUS_SIGNIN)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		insertConsumerContactsLog(tx,employeeId, "无班签到", consumerId, CONTACTS_LOG_TYPE_SYSTEM)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

/*
select csd.id,class_id classId,teacher_id teacherId,assistant_id assistantId,course_id courseId,center_id centerId,time_id timeId,room_id roomId,day_date dayDate,week,capacity,start_time startTime,end_time endTime,status,cour.name courseName,r.name roomName
			from class_schedule_detail csd
			left join course cour on cour.cid=csd.course_id
			left join room r on r.rid = csd.room_id
			where csd.id=?
*/
func getScheduleDetailId(id string) (map[string]string, error) {

	sql := `
			select csd.id,csd.class_id classId,csd.teacher_id teacherId,csd.assistant_id assistantId,csd.course_id courseId,csd.center_id centerId,csd.time_id timeId,csd.room_id roomId,csd.day_date dayDate,csd.week,csd.capacity,csd.start_time startTime,csd.end_time endTime,csd.status,cour.name courseName,r.name roomName
			from class_schedule_detail csd
			left join course cour on cour.cid=csd.course_id
			left join room r on r.rid = csd.room_id
			where csd.id=?
	    	`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var dataMap map[string]string

	if rows.Next() {
		dataMap, err = lessgo.GetDataMap(rows)
	}

	if err != nil {
		return nil, err
	}

	return dataMap, nil
}

func AddChildForNormalTempelate(childIds, scheduleId, employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	idList := strings.Split(childIds, ",")

	for _, childId := range idList {
		scheduleTempDataMap,err := getScheduleTmpByScheduelDetailId(scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if scheduleTempDataMap["id"] == "" {
			return false, "该课表不是模板课表，无法跟班", nil
		}

		timeId := scheduleTempDataMap["time_id"]
		roomId := scheduleTempDataMap["room_id"]
		week := scheduleTempDataMap["week"]
		startTime := scheduleTempDataMap["start_time"]
		courseId := scheduleTempDataMap["course_id"]

		furtherScheduleIds ,err := getFurtherScheduleIds(timeId,roomId,week,startTime,courseId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		for _,furtherScheduleId := range furtherScheduleIds{
			childScheduleFlag,err := checkChildInSchedule(childId,furtherScheduleId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			contractId,_, err := getContractIdByChildIdAndScheduleId(childId, furtherScheduleId)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			//统一合同号没值的适合，为0
			if contractId == "" {
				contractId = "0"
			}

			if !childScheduleFlag{
				err = insertScheduleChild(tx,childId,furtherScheduleId,"",employeeId,contractId,IS_FREE_NO)
				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}
			}
		}

		scheduleTempChildExistFlag,err := checkScheduleTempChildExist(childId,scheduleTempDataMap["id"])

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if !scheduleTempChildExistFlag{
			err = insertScheduleTempChild(tx,childId,scheduleTempDataMap["id"])
			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}
		}
	}

	tx.Commit()

	return true, "", nil
}

func AddChildForNormalOnce(childIds, scheduleId, employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	idList := strings.Split(childIds, ",")

	for _, childId := range idList {
		childScheduleFlag,err := checkChildInSchedule(childId,scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if !childScheduleFlag{
			err = insertScheduleChild(tx,childId,scheduleId,"",employeeId,"0",IS_FREE_NO)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}
		}
	}

	tx.Commit()

	return true, "", nil
}

func RemoveChildFromSchedule(childId,scheduleId,classId,dataType,employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	scheduleChildDataMap,err := getScheduleChildByChildIdAndScheduleId(childId,scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	signInExistFlag,err := checkSignInExist(childId,scheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if signInExistFlag {
		return false,"已签到/请假/旷课的无法剔除",nil
	}

	if scheduleChildDataMap["is_free"] == IS_FREE_YES {
		if dataType != "all" {//主管级别的是可以直接从班级中剔除任何人的，其他人只能剔除自己邀约的
			if employeeId != scheduleChildDataMap["create_user"] {
				return false,"您不能剔除其他人邀请的学员",nil
			}
		}

		err = deleteScheduleChild(tx,childId,scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		consumerDataMap,err := getConsumerByChildId(childId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		consumerId := consumerDataMap["id"]
		if consumerId == "" {
			return false,"数据关联出错，childId:"+childId,nil
		}

		if consumerDataMap["contactStatus"] != CONSUMER_STATUS_SIGNIN {
			consumerUpdateMap := make(map[string]interface{})
			consumerUpdateMap["contact_status"] = CONSUMER_STATUS_WAIT

			err = updateConsumer(tx,consumerUpdateMap,consumerId)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			_,err = insertConsumerStatusLog(tx,consumerId, employeeId, CONSUMER_STATUS_NO_SIGNIN, CONSUMER_STATUS_WAIT)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			insertConsumerContactsLog(tx,employeeId, "无班签到", consumerId, CONTACTS_LOG_TYPE_SYSTEM)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}

			note := ""

			if classId == ""{//常规课
				scheduleDataMap,err := getScheduleDetailId(scheduleId)

				if err != nil {
					return false, "", err
				}
				note = fmt.Sprintf("从常规课[时间]%v[教室]%v[课程]%v中剔除", scheduleDataMap["startTime"], scheduleDataMap["roomName"], scheduleDataMap["courseName"])
			}else{
				classDataMap,err := getWyClassById(classId)

				if err != nil {
					lessgo.Log.Error(err.Error())
					return false, "", err
				}

				note = "从"+classDataMap["start_time"]+classDataMap["name"]+"中剔除"
			}

			insertConsumerContactsLog(tx,employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

			if err != nil {
				lessgo.Log.Error(err.Error())
				return false, "", err
			}
		}

	}else{
		err = deleteScheduleChild(tx,childId,scheduleId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

func ChangeClassSchedule(childId, newScheduleId, oldScheduleId, employeeId string) (flag bool, msg string, err error){

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	signInExistFlag,err := checkSignInExist(childId,oldScheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if signInExistFlag{
		return false,"已签到/请假/旷课无法调班",nil
	}

	childInScheduleFlag,err := checkChildInSchedule(childId,newScheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if childInScheduleFlag{
		return false,"学生已在新班级中，无需重复排课",nil
	}

	scheduleChildDataMap,err := getScheduleChildByChildIdAndScheduleId(childId,oldScheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	err = deleteScheduleChild(tx,childId,oldScheduleId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	err = insertScheduleChild(tx, childId, newScheduleId, scheduleChildDataMap["wyclass_id"], scheduleChildDataMap["create_user"],scheduleChildDataMap["contract_id"],scheduleChildDataMap["is_free"])

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	newScheduleDataMpa,err := getScheduleDetailId(newScheduleId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}


	oldScheduleDataMpa,err := getScheduleDetailId(oldScheduleId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	note := "从"

	if oldScheduleDataMpa["classId"] == ""{//常规课
		note += fmt.Sprintf("常规课[时间]%v[教室]%v[课程]%v中", oldScheduleDataMpa["startTime"], oldScheduleDataMpa["roomName"], oldScheduleDataMpa["courseName"])
	}else{
		classDataMap,err := getWyClassById(oldScheduleDataMpa["classId"])

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		note += classDataMap["start_time"]+classDataMap["name"]
	}

	note += "调至"

	if newScheduleDataMpa["classId"] == ""{//常规课
		note += fmt.Sprintf("常规课[时间]%v[教室]%v[课程]%v中", newScheduleDataMpa["startTime"], newScheduleDataMpa["roomName"], newScheduleDataMpa["courseName"])
	}else{
		classDataMap,err := getWyClassById(newScheduleDataMpa["classId"])

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		note += classDataMap["start_time"]+classDataMap["name"]
	}

	consumerDataMap,err := getConsumerByChildId(childId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	_,err = insertConsumerContactsLog(tx,employeeId,note, consumerDataMap["id"], CONTACTS_LOG_TYPE_SYSTEM)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}
