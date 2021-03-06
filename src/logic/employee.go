// Title：员工相关服务
//
// Description:
//
// Author:Samurai
//
// Createtime:2013-09-24 13:04
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
package logic

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"strconv"
)

func CheckPwd(username, password string) (bool, lessgo.Employee, string) {

	var employee lessgo.Employee
	var dbPwd string

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select user_id,username,password,really_name,department_id from employee where username=?"
	rows, err := db.Query(sql, username)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, employee, " 数据库异常!"
	}
	if rows.Next() {
		err := rows.Scan(&employee.UserId, &employee.UserName,
			&dbPwd, &employee.ReallyName, &employee.DepartmentId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, employee, " 数据库异常!"
		}
	}

	if password != dbPwd {
		lessgo.Log.Warn("密码错误:", username, " : ", password)
		return false, employee, "密码错误"
	}

	sql = "select r.role_id,r.level,r.code from employee_role er left join role r on r.role_id=er.role_id where er.user_id=?"
	rows, err = db.Query(sql, employee.UserId)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return false, employee, " 数据库异常!"
	}
	var roleId, roleCode, roleLevel string
	for rows.Next() {

		err := commonlib.PutRecord(rows, &roleId, &roleLevel, &roleCode)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return false, employee, " 数据库异常!"
		}
		employee.RoleId += roleId + ","
		employee.RoleCode += roleCode + ","
		employee.RoleLevel += roleLevel + ","
	}
	lessgo.Log.Info(employee)
	return true, employee, ""
}

type Employee struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	ReallyName   string `json:"reallyName"`
	DepartmentId string `json:"departmentId"`
	CenterId     string `json:"centerId"`
}

//根据id获取员工信息
func FindEmployeeById(id int) (Employee, error) {

	var employee Employee

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select user_id,username,really_name,department_id,center_id from employee where user_id=?"

	rows, err := db.Query(sql, id)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return employee, err
	}

	if rows.Next() {
		err := commonlib.PutRecord(rows, &employee.UserId, &employee.UserName, &employee.ReallyName, &employee.DepartmentId, &employee.CenterId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return employee, err
		}
	}

	return employee, nil
}

//根据角色ID获取员工列表
func EmployeeListByRoleIdAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	id := r.FormValue("id")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select e.user_id,e.really_name from (select distinct(user_id) uid from employee_role where role_id =?) a left join employee e  on a.uid = e.user_id"

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	employees := []lessgo.Employee{}

	for rows.Next() {
		employee := lessgo.Employee{}

		err := commonlib.PutRecord(rows, &employee.UserId, &employee.ReallyName)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		employees = append(employees, employee)
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = employees

	commonlib.OutputJson(w, m, " ")
}

//获取中心下面没有离职的员工列表
func EmployeeListByCenterIdAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	id := r.FormValue("id")

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select e.user_id,e.really_name from employee e where e.center_id=? and e.is_leave=0 "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	employees := []lessgo.Employee{}

	for rows.Next() {
		employee := lessgo.Employee{}

		err := commonlib.PutRecord(rows, &employee.UserId, &employee.ReallyName)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		employees = append(employees, employee)
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = employees

	commonlib.OutputJson(w, m, " ")
}

//获取中心下面没有离职的员工列表
func EmployeeListInCenterAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	employee := lessgo.GetCurrentEmployee(r)

	if employee.UserId == "" {
		lessgo.Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	userId, _ := strconv.Atoi(employee.UserId)
	_employee, err := FindEmployeeById(userId)
	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select e.user_id,e.really_name from employee e where e.center_id=? and e.is_leave=0 "

	lessgo.Log.Debug(sql)

	rows, err := db.Query(sql, _employee.CenterId)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	employees := []lessgo.Employee{}

	for rows.Next() {
		employee := lessgo.Employee{}

		err := commonlib.PutRecord(rows, &employee.UserId, &employee.ReallyName)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		employees = append(employees, employee)
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = employees

	commonlib.OutputJson(w, m, " ")
}
