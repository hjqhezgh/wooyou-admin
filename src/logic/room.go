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
	"github.com/hjqhezgh/commonlib"
	"database/sql"
)

func RoomPage(centerId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from room where cid=? "
	lessgo.Log.Debug(countSql)
	countParams := []interface{}{centerId}

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := `select rid as id,name,channel from room where cid=? limit ?,?`
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{}
	dataParams = append(dataParams, centerId)
	dataParams = append(dataParams, (currPageNo-1)*pageSize)
	dataParams = append(dataParams, pageSize)

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

func InsertRoom(centerId, name, channel string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	_, err = insertRoom(tx, centerId, name, channel)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}

func insertRoom(tx *sql.Tx, centerId, name, channel string) (id int64, err error) {

	sql := "insert into room(cid,name,channel) values(?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(centerId,name,channel)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newRoomId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newRoomId, nil
}
