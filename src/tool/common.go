/**
 * Title：
 * 
 * Description: 
 *
 * Author: Ivan
 *
 * Create Time: 2014-03-12 14:21
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2014-03-12 Ivan 创建文件
*/
package tool

import (
	"github.com/hjqhezgh/lessgo"
	"strings"
	"net/http"
	"errors"
)

func GetCurrentEmployeeRoles(r *http.Request) (roleCodes []string, err error) {
	employee := lessgo.GetCurrentEmployee(r)

	if employee.UserId == "" {
		lessgo.Log.Warn("用户未登陆")
		err = errors.New("用户未登陆")
		return roleCodes, err
	}

	roleCodes = strings.Split(employee.RoleCode, ",")

	return roleCodes, nil
}
