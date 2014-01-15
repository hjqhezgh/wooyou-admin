// Title：课件相关服务
//
// Description:
//
// Author:black
//
// Createtime:2014-01-07 13:41
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-07 13:41 black 创建文档
package logic

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"time"
)

//获取课件分页列表
func CoursewareList(centerId, kw string, pageNo, pageSize int) (*commonlib.TraditionPage, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	countSql := "select count(1) from courseware where center_id=? and (name like ? or intro like ?) "
	lessgo.Log.Debug(countSql)
	countParams := []interface{}{centerId, "%" + kw + "%", "%" + kw + "%"}

	totalPage, totalNum, err := lessgo.GetTotalPage(pageSize, db, countSql, countParams)

	if err != nil {
		return nil, err
	}

	currPageNo := pageNo
	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	dataSql := "select id,name,intro from courseware where center_id=? and (name like ? or intro like ?) order by id desc limit ?,?"
	lessgo.Log.Debug(dataSql)

	dataParams := []interface{}{}
	dataParams = append(dataParams, centerId)
	dataParams = append(dataParams, "%"+kw+"%")
	dataParams = append(dataParams, "%"+kw+"%")
	dataParams = append(dataParams, (currPageNo-1)*pageSize)
	dataParams = append(dataParams, pageSize)

	pageData, err := lessgo.GetFillObjectPage(db, dataSql, currPageNo, pageSize, totalNum, dataParams)

	if err != nil {
		return nil, err
	}

	return pageData, nil
}

//课件保存
func SaveCourseware(id, centerId, name, intro string) (flag bool, msg string, err error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, "", err
	}

	if id == "" {
		sql := "insert into courseware(center_id,name,create_time,intro) values(?,?,?,?)"
		stmt, err := tx.Prepare(sql)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		_, err = stmt.Exec(centerId, name, time.Now().Format("20060102150405"), intro)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return false, "", err
		}
	} else {
		sql := "update courseware set name=?,intro=? where id=? "
		stmt, err := tx.Prepare(sql)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, "", err
		}

		_, err = stmt.Exec(name, intro, id)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return false, "", err
		}
	}

	tx.Commit()

	return true, "", nil
}

//修改表单数据读取
func LoadCourseware(id string) (loadFormObjects []lessgo.LoadFormObject, err error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select id,name,intro from courseware where id=? "
	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	var name, intro string

	if rows.Next() {
		err := commonlib.PutRecord(rows, &id, &name, &intro)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			return nil, err
		}
	}

	h1 := lessgo.LoadFormObject{"id", id}
	h2 := lessgo.LoadFormObject{"name", name}
	h3 := lessgo.LoadFormObject{"intro", intro}

	loadFormObjects = append(loadFormObjects, h1)
	loadFormObjects = append(loadFormObjects, h2)
	loadFormObjects = append(loadFormObjects, h3)

	return loadFormObjects, nil
}
