// Title：应用入口
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 13:04
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 13:04 black 创建文档
package main

import (
	"fmt"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"server"
	"strconv"
)

func main() {

	r := lessgo.ConfigLessgo()

	portString, _ := lessgo.Config.GetValue("lessgo", "port")

	port, _ := strconv.Atoi(portString)

	for url, handler := range handlers {
		r.HandleFunc(url, handler)
	}

	http.Handle("/", r)

	fmt.Println("服务器监听", portString, "端口")

	lessgo.Log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

//URL映射列表
var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/login": server.LoginAction,

	//音频相关服务
	"/consultant_phone_list.json":        server.ConsultantPhoneListAction,
	"/consultant_phone_detail_list.json": server.ConsultantPhoneDetailListAction,
	"/downloadAudio":                     server.DownloadAudioAction,
	"/audioNoteLoad.json":                server.AudioNoteLoadAction,
	"/audioNoteSave.json":                server.AudioNoteSaveAction,

	//客户相关服务
	"/consumer.json":     server.ConsumerListAction,
	"/consumerSave.json": server.ConsumerSaveAction,
	"/consumerLoad.json": server.ConsumerLoadAction,

	//Call Center
	"/callCenterStatistics.json": server.CallCenterStatisticsAction,
}
