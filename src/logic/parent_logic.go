// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 17:53
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 17:53 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
)

//根据电话获取parent表的id
func getParentIdByPhone(phone string) (int64, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select pid from parent where telephone=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, phone)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	var id int64

	if rows.Next() {

		err = commonlib.PutRecord(rows, &id)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return 0, err
		}
	}

	return id, nil
}

func insertParent(tx *sql.Tx, name, password, telephone, comeForm string) (id int64, err error) {

	sql := "insert into parent(name,password,telephone,reg_date,come_form) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(name, password, telephone, time.Now().Format("20060102150405"), comeForm)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newParentId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newParentId, nil
}
