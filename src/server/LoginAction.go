// Title：登陆web方法
//
// Description:
//
// Author:black
//
// Createtime:2013-09-29 10:00
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-09-29 10:00 black 创建文档
package server

import (
	"net/http"
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
)


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

