// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-12-02 10:19
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-12-02 10:19 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"strings"
	"time"
)

func WyClassSendSMSLoadAction(w http.ResponseWriter, r *http.Request) {

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

	id := r.FormValue("ids")
	classId := r.FormValue("classId")
	centerId := r.FormValue("centerId-eq")

	loadFormObjects := []lessgo.LoadFormObject{}

	findPhoneSql := "select cons.id,cont.phone,cons.child,ce.name,ce.intro from contacts cont left join consumer_new cons on cont.consumer_id=cons.id left join center ce on ce.cid=cons.center_id where cont.consumer_id in ("+id+")"
	lessgo.Log.Debug(findPhoneSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(findPhoneSql)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	phonesContent := ""
	phones := ""
	consumerIds := ""

	var centerName,centerAddress string

	for rows.Next() {
		var consumerId,phone,childName string
		err = commonlib.PutRecord(rows, &consumerId, &phone, &childName, &centerName, &centerAddress)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		if len(phone)!= 11{//不是手机号码
			continue
		}

		phonesContent += phone+"("+childName+"),"
		phones += phone+","
		consumerIds += consumerId+","
	}

	if phones == ""{
		m["success"] = false
		m["code"] = 100
		m["msg"] = "本次没有可以发送的号码，请重新选择"
		commonlib.OutputJson(w, m, " ")
		return
	}

	h1 := lessgo.LoadFormObject{"phonesContent", phonesContent}
	h2 := lessgo.LoadFormObject{"phones", phones}
	h4 := lessgo.LoadFormObject{"consumerIds", consumerIds}
	h5 := lessgo.LoadFormObject{"classId", classId}
	h6 := lessgo.LoadFormObject{"centerId", centerId}

	findTimeSql := "select name from wyclass where class_id=?"
	lessgo.Log.Debug(findTimeSql)

	rows, err = db.Query(findTimeSql, classId)

	if err != nil {
	lessgo.Log.Warn(err.Error())
	m["success"] = false
	m["code"] = 100
	m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
	commonlib.OutputJson(w, m, " ")
	return
	}

	var className string

	if rows.Next() {
	err = commonlib.PutRecord(rows, &className)

	if err != nil {
	lessgo.Log.Warn(err.Error())
	m["success"] = false
	m["code"] = 100
	m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
	commonlib.OutputJson(w, m, " ")
	return
	}
	}

	h3 := lessgo.LoadFormObject{"content", "欢迎您来试听"+centerName+className+"的课程，请提前10分钟到达"}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)

	m["success"] = true
	m["datas"] = loadFormObjects
	commonlib.OutputJson(w, m, " ")
	return
}

//短信发送
func WyClassSendSMSSaveAction(w http.ResponseWriter, r *http.Request) {
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

	classId := r.FormValue("classId")
	centerId := r.FormValue("centerId")
	consumerIds := r.FormValue("consumerIds")
	phones := r.FormValue("phones")
	content := r.FormValue("content")

	phoneList := strings.Split(phones,",")
	consumerIdList := strings.Split(consumerIds,",")

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	getCenterInfo := "select ce.name,ce.intro,class.name,class.start_time from wyclass class left join center ce on ce.cid=class.center_id where class.class_id=?"
	lessgo.Log.Debug(getCenterInfo)
	rows, err := db.Query(getCenterInfo, centerId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	var centerName ,centerIntro,className,classStartTime string

	if rows.Next() {
		err := rows.Scan(&centerName,&centerIntro,&className,&classStartTime)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	for index,phone := range phoneList {
		if phone!= ""{
			getConsumerInfo := "select child from consumer_new  where id=?"
			lessgo.Log.Debug(getConsumerInfo)
			rows, err = db.Query(getConsumerInfo, consumerIdList[index])

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			var childName string

			if rows.Next() {
				err := rows.Scan(&childName)

				if err != nil {
					lessgo.Log.Warn(err.Error())
					m["success"] = false
					m["code"] = 100
					m["msg"] = "系统发生错误，请联系IT部门"
					commonlib.OutputJson(w, m, " ")
					continue
				}
			}

			contentDetail := replaceSmsText(content,childName,centerIntro)
			smsResult ,err := SendMessage(phone, contentDetail)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			if smsResult.Msg == "" {//请求短信接口没有成功
				m["success"] = false
				m["code"] = 100
				m["msg"] = "短信发送错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			smsStatus := ""

			if smsResult.Result == 0 {//发送成功
				smsStatus  = "2"
			}else{
				smsStatus  = "3"
			}


			updateClassSmsStatus := "update wyclass_free_child set sms_status=? where wyclass_id=? and consumer_id=? "
			lessgo.Log.Debug(updateClassSmsStatus)

			stmt, err := tx.Prepare(updateClassSmsStatus)

			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(smsStatus,classId,consumerIdList[index])
			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}


			insertSmsLog := "insert into sms_send_log(phone,content,result_code,result_msg,send_time,send_status) values(?,?,?,?,?,?) "
			lessgo.Log.Debug(insertSmsLog)

			stmt, err = tx.Prepare(insertSmsLog)

			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门" + err.Error()
				commonlib.OutputJson(w, m, " ")
				return
			}

			_, err = stmt.Exec(phone,contentDetail,smsResult.Result,smsResult.Msg,time.Now().Format("20060102150405"),3)
			if err != nil {
				tx.Rollback()
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
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

func replaceSmsText(text,child,centerIntro string) string{
	return "欢迎您来试听厦门瑞景中心1208    19:00    英语的课程，请提前10分钟到达，您的注册验证码是13123123，请准时来临哟【吾幼儿童英语中心】"
}
