// Title：音频服务
//
// Description:
//
// Author:black
//
// Createtime:2013-09-29 10:02
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-09-29 10:02 black 创建文档
package server

import (
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
)

type Audio struct {
	Aid         int    `json:"aid"`
	Cid         int    `json:"cid"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Filename    string `json:"filename"`
	Seconds     string `json:"seconds"`
	LocalPhone  string `json:"localphone"`
	RemotePhone string `json:"remotephone"`
	Inout       string `json:"inout"`
}

//根据id查找视频
func FindAudioById(id int) (*Audio, error) {

	baseSql := "select aid,cid,start_time,end_time,filename,localphone,remotephone,seconds,`inout` from audio where aid=? "

	lessgo.Log.Debug(baseSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(baseSql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	audio := new(Audio)

	if rows.Next() {
		err := commonlib.PutRecord(rows, &audio.Aid, &audio.Cid, &audio.StartTime, &audio.EndTime, &audio.Filename, &audio.LocalPhone, &audio.RemotePhone, &audio.Seconds, &audio.Inout)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, err
		}
	}

	return audio, nil
}


