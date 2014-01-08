// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-08 11:31
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-08 11:31 black 创建文档
package logic

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
	"strconv"
	"fmt"
)

func insertChild(tx *sql.Tx, name, pid, sex, birthday, year, month, centerId string) (id int64, err error) {

	sql := "insert into child(name,pid,sex,birthday,center_id) values(?,?,?,?,?)"
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	childBirthday := "20090101"
	if birthday != "" {
		childBirthday = birthday
	} else {
		if year != "" && month != "" {
			monthInt,_ := strconv.Atoi(month)
			if monthInt > 9 {
				childBirthday = year + month + "01"
			}else{
				childBirthday = year + "0" +month + "01"
			}
		} else if year != "" {
			childBirthday = year + "0101"
		}
	}

	res, err := stmt.Exec(name, pid, sex, childBirthday, centerId)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	childId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return childId, nil
}


/*
select cid id,name,card_id cardId,pid,sex,birthday,center_id centerId,avatar
	        from child where cid=?
*/
func findChildById(id string) (map[string]string,error)  {

	sql := `select cid id,name,card_id cardId,pid,sex,birthday,center_id centerId,avatar
	        from child where cid=?`

	lessgo.Log.Debug(sql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	dataMap,err := lessgo.GetDataMap(rows)

	if err != nil {
		return nil, err
	}

	return dataMap,nil
}

func updateChild(tx *sql.Tx,childDataMap map[string]interface {},id string) error{
	sql := "update child set %v where cid=?"
	params := []interface {}{}

	setSql := ""

	for key,value := range childDataMap{
		setSql += key+"=?,"
		params = append(params,value)
	}

	params = append(params,id)

	setSql = commonlib.Substr(setSql,0,len(setSql)-1)

	sql = fmt.Sprintf(sql,setSql)
	lessgo.Log.Debug(sql)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(params...)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return err
	}

	return nil
}

//todo 目前只回去第一个孩子，将来这块的逻辑有待优化
func getChildByParentId(pid string) (int64, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select cid from child where pid=?"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, pid)

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
