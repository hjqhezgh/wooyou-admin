///**
// * Title：
// *
// * Description:
// *
// * Author: Ivan
// *
// * Create Time: 2013-09-26 11:52
// *
// * Version: 1.0
// *
// * 修改历史: 版本号 修改日期 修改人 修改说明
// *   1.0 2013-09-26 Ivan 创建文件
//*/
package server

//
//import (
//	"net/http"
//	"github.com/hjqhezgh/commonlib"
//	"github.com/hjqhezgh/lessgo"
//	"time"
//	"strings"
//)
//
//func EmployeeLoad(userId string)(*Employee, string, error){
//	return GetEmployeeById(userId)
//}
//
//func EmployeeLoadAction(w http.ResponseWriter, r *http.Request) {
//	m := make(map[string]interface {})
//
//	userId := r.FormValue("id")
//	employee, msg, err := EmployeeLoad(userId)
//	if err != nil {
//		m["success"] = false
//		m["msg"] = msg
//		commonlib.OutputJson(w, m, "")
//		return
//	}
//
//	loadFormObjects := []lessgo.LoadFormObject{}
//	datas := []*LoadJson{}
//	jsonUserId := new(LoadJson)
//	jsonUserId.Field = "user_id"
//	jsonUserId.Value = employee.UserId
//	datas = append(datas, jsonUserId)
//
//	jsonUsername := new(LoadJson)
//	jsonUsername.Field = "username"
//	jsonUsername.Value = employee.UserName
//	datas = append(datas, jsonUsername)
//
//	jsonPassword := new(LoadJson)
//	jsonPassword.Field = "password"
//	jsonPassword.Value = employee.Password
//	datas = append(datas, jsonPassword)
//
//	jsonRealName := new(LoadJson)
//	jsonRealName.Field = "realname"
//	jsonRealName.Value = employee.ReallyName
//	datas = append(datas, jsonRealName)
//
//	jsonRole := new(LoadJson)
//	jsonRole.Field = "role"
//	jsonRole.Value = employee.Role
//	datas = append(datas, jsonRole)
//
//	if err != nil {
//		m["sucess"] = false
//		m["msg"] = msg
//	} else {
//		m["success"] = true
//		m["datas"] = datas
//	}
//	commonlib.OutputJson(w, m, "")
//
//	return
//}
//
//func EmployeeSave(username, password, realName, role string)(*Employee, string, error){
//	db := lessgo.GetMySQL()
//	defer db.Close()
//
//	// 开启事务
//	tx, err := db.Begin()
//	if err != nil {
//		lessgo.Log.Error("db.Begin: ", err.Error())
//		return nil, "开启事务时，数据库异常", err
//	}
//	defer func() {
//		if err != nil && tx != nil {
//			// 回滚
//			if rbErr := tx.Rollback(); rbErr != nil {
//				lessgo.Log.Error("tx.Rollback: ", rbErr.Error())
//				return
//			}
//		}
//	}()
//	t := time.Now()
//	// 生成用户基本信息
//	sqlEmployee := "insert into employee (username, password, really_name, department_id) values (?,?,?,?);"
//	res, err := TxInsert(tx, sqlEmployee, username, password, realName, "1")
//	if err != nil { return nil, "插入用户基本数据，数据库异常", err }
//	// 获取生成用户基本信息ID
//	userId, err := res.LastInsertId()
//	if err != nil { return nil, "获取生成用户基本信息ID失败!", err}
//	// 保存用户角色信息
//	sqlEmployeeRole := "insert into employee_role(user_id, role_id) values (?,?);"
//	if role != "" {
//		roleArr := strings.Split(role, ",")
//		for _, roleId := range roleArr {
//			_, err = TxInsert(tx, sqlEmployeeRole, userId, roleId)
//			if err != nil { return nil, "插入用户角色信息，数据库异常", err }
//		}
//	}
//	lessgo.Log.Info("保存用户信息时间: ", time.Now().Sub(t))
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		lessgo.Log.Error("tx.Commit: ", err.Error())
//		return nil, "提交事务，数据库异常", err
//	}
//	// 关闭数据库连接
//	if err = db.Close(); err != nil {
//		lessgo.Log.Error("db.Close: ", err.Error())
//		return nil, "关闭数据库连接，数据库异常", err
//	}
//	// 返回当前插入用户信息
//	employee, _, err := GetEmployeeById(string(userId))
//
//	return employee, "处理成功", nil
//}
//
//func EmployeeSaveAction(w http.ResponseWriter, r *http.Request) {
//	EmployeeFormAction(w, r, FORM_ACTION_TYPE.CREATE.Key)
//}
//
//func EmployeeUpdate(userId, username, password, realName, role string)(*Employee, string, error){
//	db := lessgo.GetMySQL()
//	defer db.Close()
//
//	// 开启事务
//	tx, err := db.Begin()
//	if err != nil {
//		lessgo.Log.Error("db.Begin: ", err.Error())
//		return nil, "开启事务时，数据库异常", err
//	}
//	defer func() {
//		if err != nil && tx != nil {
//			// 回滚
//			if rbErr := tx.Rollback(); rbErr != nil {
//				lessgo.Log.Error("tx.Rollback: ", rbErr.Error())
//				return
//			}
//		}
//	}()
//	t := time.Now()
//	// 更新用户基本信息
//	sqlEmployee := "update employee set username=?, password=?, really_name=? where user_id=?;"
//	_, err = TxUpdate(tx, sqlEmployee, username, password, realName, userId)
//	if err != nil { return nil, "更新用户基本数据，数据库异常", err }
//	// 更新用户角色信息
//	// 删除用户原来的角色信息
//	sqlDelER := "delete from employee_role where user_id=?;"
//	_, err = TxDelete(tx, sqlDelER, userId)
//	if err != nil { return nil, "删除用户已有角色信息，数据库异常", err }
//	// 保存新角色信息
//	sqlInsEr := "insert into employee_role(user_id, role_id) values (?,?);"
//	if role != "" {
//		roleArr := strings.Split(role, ",")
//		for _, roleId := range roleArr {
//			_, err = TxInsert(tx, sqlInsEr, userId, roleId)
//			if err != nil { return nil, "插入用户新角色信息，数据库异常", err }
//		}
//	}
//	lessgo.Log.Info("更新用户信息时间: ", time.Now().Sub(t))
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		lessgo.Log.Error("tx.Commit: ", err.Error())
//		return nil, "提交事务，数据库异常", err
//	}
//	// 关闭数据库连接
//	if err = db.Close(); err != nil {
//		lessgo.Log.Error("db.Close: ", err.Error())
//		return nil, "关闭数据库连接，数据库异常", err
//	}
//	// 返回当前插入用户信息
//	employee, _, err := GetEmployeeById(string(userId))
//
//	return employee, "处理成功", nil
//}
//
//func EmployeeUpdateAction(w http.ResponseWriter, r *http.Request) {
//	EmployeeFormAction(w, r, FORM_ACTION_TYPE.UPDATE.Key)
//}
//
//func EmployeeDelete(userId string)(string, error){
//	db := lessgo.GetMySQL()
//	defer db.Close()
//
//	// 开启事务
//	tx, err := db.Begin()
//	if err != nil {
//		lessgo.Log.Error("db.Begin: ", err.Error())
//		return "开启事务时，数据库异常", err
//	}
//	defer func() {
//		if err != nil && tx != nil {
//			// 回滚
//			if rbErr := tx.Rollback(); rbErr != nil {
//				lessgo.Log.Error("tx.Rollback: ", rbErr.Error())
//				return
//			}
//		}
//	}()
//	t := time.Now()
//	// 删除用户基本信息
//	sqlDelEmp := "delete from employee where user_id=?;"
//	_, err = TxDelete(tx, sqlDelEmp, userId)
//	if err != nil { return "删除用户基本数据，数据库异常", err }
//	// 删除用户角色信息
//	sqlDelER := "delete from employee_role where user_id=?;"
//	_, err = TxDelete(tx, sqlDelER, userId)
//	if err != nil { return "删除用户已有角色信息，数据库异常", err }
//	lessgo.Log.Info("删除用户信息时间: ", time.Now().Sub(t))
//	// 提交事务
//	if err = tx.Commit(); err != nil {
//		lessgo.Log.Error("tx.Commit: ", err.Error())
//		return "提交用户信息事务，数据库异常", err
//	}
//	// 关闭数据库连接
//	if err = db.Close(); err != nil {
//		lessgo.Log.Error("db.Close: ", err.Error())
//		return "关闭数据库连接，数据库异常", err
//	}
//
//	return "删除成功", nil
//}
//
//func EmployeeDeleteAction(w http.ResponseWriter, r *http.Request) {
//	lessgo.Log.Debug("delete")
//	m := make(map[string]interface {})
//	err := r.ParseForm()
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		return
//	}
//	userId := r.FormValue("id")
//	lessgo.Log.Debug(userId)
//	msg, err := EmployeeDelete(userId)
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		m["sucess"] = false
//	} else {
//		m["success"] = true
//	}
//	m["msg"] = msg
//	commonlib.RenderTemplate(w, r, "notify_delete.html", m, nil, "../template/component/" + getTerminal(r.URL.Path) + "/notify_delete.html")
//
//	return
//}
//
//func EmployeeFormAction(w http.ResponseWriter, r *http.Request, actionType string) {
//	m := make(map[string]interface {})
//
//	err := r.ParseForm()
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		return
//	}
//
//	userId := r.FormValue("user_id")
//	username := r.FormValue("username")
//	password := r.FormValue("password")
//	realName := r.FormValue("realname")
//	role := r.FormValue("role")
//
//	employee := new(Employee)
//	msg := ""
//	switch actionType {
//	case FORM_ACTION_TYPE.CREATE.Key:
//		if CheckUserName(username) {
//			lessgo.Log.Debug("exist")
//			m["sucess"] = false
//			m["msg"] = "用户名已存在，请更改用户名。"
//			commonlib.OutputJson(w, m, "")
//			return
//		}
//		employee, msg, err = EmployeeSave(username, password, realName, role)
//	case FORM_ACTION_TYPE.UPDATE.Key:
//		employee, msg, err = EmployeeUpdate(userId, username, password, realName, role)
//	}
//
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		m["sucess"] = false
//		m["msg"] = msg
//	} else {
//		m["success"] = true
//		m["datas"] = employee
//	}
//	commonlib.OutputJson(w, m, "")
//
//	return
//}
//
//func GetEmployeeById(userId string) (*Employee, string, error) {
//	sql := "select e.password, e.username, e.really_name, e.department_id, GROUP_CONCAT(er.role_id SEPARATOR ',') as role from employee e left join employee_role er on e.user_id=er.user_id where e.user_id=?;"
//	lessgo.Log.Debug(sql, userId)
//	rows, err := DBSelect(sql, userId)
//
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		return nil, "查询用户信息出错", err
//	}
//
//	employee := new(Employee)
//	employee.UserId = userId
//	if rows.Next() {
//		err = commonlib.PutRecord(rows, &employee.Password, &employee.UserName, &employee.ReallyName, &employee.DepartmentId, &employee.Role)
//		if err != nil {
//			lessgo.Log.Error(err.Error())
//			return nil, "封装用户信息出错", err
//		}
//	}
//	lessgo.Log.Debug(employee)
//
//	return employee, "处理成功", nil
//}
//
////func GetEmployeeGridPanelAction(w http.ResponseWriter, r *http.Request) {
////	data := make(map[string]interface {})
////	lessgo.Log.Info("GetEmployeeGridPanelAction")
////	commonlib.OutputJson(w, data,"")
////
////	return
////}
//
//func CheckUserName(username string) bool {
//	sql := "select user_id from employee where username=?;"
//	lessgo.Log.Debug("CheckUserName: ", strings.Trim(username, " "))
//	rows, err := DBSelect(sql, strings.Trim(username, " "))
//	if err != nil {
//		lessgo.Log.Error(err.Error())
//		return false
//	}
//
//	if rows.Next() {
//		return true
//	}
//
//	return false
//}
//
