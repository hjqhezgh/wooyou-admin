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
package tool

import (
	"database/sql"
	"github.com/hjqhezgh/lessgo"
	"time"
)

/****
* 数据插入
 */
func TxInsert(tx *sql.Tx, sqlStr string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sqlStr, args ...)
}

/****
* 数据删除
 */
func TxDelete(tx *sql.Tx, sqlStr string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sqlStr, args ...)
}

/****
* 数据更新
 */
func TxUpdate(tx *sql.Tx, sqlStr string, args ...interface {}) (sql.Result, error) {
	return txOperation(tx, sqlStr, args ...)
}

/****
* 数据查询
 */
func TxSelect(tx *sql.Tx, sqlStr string, args ...interface{}) (*sql.Rows, error) {
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

	rows, err := stmt.Query(args...)

	if err != nil { return nil, err }

	return rows, err
}

/****
* 数据插入
 */
func DBInsert(sqlStr string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sqlStr, args ...)
}

/****
* 数据删除
 */
func DBDelete(sqlStr string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sqlStr, args ...)
}

/****
* 数据更新
 */
func DBUpdate(sqlStr string, args ...interface {}) (sql.Result, error) {
	return dbOperation(sqlStr, args ...)
}

/****
* 数据查询
 */
func DBSelect(sqlStr string, args ...interface{}) (*sql.Rows, error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(sqlStr, args ...)

	return rows, err
}

/****
* 数据增删改处理
 */
func dbOperation(sqlStr string, args ...interface{}) (sql.Result, error) {
	db := lessgo.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sqlStr)

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
// 事务处理
func TransactionAction(txAction func(*sql.Tx)(string, error))(string, error){
	db := lessgo.GetMySQL()
	defer db.Close()

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		lessgo.Log.Error("db.Begin: ", err.Error())
		return "开启事务时，数据库异常", err
	}
	defer func() {
		if err != nil && tx != nil {
			// 回滚
			if rbErr := tx.Rollback(); rbErr != nil {
				lessgo.Log.Error("tx.Rollback: ", rbErr.Error())
				return
			}
		}
	}()
	t := time.Now()
	msg, err := txAction(tx)
	if err != nil {
		return msg, err
	}
	lessgo.Log.Info("事务处理时间: ", time.Now().Sub(t))

	// 提交事务
	if err = tx.Commit(); err != nil {
		lessgo.Log.Error("tx.Commit: ", err.Error())
		return "提交事务，数据库异常", err
	}
	// 关闭数据库连接
	if err = db.Close(); err != nil {
		lessgo.Log.Error("db.Close: ", err.Error())
		return "关闭数据库连接，数据库异常", err
	}

	return msg, nil
}

