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

	//todo 数据验证

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if id == "" {
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
			_, err = insertConsumerContactsLog(tx, createUser, remark, fmt.Sprint(consumerId))

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

			updateConsumerMap := make(map[string]interface {})
			updateConsumerMap["parent_id"] = parentId

			err = updateConsumer(tx,updateConsumerMap,fmt.Sprint(consumerId))
			if err != nil {
				return false, "", err
			}
		} else {
			updateConsumerMap := make(map[string]interface {})
			updateConsumerMap["parent_id"] = parentId

			err = updateConsumer(tx,updateConsumerMap,fmt.Sprint(consumerId))
			if err != nil {
				return false, "", err
			}
		}

	} else {
		consumerDataMap,err := findConsumerById(id)

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

		updateConsumerMap := make(map[string]interface {})
		updateConsumerMap["child"] = child
		updateConsumerMap["year"] = year
		updateConsumerMap["month"] = month
		updateConsumerMap["home_phone"] = homePhone
		updateConsumerMap["birthday"] = birthday
		updateConsumerMap["come_from_id"] = comeFromId
		updateConsumerMap["parttime_name"] = parttimeName
		updateConsumerMap["level"] = level

		err = updateConsumer(tx,updateConsumerMap,id)
		if err != nil {
			return false, "", err
		}

		childId,err := getChildByParentId(consumerDataMap["parentId"])
		if err != nil {
			return false, "", err
		}

		updateChildMap := make(map[string]interface {})
		updateChildMap["name"] = child

		err = updateChild(tx,updateChildMap,fmt.Sprint(childId))
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
func findConsumerById(id string) (map[string]string,error)  {

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

	dataMap,err := lessgo.GetDataMap(rows)

	if err != nil {
		return nil, err
	}

	return dataMap,nil
}

//更新consumer
func updateConsumer(tx *sql.Tx,consumerDataMap map[string]interface {},id string) error{
	sql := "update consumer_new set %v where id=?"
	params := []interface {}{}

	setSql := ""

	for key,value := range consumerDataMap{
		setSql += key+"=?,"
		params = append(params,value)
	}

	params = append(params,id)

	setSql = commonlib.Substr(setSql,0,len(setSql)-1)

	sql = fmt.Sprintf(sql,setSql)
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

func ConsumerPage(paramsMap map[string]string,dataType,employeeId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {/*
	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from consumer_new cons where 1=1 "
	countParams := []interface{}{scheduleId}
	lessgo.Log.Debug(countSql)

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `
				select sdc.child_id id,ch.name childName,p.telephone phone,si.type signType,si.sign_time signTime,cour.name courseName,contr.id as contractId,contr.contract_no contractNo,contr.apply_time applyTime,cons.id consumerId,cons.level level,d.remark
	 		    from (select * from schedule_detail_child where schedule_detail_id=? order by id desc limit ?,?) sdc
	 			left join child ch on ch.cid=sdc.child_id
	 			left join parent p on p.pid=ch.pid
	 			left join consumer_new cons on cons.parent_id=ch.pid
	            left join (select consumer_id,GROUP_CONCAT(concat(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i'),' ',note) ORDER BY id DESC SEPARATOR '<br/>') remark from consumer_contact_log group by consumer_id) d on d.consumer_id=cons.id
	 			left join sign_in si on si.child_id=sdc.child_id and sdc.schedule_detail_id=si.schedule_detail_id
	 			left join class_schedule_detail csd on csd.id=sdc.schedule_detail_id
	 			left join contract contr on contr.id=sdc.contract_id
	 			left join course cour on cour.cid=contr.course_id`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{scheduleId, (currPageNo-1)*pageSize, pageSize}

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}*/

	return nil, nil
}
