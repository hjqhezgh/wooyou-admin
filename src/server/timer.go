// Title：
//
// Description:
//
// Author:Samurai
//
// Createtime:2013-10-14 15:16
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-10-14 15:16 Samurai 创建文档package timer


package server

import (
	"time"
	"fmt"
	"database/sql"
	"github.com/hjqhezgh/lessgo"
)

const (
	CLASS_STATUS_INIT = 1
	CLASS_STATUS_SUCCESS = 2
	CLASS_STATUS_FAILED = 3
)
//时间串转换 2013-09-25 18:58:42 ---> 20130925185842
func TimeString(time string) string{
	var bb []byte
	b := []byte(time)
	for _, v := range(b) {
		if v >= 48 && v <= 57 {
			bb = append(bb, v)
		}
	}
	return string(bb)
}


func updateSatus(sid, vid, status int, db *sql.DB) {
	sql := "update video set schedule_detail_id=?,data_rel_status=? where vid=?"
	tx, err := db.Begin()

	_, err = tx.Exec(sql, sid, status, vid)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()
}

func VedioTimer() {

	type Ret struct {
		id 		int
		Btime	string
		Etime	string
	}

	var i int
	var Btime, Etime string
	var video_rets,class_rets []Ret

	sqlVideo := "select vid,start_time,end_time from video where data_rel_status=?";
	sqlClass := "select id,real_start_time,real_end_time from class_schedule_detail";

	for {

		i = 0

		db := lessgo.GetMySQL()
		defer db.Close()

		video_rets = []Ret{}
		class_rets = []Ret{}

		//查询
		rows, err := db.Query(sqlVideo, CLASS_STATUS_INIT)
		if err != nil {
			lessgo.Log.Warn(err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		for rows.Next() {
			vRet := new(Ret)
			err := rows.Scan(&vRet.id, &Btime, &Etime)
			if err != nil {
			}
			vRet.Btime = TimeString(Btime)
			vRet.Etime = TimeString(Etime)
			video_rets = append(video_rets, *vRet)
			i++
		}

		//i >0 表示有未登记的video文件，查询课程安排表关联
		if i > 0 {
			rows, err = db.Query(sqlClass)
			if err != nil {
				lessgo.Log.Warn(err.Error())
				time.Sleep(time.Second * 10)
				continue
			}
			for rows.Next() {
				cRet := new(Ret)
				err := rows.Scan(&cRet.id, &Btime, &Etime)
				if err != nil {
				}
				cRet.Btime = TimeString(Btime)
				cRet.Etime = TimeString(Etime)
				class_rets = append(class_rets, *cRet)
			}

			//关联更新
			for _, v := range(video_rets) {
				for _, c := range(class_rets) {
					if (c.Btime<v.Btime && v.Etime<c.Etime) 									||
							(v.Btime<c.Btime && c.Etime<v.Etime)								||
								(c.Btime>v.Btime && v.Etime<c.Etime && v.Etime>c.Btime) 	 	||
									(c.Btime<v.Btime && v.Etime>c.Etime && v.Btime<c.Etime) {
						fmt.Println("关联到 视频id：", v.id, "课程时间id：", c.id)
						updateSatus (c.id, v.id, CLASS_STATUS_SUCCESS, db)
					}
				}
			}
		}

		db.Close()
		//更新
		time.Sleep(time.Second * 60)
	}
}
