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
package server

import (
	"github.com/hjqhezgh/lessgo"
	"github.com/hjqhezgh/commonlib"
)

func CheckPwd(username, password string) (bool, lessgo.Employee, string){

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

	return true, employee, ""
}

type Employee struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	ReallyName   string `json:"reallyName"`
	DepartmentId string `json:"departmentId"`
	CenterId 	 string `json:"centerId"`
}

//根据id获取员工信息
func FindEmployeeById(id int) (Employee,error) {

	var employee Employee

	db := lessgo.GetMySQL()
	defer db.Close()

	sql := "select user_id,username,really_name,department_id,center_id from employee where user_id=?"

	rows, err := db.Query(sql, id)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return  employee, err
	}

	if rows.Next() {
		err := commonlib.PutRecord(rows,&employee.UserId, &employee.UserName,&employee.ReallyName, &employee.DepartmentId, &employee.CenterId)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return  employee, err
		}
	}

	return employee , nil
}

