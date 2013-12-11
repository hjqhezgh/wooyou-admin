// Title：登陆相关
//
// Description:
//
// Author:Samurai
//
// Createtime:2013-11-26 10:02
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明

package server

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//md5加密方法
func Md5Str(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	strMd5 := fmt.Sprintf("%x", h.Sum(nil))
	return strMd5
}

func getTimeString() string {
	t := time.Now()
	str := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	return str[2:]
}

func GenerateXml(mobile, msg string) string {
	cpid := "020000001832"
	port := "0723"
	secret_key := "87bd337f7640aab36ab613f806c3b062"
	now_time := getTimeString()
	cpmid := "12345"
	signature := Md5Str(secret_key + now_time)

	xml := `<?xml version="1.0" encoding="UTF-8"?>
	<MtPacket>
		<cpid>%s</cpid>
		<mid>0</mid>
		<cpmid>%s</cpmid>
		<mobile>%s</mobile>
		<port>%s</port>
		<msg>%s</msg>
		<signature>%s</signature>
		<timestamp>%s</timestamp>
		<validtime></validtime>
	</MtPacket>
`
	return fmt.Sprintf(xml, cpid, cpmid, mobile, port, msg, signature, now_time)
}

func SendMessage(mobile, msg string) (SmsResult, error) {
	trim_mgs := strings.TrimLeft(strings.TrimRight(msg, " "), " ")
	data := GenerateXml(mobile, trim_mgs)
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", "http://221.179.216.74/providermt", strings.NewReader(data))

	reqest.Header.Set("Content-Type", "text/xml; charset=GBK")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")

	response, _ := client.Do(reqest)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		fmt.Println("****************", bodystr)
		return GetResult(bodystr)
	}

	return SmsResult{}, nil
}

type SmsResult struct {
	MtResponse string `xml:"MtResponse"`
	Mid        string `xml:"mid"`
	Cpmid      string `xml:"cpmid"`
	Result     int    `xml:"result"`
	Msg        string
}

func GetResult(data string) (SmsResult, error) {

	var ret SmsResult

	err := xml.Unmarshal([]byte(data), &ret)

	if err != nil {
		fmt.Printf("error: %v", err)
		return ret, err
	}

	if ret.Result == 0 {
		ret.Msg = "发送成功"
	} else if ret.Result == 1001 {
		ret.Msg = "重要字段有空值"
	} else if ret.Result == 1002 {
		ret.Msg = "帐号密码验证错误"
	} else if ret.Result == 1003 {
		ret.Msg = "IP限制访问"
	} else if ret.Result == 1004 {
		ret.Msg = "内容超长，超出180个字"
	} else if ret.Result == 1005 {
		ret.Msg = "群发号码过多，超出100个"
	} else if ret.Result == 1006 {
		ret.Msg = "子端口不合法"
	} else if ret.Result == 1007 {
		ret.Msg = "信息内容非法"
	} else if ret.Result == 1008 {
		ret.Msg = "非法CPID"
	} else if ret.Result == 9999 {
		ret.Msg = "未知错误"
	}

	return ret, nil
}
