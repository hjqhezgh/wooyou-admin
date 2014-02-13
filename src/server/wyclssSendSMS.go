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
	"fmt"
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
	scheduleId := r.FormValue("scheduleId")
	centerId := r.FormValue("centerId-eq")

	loadFormObjects := []lessgo.LoadFormObject{}

	findPhoneSql := "select ch.cid,p.telephone,ch.name from child ch left join parent p on p.pid=ch.pid where ch.cid in (" + id + ")"
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

	for rows.Next() {
		var consumerId, phone, childName string
		err = commonlib.PutRecord(rows, &consumerId, &phone, &childName)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		if len(phone) != 11 { //不是手机号码
			continue
		}

		phonesContent += phone + "(" + childName + "),"
		phones += phone + ","
		consumerIds += consumerId + ","
	}

	if phones == "" {
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
	h7 := lessgo.LoadFormObject{"scheduleId", scheduleId}

	getCenterInfo := "select ce.name,ce.intro,class.name,class.start_time from wyclass class left join center ce on ce.cid=class.center_id where class.class_id=?"
	lessgo.Log.Debug(getCenterInfo)
	rows, err = db.Query(getCenterInfo, classId)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	var centerName, centerIntro, className, classStartTime string

	if rows.Next() {
		err := commonlib.PutRecord(rows, &centerName, &centerIntro, &className, &classStartTime)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	h3 := lessgo.LoadFormObject{"content", getSmsTmpText(employee.ReallyName, "$child", centerIntro, classStartTime)}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)
	loadFormObjects = append(loadFormObjects, h4)
	loadFormObjects = append(loadFormObjects, h5)
	loadFormObjects = append(loadFormObjects, h6)
	loadFormObjects = append(loadFormObjects, h7)

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
	scheduleId := r.FormValue("scheduleId")
	consumerIds := r.FormValue("consumerIds")
	phones := r.FormValue("phones")
	content := r.FormValue("content")

	phoneList := strings.Split(phones, ",")
	consumerIdList := strings.Split(consumerIds, ",")

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

	for index, phone := range phoneList {
		if phone != "" {
			getConsumerInfo := "select name from child  where cid=?"
			lessgo.Log.Debug(getConsumerInfo)
			rows, err := db.Query(getConsumerInfo, consumerIdList[index])

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

			contentDetail := getSmsContent(content, childName)
//			smsResult := SmsResult{Msg: "asdas", Result: 0}
			smsResult, err := SendMessage(phone, contentDetail)
			fmt.Println(contentDetail)
			if err != nil {
				lessgo.Log.Warn(err.Error())
				m["success"] = false
				m["code"] = 100
				m["msg"] = "系统发生错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			if smsResult.Result != 0 { //请求短信接口没有成功
				m["success"] = false
				m["code"] = 100
				m["msg"] = "短信发送错误，请联系IT部门"
				commonlib.OutputJson(w, m, " ")
				return
			}

			smsStatus := ""

			if smsResult.Result == 0 { //发送成功
				smsStatus = "2"
			} else {
				smsStatus = "3"
			}

			updateClassSmsStatus := "update schedule_detail_child set sms_status=? where wyclass_id=? and child_id=? and schedule_detail_id=? "
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

			_, err = stmt.Exec(smsStatus, classId, consumerIdList[index], scheduleId)
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

			_, err = stmt.Exec(phone, contentDetail, smsResult.Result, smsResult.Msg, time.Now().Format("20060102150405"), 3)
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

//获取短信模板文本
func getSmsTmpText(employeeName, child, centerIntro, startTime string) string {
	st, _ := time.ParseInLocation("20060102150405", startTime, time.Local)
	content := child + "家长，您好，我是您吾幼儿童社区的老师" + employeeName + "。咱们约定的时间：" + st.Format("2006-01-02 15:04") + "。" + centerIntro + "到时找不到地址直接打中心电话确认。咱们官网：www.wooyou.com.cn您可以提前上网了解。祝生活愉快！"
	content += "【吾幼英语美术社区】"
	return content
}

//替换为具体的短信模板
func getSmsContent(content, childName string) string {
	return strings.Replace(content, "$child", childName, -1)
}
