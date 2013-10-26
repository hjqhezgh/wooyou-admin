/**
 * Title：
 *
 * Description:
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-27 11:30
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-27 Ivan 创建文件
 */
package wooyousite

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"io"
	"os"
	"strings"
)

//分析URL得出当前url访问的终端类型
func getTerminal(url string) (terminal string) {
	strs := strings.Split(url, "/")
	if len(strs) > 1 {
		terminal = strs[1]
	} else {
		terminal = "web"
	}

	return terminal
}

//分析URL得出当前url访问的实体ID
func getEntityId(url string) (entityId string) {
	strs := strings.Split(url, "/")
	if len(strs) > 2 {
		entityId = strings.Split(url, "/")[2]
	} else {
		entityId = ""
	}

	return entityId
}

func UploadImage(imageName, disName, savePath string) error {
	if imageName != "" {
		tmpFile, err := os.OpenFile("../tmp/"+imageName, os.O_RDWR, 0777)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return err
		}

		_, err = os.Stat(savePath)

		if err != nil && os.IsNotExist(err) {
			lessgo.Log.Error(savePath, "文件夹不存在，创建")
			os.MkdirAll(savePath, 0777)
		}

		disFile, err := os.Create(savePath + "/" + disName)

		if err != nil {
			lessgo.Log.Error(err.Error())
			return err
		}

		io.Copy(disFile, tmpFile)

		os.Remove("../tmp/" + imageName)
	}

	return nil
}

func GetImageFileName(imageName string) string {
	return commonlib.Substr(imageName, strings.LastIndex(imageName, "/")+1, len(imageName))
}

func GetImageFileSuffix(imageName string) string {
	return commonlib.Substr(imageName, strings.LastIndex(imageName, ".")+1, len(imageName))
}
