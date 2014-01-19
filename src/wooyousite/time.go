/**
 * User: Samurai
 * Date: 13-8-2
 * Time: 下午7:02
 *
 * 包说明：公共函数包，与业务无关的通用方法都写在此。
 */
package wooyousite

import (
	"fmt"
	"time"
)

/*
 * 函数说明：获取当前日期的字符串
 * 传入参数：
 * 返回值：  yyyyMMdd
 */
func GetDateString() string {
	t := time.Now()
	str := fmt.Sprintf("%04d%02d%02d", t.Year(), int(t.Month()), t.Day())
	return str
}

/*
 * 函数说明：获取当前时间的字符串
 * 传入参数：
 * 返回值：  yyyyMMddhhmmss
 */
func GetTimeString() string {
	t := time.Now()
	str := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	return str
}

/*
 * 函数说明：获取当前时间的字符串
 * 传入参数：
 * 返回值：  yyyyMMddhhmmss
 */
func GetNoDateString() string {
	t := time.Now()
	str := fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
	return str
}

/*
 * 函数说明：字符串截取
 * 传入参数：str		:  待截取的字符串
 			start  	:  截取的起始位置
 			length  :  截取长度
 * 返回值：  截取的字符串
*/
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func TimeFormat(tm string) string {
	b := []byte(tm)
	var bb []byte
	for k, v := range b {
		bb = append(bb, v)
		if k == 3 || k == 5 {
			bb = append(bb, '-')
		}
		if k == 7 {
			bb = append(bb, ' ')
		}
		if k == 9 || k == 11 {
			bb = append(bb, ':')
		}
	}
	return string(bb)
}

func GetTimeNow() string {
	t := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	return str
}
