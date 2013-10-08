// Title：音频下载相关服务
//
// Description:
//
// Author:black
//
// Createtime:2013-09-29 09:56
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-09-29 09:56 black 创建文档
package server

import (
	"bufio"
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//顾问分页数据服务
func DownloadAudioAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	idString := r.FormValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		lessgo.Log.Error("录音下载id格式错误：", idString)
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, "")
	}

	audio, err := FindAudioById(id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, "")
		return
	}

	date := strings.Split(audio.StartTime, " ")[0]

	sendTime, err := time.ParseInLocation("2006-1-2", date, time.Local)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, "")
		return
	}

	audioRootPath, _ := lessgo.Config.GetValue("wooyou", "audioRootPath")

	filePath := fmt.Sprint(audioRootPath, "/", audio.Cid, "/", sendTime.Format("2006-1-2"), "/", audio.Filename)

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, "")
		return
	}

	datas, err := ioutil.ReadAll(bufio.NewReader(f))

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, "")
		return
	}

	w.Header().Set("mimetype", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename="+audio.Filename)

	w.Write(datas)
}
