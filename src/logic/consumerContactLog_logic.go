// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 17:46
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 17:46 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
	"time"
)

func insertConsumerContactsLog(tx *sql.Tx, createUser, note, consumerId,contactsType string) (id int64, err error) {

	sql := "insert into consumer_contact_log(create_user,create_time,note,consumer_id,type) values(?,?,?,?,?) "
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(createUser, time.Now().Format("20060102150405"), note, consumerId, contactsType)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	consumerContactsLogId, err := res.LastInsertId()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return consumerContactsLogId, nil
}
