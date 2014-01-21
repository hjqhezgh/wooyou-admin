// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-20 10:16
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-20 10:16 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
	"time"
)

func insertPayLog(tx *sql.Tx, consumerId, employeeId string) error {

	sql := "insert into pay_log(consumer_id,pay_time,employee_id) values(?,?,?)"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(consumerId, time.Now().Format("20060102150405"), employeeId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}
