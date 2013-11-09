// Title：
//
// Description:
//
// Author:samurai
//
// Createtime:2013-10-14 14:14
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-14 14:14 user04 创建文档package test

package server

import (
	//"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	//"fmt"
	"tool"
	"reflect"
)

type Model struct {
	value interface{}
}

type ModelList struct {
	model []*Model
}

type ReportColumn struct {
	sqlName    string
	columnName string
	index      int
}

func ReportSend() {

	var report_name, report_remark, report_sql, employee_name, employee_mail string

	sql := `select r.name,r.remark,r.sql,e.really_name,e.mail from report_manage rm left join employee e
					on rm.employee_id=e.user_id left join report r on rm.report_id=r.id where rm.send_flag=1`

	html := `<style>body {color:#4f6b72;font-size:14px;}th{background:#CAE8EA;height:45px;padding:0 12px;line-height:45px;font-size:16px;text-align:center;}tr{background:#fff;}td {padding: 10px 6px 10px 12px;}tr.alt{background:#f5fafa;}td.spec{font-weight:700;color:#313131;}</style>`
	var table string
	str := ""
	td := ""

	db := lessgo.GetMySQL()
	defer db.Close()
	rows, err := db.Query(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}
	for rows.Next() {
		table = ""
		err := rows.Scan(&report_name, &report_remark, &report_sql, &employee_name, &employee_mail)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}
		//查询某张报表
		rows2, err := db.Query(report_sql)
		if err != nil {
			lessgo.Log.Error(err.Error())
			continue
		}
		table += "<table table=‘table’ border=‘0’ cellpadding=‘0’ cellspacing=‘1’ bgcolor=‘#C1DAD7’>"
		column, _ := rows2.Columns()
		table += "<tr>"
		params := []interface{}{}
		for _, c := range column {
			table += "<th>"
			table += c
			table += "</th>"
			m := new(Model)
			params = append(params, &m.value)
		}
		table += "</tr>"
		for rows2.Next() {

			err := rows2.Scan(params...)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return
			}

			table += "<tr>"
			//拼接table数据
			//lessgo.Log.Debug(params...)
			for _, v := range params {
				tmp := reflect.Indirect(reflect.ValueOf(v)).Interface()
				if tmp != nil {
					str = string(tmp.([]byte))
					//lessgo.Log.Debug(string(tmp.([]byte)))
					td = "<td>" + str + "</td>"
				} else {
					td = "<td></td>"
				}
				table += td
			}
			table += "</tr>"
		}

		html += "</tr>"
		table += "</table>"
		html += table

	}

	sendMail(employee_mail, employee_name, report_name, html)
	return
}

func sendMail(account, userName, subject, content string) (err error) {
	user := "ibignose@126.com"
	password := "ibignose12345"
	host := "smtp.126.com:25"
	to := account
	body := "<strong>" + userName + "</strong>,  您好！: <br><br> "

	body += content

	err = tool.SendMail(user, password, host, to, subject, body, "html")

	if err != nil {
		lessgo.Log.Error(err.Error())
	} else {
		lessgo.Log.Debug("send mail success!")
	}
	return
}
