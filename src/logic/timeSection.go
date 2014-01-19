// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-17 17:02
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-17 17:02 black 创建文档
package logic

import (
	"github.com/hjqhezgh/lessgo"
)

func getTimeSectionById(id string) (map[string]string, error) {

	sql := `select id,center_id,start_time,end_time,lesson_no from time_section where id=?`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	var dataMap map[string]string

	if rows.Next() {
		dataMap, err = lessgo.GetDataMap(rows)
	}

	if err != nil {
		return nil, err
	}

	return dataMap, nil
}
