// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-09 11:14
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-09 11:14 black 创建文档
package server

import (
	"net/http"
	"github.com/hjqhezgh/commonlib"
)

func HeHe(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	state := r.FormValue("state")
	suffix := r.FormValue("suffix")
	fileName := r.FormValue("fileName")
	srcFileName := r.FormValue("srcFileName")
	fileSize := r.FormValue("fileSize")

	if state == "SUCCESS" {
		m["success"] = true
		m["suffix"] = suffix
		m["fileName"] = fileName
		m["srcFileName"] = srcFileName
		m["fileSize"] = fileSize
		commonlib.OutputJson(w, m, " ")
		return
	}else{
		m["success"] = false
		m["msg"] = "文件上传发生错误"
		commonlib.OutputJson(w, m, " ")
		return
	}
}
