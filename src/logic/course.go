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
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
)

func InsertCourse(centerId, courseName, courseType string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	courseNameFlag, err := checkCourseNameExist(centerId, courseName)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if !courseNameFlag {
		_, err = insertCourse(tx, centerId, courseName, courseType)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}
	} else {
		return false, "课程名字已存在", nil
	}

	tx.Commit()

	return true, "", nil
}

func UpdateCourse(id, courseName, courseType string) (flag bool, msg string, err error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	updateCourseMap := make(map[string]interface{})
	updateCourseMap["name"] = courseName
	updateCourseMap["type"] = courseType

	err = updateCourse(tx, updateCourseMap, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	tx.Commit()

	return true, "", nil
}

func insertCourse(tx *sql.Tx, centerId, courseName, courseType string) (id int64, err error) {

	sql := "insert into course(name,center_id,price,is_probation,type,begin_age,end_age,intro,app_display_level,create_time,lesson_num) values(?,?,?,?,?,?,?,?,?,?,?)"
	lessgo.Log.Debug(sql)
	stmt, err := tx.Prepare(sql)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(courseName, centerId, "0", "0", courseType, "0", "60", courseName, "0", time.Now().Format("20060102150405"), "12")

	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	newCourseId, err := res.LastInsertId()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return 0, err
	}

	return newCourseId, nil
}

func checkCourseNameExist(centerId, courseName string) (bool, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	dataSql := ""

	dataSql = "select count(1) from course where center_id=? and name=?"

	lessgo.Log.Debug(dataSql)

	var rows *sql.Rows
	var err error

	rows, err = db.Query(dataSql, centerId, courseName)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, err
	}

	num := 0

	if rows.Next() {

		err = commonlib.PutRecord(rows, &num)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, err
		}
	}

	if num > 0 {
		return true, nil
	}

	return false, nil
}

func updateCourse(tx *sql.Tx, courseDataMap map[string]interface{}, id string) error {
	sql := "update course set %v where cid=?"
	params := []interface{}{}

	setSql := ""

	for key, value := range courseDataMap {
		setSql += key + "=?,"
		params = append(params, value)
	}

	params = append(params, id)

	setSql = commonlib.Substr(setSql, 0, len(setSql)-1)

	sql = fmt.Sprintf(sql, setSql)
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

func CoursePage(centerId string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := " select count(1) from course where center_id=? "
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

	dataSql := `select cid as id,name,type from course where center_id=? order by cid desc limit ?,?`
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
