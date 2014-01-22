// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 16:23
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 16:23 black 创建文档
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

const (
	CONSUMER_STATUS_NO_CONTACT = "1" //未联系
	CONSUMER_STATUS_WAIT       = "2" //待确认
	CONSUMER_STATUS_ABANDON    = "3" //废弃
	CONSUMER_STATUS_NO_SIGNIN  = "4" //已邀约
	CONSUMER_STATUS_SIGNIN     = "5" //已签到
)

//客户保存
func SaveConsumer(paramsMap map[string]string) (flag bool, msg string, err error) {

	id := paramsMap["id"]
	phone := paramsMap["phone"]
	contactsName := paramsMap["contactsName"]
	homePhone := paramsMap["homePhone"]
	child := paramsMap["child"]
	year := paramsMap["year"]
	month := paramsMap["month"]
	birthday := paramsMap["birthday"]
	comeFromId := paramsMap["comeFromId"]
	centerId := paramsMap["centerId"]
	remark := paramsMap["remark"]
	parttimeName := paramsMap["parttimeName"]
	level := paramsMap["level"]
	createUser := paramsMap["createUser"]

	// todo 数据验证

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if id == "" {

		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if centerId == "1" {
			centerId = "7"
		}

		homePhoneFlag, err := checkConsumerPhoneExist(homePhone)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if homePhoneFlag {
			return false, "家庭电话已经在系统中存在，无需重复录入", nil
		}

		phoneFlag, err := checkConsumerPhoneExist(phone)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if phoneFlag {
			return false, "联系人电话已经存在，无需重复录入", nil
		}

		contactStatus := CONSUMER_STATUS_NO_CONTACT
		lastContactTime := ""

		if remark != "" {
			contactStatus = CONSUMER_STATUS_WAIT
			lastContactTime = time.Now().Format("20060102150405")
		}

		consumerId, err := insertConsumer(tx, centerId, contactStatus, homePhone, child, year, month, birthday, lastContactTime, comeFromId, createUser, parttimeName)

		if err != nil {
			return false, "", err
		}

		_, err = insertContacts(tx, contactsName, phone, fmt.Sprint(consumerId))

		if err != nil {
			return false, "", err
		}

		if remark != "" {
			_, err = insertConsumerContactsLog(tx, createUser, remark, fmt.Sprint(consumerId), CONTACTS_LOG_TYPE_PHONE)

			if err != nil {
				return false, "", err
			}
		}

		parentId, err := getParentIdByPhone(phone)

		if err != nil {
			return false, "", err
		}

		if parentId == 0 { //parent不存在
			newParentName := child + "家长"
			if contactsName != "" {
				newParentName = contactsName
			}

			parentId, err = insertParent(tx, newParentName, "123456", phone, "2")

			if err != nil {
				return false, "", err
			}

			_, err = insertChild(tx, child, fmt.Sprint(parentId), "1", birthday, year, month, centerId)

			if err != nil {
				return false, "", err
			}

			updateConsumerMap := make(map[string]interface{})
			updateConsumerMap["parent_id"] = parentId

			err = updateConsumer(tx, updateConsumerMap, fmt.Sprint(consumerId))
			if err != nil {
				return false, "", err
			}
		} else {
			updateConsumerMap := make(map[string]interface{})
			updateConsumerMap["parent_id"] = parentId

			err = updateConsumer(tx, updateConsumerMap, fmt.Sprint(consumerId))
			if err != nil {
				return false, "", err
			}
		}

	} else {
		consumerDataMap, err := getConsumerById(id)

		if err != nil {
			return false, "", err
		}

		oldHomePhone := consumerDataMap["homePhone"]

		if oldHomePhone != homePhone {
			homePhoneFlag, err := checkConsumerPhoneExist(homePhone)
			if err != nil {
				return false, "", err
			}

			if homePhoneFlag {
				return false, "家庭电话已经在系统中存在，无需重复录入", nil
			}
		}

		updateConsumerMap := make(map[string]interface{})
		updateConsumerMap["child"] = child
		updateConsumerMap["year"] = year
		updateConsumerMap["month"] = month
		updateConsumerMap["home_phone"] = homePhone
		updateConsumerMap["birthday"] = birthday
		updateConsumerMap["come_from_id"] = comeFromId
		updateConsumerMap["parttime_name"] = parttimeName
		updateConsumerMap["level"] = level

		err = updateConsumer(tx, updateConsumerMap, id)
		if err != nil {
			return false, "", err
		}

		childId, err := getChildByParentId(consumerDataMap["parentId"])
		if err != nil {
			return false, "", err
		}

		updateChildMap := make(map[string]interface{})
		updateChildMap["name"] = child

		err = updateChild(tx, updateChildMap, fmt.Sprint(childId))
		if err != nil {
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

//判断客户电话是否已经存在
func checkConsumerPhoneExist(phone string) (bool, error) {

	if phone == "" {
		return false, nil
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select count(1) from consumer_new where home_phone=? "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	sql = "select count(1) from contacts where phone=? "

	lessgo.Log.Debug(sql)

	rows, err = db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num = 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	return false, nil
}

func insertConsumer(tx *sql.Tx, centerId, contactStatus, homePhone, child, year, month, birthday, lastContactTime, comeFromId, createUser, parttimeName string) (id int64, err error) {

	sql := "insert into consumer_new(center_id,contact_status,home_phone,create_time,child,year,month,birthday,last_contact_time,come_from_id,create_user_id,parttime_name) values(?,?,?,?,?,?,?,?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(centerId, contactStatus, homePhone, time.Now().Format("20060102150405"), child, year, month, birthday, lastContactTime, comeFromId, createUser, parttimeName)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	consumerId, err := res.LastInsertId()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return consumerId, err
}

//根据id获取consumer数据Map
/*
select id,center_id centerId,contact_status contactStatus,home_phone homePhone,parent_id parentId,child,year,month,birthday,
	    	       last_tmk_id lastTMKId,is_own_by_tmk isOwnByTmk,come_from_id comeFromId,current_tmk_id,currentTMKId,sign_in_time signInTime,
	    	       pay_time,payTime,pay_status payStatus,parttime_name parttimeName,level
	        from consumer_new where id=?
*/
func getConsumerById(id string) (map[string]string, error) {

	sql := `select id,center_id centerId,contact_status contactStatus,home_phone homePhone,parent_id parentId,child,year,month,birthday,
	    	       last_tmk_id lastTMKId,is_own_by_tmk isOwnByTmk,come_from_id comeFromId,current_tmk_id currentTMKId,sign_in_time signInTime,
	    	       pay_time payTime,pay_status payStatus,parttime_name parttimeName,level
	        from consumer_new where id=?`

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

//更新consumer
func updateConsumer(tx *sql.Tx, consumerDataMap map[string]interface{}, id string) error {
	sql := "update consumer_new set %v where id=?"
	params := []interface{}{}

	setSql := ""

	for key, value := range consumerDataMap {
		setSql += key + "=?,"
		params = append(params, value)
	}

	params = append(params, id)

	setSql = commonlib.Substr(setSql, 0, len(setSql)-1)

	sql = fmt.Sprintf(sql, setSql)
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(params...)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

func ConsumerPage(paramsMap map[string]string, dataType, employeeId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	whereSql := ""

	countSql := `
		select count(1) from (
			select c.consumer_id from contacts c
			left join consumer_new b on b.id=c.consumer_id
			where (c.phone like ? or c.name like ? or b.child like ?) %v
			group by c.consumer_id
		) num
	`
	countParams := []interface{}{"%" + paramsMap["kw"] + "%", "%" + paramsMap["kw"] + "%", "%" + paramsMap["kw"] + "%"}

	dataSql := `
		select cons.id id,ce.name as centerName,e.really_name tmkName,cont.name contactName,cont.phone phone ,cons.home_phone homePhone,cons.child childName,cons.contact_status contactStatus,cons.parent_id parentId,b.remark remark,cons.pay_status payStatus,cons.pay_time payTime,cf.name comeFromName, cons.parttime_name partTimeName
		from (
			select c.consumer_id,min(c.id) contacts_id from contacts c
			left join consumer_new b on b.id=c.consumer_id
			where (c.phone like ? or c.name like ? or b.child like ?) %v
			group by c.consumer_id %v limit ?,?)a
		left join consumer_new cons on cons.id=a.consumer_id
		left join contacts cont on cont.id=a.contacts_id
		left join center ce on ce.cid=cons.center_id
		left join employee e on e.user_id=cons.current_tmk_id
		left join come_from cf on cf.id=cons.come_from_id
	`
	dataParams := []interface{}{"%" + paramsMap["kw"] + "%", "%" + paramsMap["kw"] + "%", "%" + paramsMap["kw"] + "%"}

	if paramsMap["status"] != "" {
		countParams = append(countParams, paramsMap["status"])
		dataParams = append(dataParams, paramsMap["status"])
		whereSql += " and b.contact_status=? "
	}

	if paramsMap["lastContractStartTime"] != "" && paramsMap["timeType"] == "1" {
		countParams = append(countParams, paramsMap["lastContractStartTime"])
		dataParams = append(dataParams, paramsMap["lastContractStartTime"])
		whereSql += " and b.sign_in_time>=? "
	}

	if paramsMap["lastContractStartTime"] != "" && paramsMap["timeType"] == "2" {
		countParams = append(countParams, paramsMap["lastContractStartTime"])
		dataParams = append(dataParams, paramsMap["lastContractStartTime"])
		whereSql += " and b.pay_time>=? "
	}

	if paramsMap["lastContractEndTime"] != "" && paramsMap["timeType"] == "1" {
		countParams = append(countParams, paramsMap["lastContractEndTime"])
		dataParams = append(dataParams, paramsMap["lastContractEndTime"])
		whereSql += " and b.sign_in_time<=? "
	}

	if paramsMap["lastContractEndTime"] != "" && paramsMap["timeType"] == "2" {
		countParams = append(countParams, paramsMap["lastContractEndTime"])
		dataParams = append(dataParams, paramsMap["lastContractEndTime"])
		whereSql += " and b.pay_time<=? "
	}

	if paramsMap["payStatus"] == "1" {
		whereSql += " and b.pay_time is not null and b.pay_time != '' and b.pay_status=1 "
	} else if paramsMap["payStatus"] == "2" {
		whereSql += " and b.pay_time is not null and b.pay_time != '' and b.pay_status=2 "
	} else if paramsMap["payStatus"] == "3" {
		whereSql += " and (b.pay_time is null or b.pay_time ='')"
	}

	if paramsMap["tmkId1"] != "" {
		countParams = append(countParams, paramsMap["tmkId1"])
		dataParams = append(dataParams, paramsMap["tmkId1"])
		whereSql += " and b.current_tmk_id=? "
	}

	if paramsMap["tmkId2"] != "" {
		countParams = append(countParams, paramsMap["tmkId2"])
		dataParams = append(dataParams, paramsMap["tmkId2"])
		whereSql += " and b.current_tmk_id=? "
	}

	if dataType == "center" {
		userId, _ := strconv.Atoi(employeeId)
		_employee, err := FindEmployeeById(userId)
		if err != nil {
			return nil, err
		}

		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if _employee.CenterId == "1" {
			_employee.CenterId = "7"
		}

		countParams = append(countParams, _employee.CenterId)
		dataParams = append(dataParams, _employee.CenterId)

		whereSql += " and b.center_id=? "
	}

	if paramsMap["centerId1"] != "" && dataType == "all" {
		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if paramsMap["centerId1"] == "1" {
			countParams = append(countParams, "7")
			dataParams = append(dataParams, "7")
		} else {
			countParams = append(countParams, paramsMap["centerId1"])
			dataParams = append(dataParams, paramsMap["centerId1"])
		}

		whereSql += " and b.center_id=? "
	}

	if paramsMap["centerId2"] != "" && dataType == "all" {
		//番茄田逻辑补丁，番茄田添加的用户都属于福州台江中心
		if paramsMap["centerId2"] == "1" {
			countParams = append(countParams, "7")
			dataParams = append(dataParams, "7")
		} else {
			countParams = append(countParams, paramsMap["centerId2"])
			dataParams = append(dataParams, paramsMap["centerId2"])
		}
		whereSql += " and b.center_id=? "
	}

	if paramsMap["parttimeName"] != "" {
		countParams = append(countParams, paramsMap["parttimeName"])
		dataParams = append(dataParams, paramsMap["parttimeName"])
		whereSql += " and b.parttime_name=? "
	}

	if paramsMap["comeFromId"] != "" {
		countParams = append(countParams, paramsMap["comeFromId"])
		dataParams = append(dataParams, paramsMap["comeFromId"])
		whereSql += " and b.come_from_id=? "
	}

	orderSql := ""

	if paramsMap["sort"] == "" || paramsMap["sort"] == "create_time" {
		orderSql += " order by b.id desc "
	} else if paramsMap["sort"] == "last_time" {
		orderSql += " order by b.last_contact_time desc "
	}

	countSql = fmt.Sprintf(countSql, whereSql)
	dataSql = fmt.Sprintf(dataSql, whereSql, orderSql)

	dataSql += ` left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note)  ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) b on b.consumer_id=cons.id `

	lessgo.Log.Debug(countSql)
	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataParams = append(dataParams, (currPageNo-1)*pageSize, pageSize)

	lessgo.Log.Debug(dataSql)
	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

//根据ChildId获取consumer数据Map
/*
select cons.id,cons.center_id centerId,cons.contact_status contactStatus,cons.home_phone homePhone,cons.parent_id parentId,cons.child,cons.year,cons.month,cons.birthday,cons.last_tmk_id lastTMKId,cons.is_own_by_tmk isOwnByTmk,cons.come_from_id comeFromId,cons.current_tmk_id currentTMKId,cons.sign_in_time signInTime,cons.pay_time payTime,cons.pay_status payStatus,cons.parttime_name parttimeName,cons.level
	    	from consumer_new cons left join child ch on ch.pid=cons.parent_id where ch.cid=?
*/
func getConsumerByChildId(id string) (map[string]string, error) {

	sql := `
			select cons.id,cons.center_id centerId,cons.contact_status contactStatus,cons.home_phone homePhone,cons.parent_id parentId,cons.child,cons.year,cons.month,cons.birthday,cons.last_tmk_id lastTMKId,cons.is_own_by_tmk isOwnByTmk,cons.come_from_id comeFromId,cons.current_tmk_id currentTMKId,cons.sign_in_time signInTime,cons.pay_time payTime,cons.pay_status payStatus,cons.parttime_name parttimeName,cons.level
	    	from consumer_new cons left join child ch on ch.pid=cons.parent_id where ch.cid=?
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

func ConsumerPay(consumerIds, payType, employeeId string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	idList := strings.Split(consumerIds, ",")

	for _, consumerId := range idList {
		flag, msg, err = childPay(tx, "", consumerId, "", "", payType, employeeId)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		if !flag {
			return false, msg, nil
		}
	}

	tx.Commit()

	return true, "", nil
}
