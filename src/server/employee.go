// Title：登陆相关
//
// Description:
//
// Author:Samurai
//
// Createtime:2013-09-24 13:04
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
package server

import (
	"net/http"
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
)

func CheckPwd(username, password string) (bool, lessgo.Employee, string){

	var employee lessgo.Employee
	var dbPwd string

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select user_id,username,password,really_name,department_id from employee where username=?"

	rows, err := db.Query(sql, username)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, employee, " 数据库异常!"
	}

	if rows.Next() {
		err := rows.Scan(&employee.UserId, &employee.UserName,
						&dbPwd, &employee.ReallyName, &employee.DepartmentId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, employee, " 数据库异常!"
		}
	}

	if password != dbPwd {
		lessgo.Log.Warn("密码错误:", username, " : ", password)
		return false, employee, "密码错误"
	}

	return true, employee, ""
}

func LoginAction(w http.ResponseWriter,r *http.Request ) {

	data := make(map[string]interface {})

	username := r.FormValue("username")
	if username == "" {
		lessgo.Log.Warn("username is NULL!")
		return
	}

	password := r.FormValue("password")
	if password == "" {
		lessgo.Log.Warn("password is NULL!")
		return
	}

	ret, employee, msg := CheckPwd(username, password)

	if ret{
		//密码正确
		data["success"] = true
		lessgo.SetCurrentEmployee(employee,w,r)
	} else {
		data["success"] = false
		data["msg"] = msg
	}

	commonlib.OutputJson(w, data,"")

	return
}


