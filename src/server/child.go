// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-11-11 17:49
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-11-11 17:49 black 创建文档
package server

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
)

type Child struct {
	Cid      int
	Name     string
	CardId   string
	Pid      int
	Sex      int
	Birthday string
	Hobby    string
	CenterId int
	Avatar   string
}

func FindChildById(id string) (Child, error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select cid,name,card_id,pid,sex,birthday,hobby,center_id,avatar from child where cid=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		return Child{}, err
	}

	child := Child{}

	if rows.Next() {
		err = commonlib.PutRecord(rows, &child.Cid, &child.Name, &child.CardId, &child.Pid, &child.Sex, &child.Birthday, &child.Hobby, &child.CenterId, &child.Avatar)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return Child{}, err
		}
	}

	return child, nil
}
