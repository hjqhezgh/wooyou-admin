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

func TimeSectionPage(centerId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from time_section where center_id=? "
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

	dataSql := `select id,start_time,end_time from time_section where center_id=? order by start_time,end_time limit ?,?`
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

func InsertTimeSection(centerId, startTime, endTime string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	_, err = insertTimeSection(tx, centerId, startTime, endTime)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}

func insertTimeSection(tx *sql.Tx, centerId, startTime, endTime string) (id int64, err error) {

	sql := "insert into time_section(center_id,start_time,end_time,lesson_no) values(?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(centerId,startTime,endTime,"1")

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newTimeSectionId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newTimeSectionId, nil
}
