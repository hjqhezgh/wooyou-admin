// Title：新闻正文的图片上传功能
//
// Description:
//
// Author:black
//
// Createtime:2013-08-28 09:55
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-28 09:55 black 创建文档
package wooyousite

import (
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

/*****
 * 获取上传图片的随机不重复文件名
 */
func findRandomFileName(sourceFileName string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := ""

	for i := 0; i < 4; i++ {
		str += fmt.Sprint(r.Intn(10))
	}

	return fmt.Sprint(time.Now().UnixNano(), str)
}

func NewsImageUplaodAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	newsImgDir, _ := lessgo.Config.GetValue("wooyou", "newsImgDir")

	_, err := os.Stat(newsImgDir)

	if err != nil && os.IsNotExist(err) {
		lessgo.Log.Info(newsImgDir,"，创建")
		os.Mkdir(newsImgDir, 0777)
	}

	fn, header, err := r.FormFile("imgFile")

	if err != nil && os.IsNotExist(err) {
		m["error"] = 1
		m["message"] = err.Error()
		lessgo.Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	newFileName := findRandomFileName(header.Filename)

	f, err := os.Create(newsImgDir+"/" + newFileName)

	if err != nil {
		m["error"] = 1
		m["message"] = err.Error()
		lessgo.Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	defer f.Close()

	io.Copy(f, fn)

	m["error"] = 0
	m["url"] = "/newsimg/" + newFileName

	commonlib.OutputJson(w, m, " ")
}
