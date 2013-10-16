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
	"fmt"
	"time"
	"database/sql"
	"github.com/hjqhezgh/lessgo"
)

const (
	//视频和课程关联的状态, 1：未关联   2：关联成功   3：关联失败
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

//更新状态
func updateSatus(sid, vid, status int, db *sql.DB) {
	sql := "update video set schedule_detail_id=?,data_rel_status=? where vid=?"
	tx, err := db.Begin()

	_, err = tx.Exec(sql, sid, status, vid)
	if err != nil {
		lessgo.Log.Error(err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
}

func GetDateString() string {
	t := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d", t.Year(), int(t.Month()), t.Day())
	return str
}

//获取一周前的日期时间
func getWeekAgoDate(db *sql.DB) string{
	var dateString string
	sql := "select date_sub(now(),interval 1 week)"
	//查询
	rows, err := db.Query(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return	""
	}
	if rows.Next() {
		err := rows.Scan(&dateString)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return ""
		}
	}
	return dateString
}

func UpdateVideoStatus() {

	type Ret struct {
		id 			int
		rid         int
		cid			int
		begin_time	string
		end_time	string
	}

	flag := false
	var begin_time, end_time string
	var video_rets,class_rets []Ret
	var week_ago_date string

	sql_video := "select vid,rid,cid,start_time,end_time from video where data_rel_status=?";
	sql_class := "select id,room_id,center_id,real_start_time,real_end_time from class_schedule_detail where day_date>?";

	db := lessgo.GetMySQL()
	defer db.Close()
	//初始化时或日期变化时统一更新
	week_ago_date = getWeekAgoDate(db)

	//一周前的视频没关联更新状态为失败
	sql := "update video set data_rel_status=? where data_rel_status=? and start_time<?"
	tx, err := db.Begin()
	_, err = tx.Exec(sql, CLASS_STATUS_FAILED, CLASS_STATUS_INIT, week_ago_date)
	if err != nil {
		lessgo.Log.Error(err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()

	//查询未关联的视频
	rows, err := db.Query(sql_video, CLASS_STATUS_INIT)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}
	for rows.Next() {
		ret := new(Ret)
		err := rows.Scan(&ret.id, &ret.rid, &ret.cid ,&begin_time, &end_time)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}
		ret.begin_time = TimeString(begin_time)
		ret.end_time = TimeString(end_time)
		video_rets = append(video_rets, *ret)
		flag = true
	}
	fmt.Println(video_rets)
	//flag 表示有未登记的video文件，查询课程安排表关联
	if flag {
		rows, err = db.Query(sql_class, TimeString(week_ago_date))
		if err != nil {
			lessgo.Log.Error(err.Error())
			return
		}
		for rows.Next() {
			ret := new(Ret)
			err := rows.Scan(&ret.id, &ret.rid, &ret.cid ,&begin_time, &end_time)
			if err != nil {
				lessgo.Log.Error(err.Error())
				return
			}
			ret.begin_time = TimeString(begin_time)
			ret.end_time = TimeString(end_time)
			class_rets = append(class_rets, *ret)
		}
		fmt.Println(class_rets)
		//关联更新
		for _, v := range(video_rets) {
			for _, c := range(class_rets) {
				if v.cid==c.cid && v.rid==c.rid {
					if (c.begin_time<v.begin_time && v.end_time<c.end_time) 										||
							(v.begin_time<c.begin_time && c.end_time<v.end_time)									||
								(c.begin_time>v.begin_time && v.end_time<c.end_time && v.end_time>c.begin_time)		||
									(c.begin_time<v.begin_time && v.end_time>c.end_time && v.begin_time<c.end_time) {
						lessgo.Log.Debug("新视频id：", v.id, "关联的课程时间id：", c.id)
						updateSatus (c.id, v.id, CLASS_STATUS_SUCCESS, db)
					}
				}
			}
		}
	}
}
