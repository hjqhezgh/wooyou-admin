/**
 * Title：
 * 
 * Description: 
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-28 10:10
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-28 Ivan 创建文件
*/
package server

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
)

/****
* 数据插入
 */
func TxInsert(tx *sql.Tx, sql string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sql, args ...)
}

/****
* 数据删除
 */
func TxDelete(tx *sql.Tx, sql string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sql, args ...)
}

/****
* 数据更新
 */
func TxUpdate(tx *sql.Tx, sql string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sql, args ...)
}

/****
* 数据查询
 */
func TxSelect(tx *sql.Tx, sql string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		lessgo.Log.Error("tx.Prepare: ", err.Error())
		return nil, err
	}
	defer func() {
		if stmtErr := stmt.Close(); stmtErr != nil {
			lessgo.Log.Error("stmt.Close: ", stmtErr.Error())
			return
		}
	}()

	rows, err := stmt.Query(args...)

	if err != nil { return nil, err }

	return rows, err
}

/****
* 数据插入
 */
func DBInsert(sql string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sql, args ...)
}

/****
* 数据删除
 */
func DBDelete(sql string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sql, args ...)
}

/****
* 数据更新
 */
func DBUpdate(sql string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sql, args ...)
}

/****
* 数据查询
 */
func DBSelect(sql string, args ...interface{}) (*sql.Rows, error) {
	//log.Debug(sql, "  args: ", args)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sql, args ...)

	return rows, err
}

/****
* 数据增删改处理
 */
func dbOperation(sql string, args ...interface{}) (sql.Result, error) {
	//	log.Debug(sql, "  args: ", args)
	db := lessgo.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)

	if err != nil { return nil, err }

	result, err := stmt.Exec(args...)

	if err != nil { return nil, err }

	return result, err
}
/****
* 数据增删改处理(事务)
 */
func txOperation(tx *sql.Tx, sqlStr string, args ...interface{}) (sql.Result, error) {
	//	log.Debug(sqlStr, "  args: ", args)
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		lessgo.Log.Error("tx.Prepare: ", err.Error())
		return nil, err
	}
	defer func() {
		if stmtErr := stmt.Close(); stmtErr != nil {
			lessgo.Log.Error("stmt.Close: ", stmtErr.Error())
			return
		}
	}()

	result, err := stmt.Exec(args...)

	if err != nil { return nil, err }

	return result, err
}

