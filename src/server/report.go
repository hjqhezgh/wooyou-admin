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
)

type Model struct {
	key   string
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

	var modelLists []ModelList
	var report_name, report_remark, report_sql, employee_name, employee_mail string

	sql := `select r.name,r.remark,r.sql,e.really_name,e.mail from report_manage rm left join employee e
					on rm.employee_id=e.user_id left join report r on rm.report_id=r.id where rm.send_flag=1`

	lessgo.Log.Debug(sql)
	//var rcs []ReportColumn

	db := lessgo.GetMySQL()
	defer db.Close()
	rows, err := db.Query(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(&report_name, &report_remark, &report_sql, &employee_name, &employee_mail)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}
		//查询某张报表
		rows2, err := db.Query(report_sql)
		lessgo.Log.Debug(report_sql)
		for rows2.Next() {
			params := []interface{}{}
			modelList := new(ModelList)
			column, _ := rows2.Columns()

			for _, c := range column {
				m := new(Model)
				m.key = c
				params = append(params, &m.value)
				modelList.model = append(modelList.model, m)
			}
			err := rows2.Scan(params...)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return
			}
			//保存模型数据
			modelLists = append(modelLists, *modelList)
		}

		html := "<table table width='38%' border='1' cellpadding='2' cellspacing='0'><tr>"
		rows3, err := db.Query("select rc.sql_column,rc.column_name,rc.index from report_column rc where rc.report_id=? order by rc.`index`", 1)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}
		var ms []ReportColumn
		var t1, t2, t3 string
		for rows3.Next() {
			m := new(ReportColumn)
			err = rows3.Scan(&t1, &t2, &t3)
			html += "<th>"
			html += t2
			html += "</th>"
			ms = append(ms, *m)
		}
		lessgo.Log.Error(ms)
		html += "</tr>"

		str := ""
		td := ""
		for _, v := range modelLists {
			lessgo.Log.Debug(v)
			html += "<tr>"
			for _, vv := range v.model {
				if vv.value != nil {
					str = string(vv.value.([]byte))
					td = "<td>" + str + "</td>"
				} else {
					td = "<td></td>"
				}
				html += td
				lessgo.Log.Error(vv.value)

			}
			html += "</tr>"

		}
		html += "</table>"
		lessgo.Log.Error(html)

		sendMail(employee_mail, employee_name, report_name, html)
	}
	return
}

func sendMail(account, userName, subject, content string) (err error) {
	lessgo.Log.Debug("begin sendMail")
	user := "ibignose@126.com"
	password := "ibignose12345"
	host := "smtp.126.com:25"
	to := account
	body := "<strong>" + userName + "</strong>,  您好！: <br><br> "

	body += content
	lessgo.Log.Debug("sending email...")

	err = tool.SendMail(user, password, host, to, subject, body, "html")

	if err != nil {
		lessgo.Log.Error("send mail error!")
		lessgo.Log.Error(err.Error())
	} else {
		lessgo.Log.Debug("send mail success!")
	}
	return
}
