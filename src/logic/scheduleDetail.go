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
	"github.com/hjqhezgh/lessgo"
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

		consumerDataMap, err := findConsumerByChildId(id)

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

			_,err = insertConsumerContactsLog(tx, employeeId, note, consumerId, CONTACTS_LOG_TYPE_SYSTEM)

			if err != nil {
				return false, "", err
			}

		} else {
			return false, "数据关联有问题，请联系IT部门,chilId:"+id, nil
		}
	}

	tx.Commit()

	return true, "", nil
}

func AddChildToClassQuickAction(classId, scheduleId, consumerId, employeeId string) (flag bool, msg string, err error){

	childId,err := getChildByConsumerId(consumerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	return AddChildToClass(classId, scheduleId,fmt.Sprint(childId),employeeId)
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
