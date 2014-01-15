// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-15 21:19
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-15 21:19 black 创建文档
package logic

import (
	"github.com/hjqhezgh/lessgo"
)

/*
select class_id id,name,start_time,end_time,code,center_id,child_num from wyclass where class_id=?
*/
func findWyClassById(id string) (map[string]string, error) {

	sql := `
			select class_id id,name,start_time,end_time,code,center_id,child_num from wyclass where class_id=?
	    	`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	dataMap, err := lessgo.GetDataMap(rows)

	if err != nil {
		return nil, err
	}

	return dataMap, nil
}
