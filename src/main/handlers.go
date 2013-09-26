// Title：请求处理器配置文件
//
// Description:
//
// Author:black
//
// Createtime:2013-07-30 10:13
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-07-30 10:13 black 创建文档
package main

import (
	"net/http"
	"server"
)

//URL映射列表
var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/login": server.LoginAction,

	//音频相关服务
	"/consultant_phone_list.json": server.ConsultantPhoneListAction,
}
