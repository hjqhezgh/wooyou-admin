// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-15 21:10
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-15 21:10 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
	"time"
)

func insertConsumerStatusLog(tx *sql.Tx, consumerId, employeeId, oldStatus, newStatus string) (id int64, err error) {

	sql := "insert into consumer_status_log(consumer_id,employee_id,create_time,old_status,new_status) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(consumerId, employeeId, time.Now().Format("20060102150405"), oldStatus, newStatus)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	logId, err := res.LastInsertId()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return logId, err
}
