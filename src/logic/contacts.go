// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 17:41
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 17:41 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
)

func insertContacts(tx *sql.Tx, name, phone, consumerId string) (id int64, err error) {

	sql := "insert into contacts(name,phone,is_default,consumer_id) values(?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(name, phone, "1", consumerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	contactsId, err := res.LastInsertId()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return contactsId, nil
}
