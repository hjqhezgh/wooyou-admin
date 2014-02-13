// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-02-11 10:06
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-02-11 10:06 black 创建文档
package tool

import (
	"bufio"
	"github.com/hjqhezgh/lessgo"
	"os"
	"strings"
	"fmt"
)

func SendMsg(){

	smsTemp := `$name家长：吾幼儿童英语定于2月15日后正式开班啦！为避免插班，等待下一个班情况，请欲报名的小朋友们速速报名，抢占新班位置，报名咨询：88609091【吾幼英语社区】`

	f, err := os.Open("phone.txt")

	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}

	defer f.Close()

	br := bufio.NewReader(f)

	phones := []string{}
	errPhones := []string{}

	for {
		//每次读取一行
		line, err := br.ReadString('\n')

		if err != nil {
			if err.Error() != "EOF" {
				lessgo.Log.Error(err.Error())
				return
			} else {
				break
			}
		}

		infos := strings.Split(line,",")
		name := infos[0]
		phone := infos[1]

		phone = strings.Replace(phone,"\n","",-1)

		if checkPhone(phones,phone){
			continue
		}else{
			if len(phone) != 11 {
				errPhones = append(errPhones,phone)
				fmt.Println("错误号码：",phone,"错误原因：非法的手机号")
				continue
			}

			phones = append(phones,phone)

			msg := strings.Replace(smsTemp,"$name",name,-1)

			fmt.Println(phone,msg)

			smsResult, err := SendMessage(phone, msg)

			if err != nil {
				errPhones = append(errPhones,phone)
				fmt.Println("错误号码：",phone,"错误原因："+err.Error())
				continue
			}

			if smsResult.Result != 0 { //请求短信接口没有成功
				errPhones = append(errPhones,phone)
				fmt.Println("错误号码：",phone,"错误原因："+smsResult.Msg)
				continue
			}
		}
	}

	fmt.Println("错误号码：",errPhones)
}

func checkPhone(phones []string ,phone string) bool{
	for _,value := range phones{
		if value == phone {
			return true
		}
	}

	return false
}
