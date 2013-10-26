/**
 * Title：
 *
 * Description:
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-26 11:52
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-26 Ivan 创建文件
 */
package wooyousite

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewsLoad(Id string) (*News, string, error) {
	return GetNewsById(Id)
}

func NewsLoadAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	id := r.FormValue("id")
	news, msg, err := NewsLoad(id)
	if err != nil {
		m["success"] = false
		m["msg"] = msg
		commonlib.OutputJson(w, m, "")
		return
	}

	datas := []lessgo.LoadFormObject{}
	jsonId := lessgo.LoadFormObject{"id", id}
	datas = append(datas, jsonId)
	jsonTitle := lessgo.LoadFormObject{"title", news.Title}
	datas = append(datas, jsonTitle)
	jsonNewsContent := lessgo.LoadFormObject{"news_content", news.NewsContent}
	datas = append(datas, jsonNewsContent)
	previewImage := "/pic/news/" + news.PreviewImage
	jsonPreviewImage := lessgo.LoadFormObject{"preview_image", previewImage}
	datas = append(datas, jsonPreviewImage)
	jsonPreviewText := lessgo.LoadFormObject{"preview_text", news.PreviewText}
	datas = append(datas, jsonPreviewText)
	jsonIsCarousel := lessgo.LoadFormObject{"is_carousel", news.IsCarousel}
	datas = append(datas, jsonIsCarousel)
	carouselImage := "/pic/news/" + news.CarouselImage
	jsonCarouselImage := lessgo.LoadFormObject{"carousel_image", carouselImage}
	datas = append(datas, jsonCarouselImage)
	if err != nil {
		m["sucess"] = false
		m["msg"] = msg
	} else {
		m["success"] = true
		m["datas"] = datas
	}
	commonlib.OutputJson(w, m, "")

	return
}

func NewsSave(news *News) (string, error) {
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
	// 生成新闻信息
	sqlNews := "insert into news (title, news_content, preview_image, preview_text, is_carousel, carousel_image, upload_time) values (?,?,?,?,?,?,?);"
	res, err := TxInsert(tx, sqlNews, news.Title, news.NewsContent, news.PreviewImage, news.PreviewText, news.IsCarousel, news.CarouselImage, GetTimeString())
	lessgo.Log.Debug(sqlNews, news.Title, news.NewsContent, news.PreviewImage, news.PreviewText, news.IsCarousel, news.CarouselImage)
	if err != nil {
		return "插入用户基本数据，数据库异常", err
	}
	// 获取生成新闻信息ID
	id, err := res.LastInsertId()
	if err != nil {
		return "获取生成新闻信息ID失败!", err
	}
	news.Id = int(id)
	lessgo.Log.Info("保存用户信息时间: ", time.Now().Sub(t))
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

	return "处理成功", nil
}

func NewsSaveAction(w http.ResponseWriter, r *http.Request) {
	NewsFormAction(w, r, FORM_ACTION_TYPE.CREATE.Key)
}

func NewsUpdate(news *News) (string, error) {
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
	// 更新新闻信息
	sqlNews := "update news set news_content=?, preview_text=?, is_carousel=? where id=?;"
	_, err = TxUpdate(tx, sqlNews, news.NewsContent, news.PreviewText, news.IsCarousel, news.Id)
	if err != nil {
		return "更新新闻信息，数据库异常", err
	}
	if news.PreviewImage != "" {
		sqlUpdPre := "update news set preview_image=? where id=?;"
		_, err = TxUpdate(tx, sqlUpdPre, news.PreviewImage, news.Id)
		if err != nil {
			return "更新新闻信息，数据库异常", err
		}
	}
	if news.CarouselImage != "" {
		sqlUpdCar := "update news set carousel_image=? where id=?;"
		_, err = TxUpdate(tx, sqlUpdCar, news.CarouselImage, news.Id)
		if err != nil {
			return "更新新闻信息，数据库异常", err
		}
	}
	lessgo.Log.Info("更新时间: ", time.Now().Sub(t))
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

	return "处理成功", nil
}

func NewsUpdateAction(w http.ResponseWriter, r *http.Request) {
	NewsFormAction(w, r, FORM_ACTION_TYPE.UPDATE.Key)
}

func NewsDelete(id string) (string, error) {
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
	// 删除新闻信息
	sqlDelNews := "delete from news where id=?;"
	_, err = TxDelete(tx, sqlDelNews, id)
	if err != nil {
		return "删除新闻信息，数据库异常", err
	}
	lessgo.Log.Info("删除用户信息时间: ", time.Now().Sub(t))
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

	return "删除成功", nil
}

func NewsDeleteAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	err := r.ParseForm()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}
	id := r.FormValue("id")
	msg, err := NewsDelete(id)
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["sucess"] = false
	} else {
		m["success"] = true
	}
	m["msg"] = msg
	commonlib.RenderTemplate(w, r, "notify_delete.html", m, nil, "../template/component/"+getTerminal(r.URL.Path)+"/notify_delete.html")

	return
}

func NewsFormAction(w http.ResponseWriter, r *http.Request, actionType string) {
	m := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}

	news := new(News)
	news.Id, _ = strconv.Atoi(r.FormValue("id"))
	news.Title = r.FormValue("title")
	news.NewsContent = r.FormValue("news_content")
	previewImage := r.FormValue("preview_image")
	news.PreviewImage = commonlib.Substr(previewImage, strings.LastIndex(previewImage, "/")+1, len(previewImage))
	news.PreviewText = r.FormValue("preview_text")
	news.IsCarousel = r.FormValue("is_carousel")
	carouselImage := r.FormValue("carousel_image")
	news.CarouselImage = commonlib.Substr(carouselImage, strings.LastIndex(carouselImage, "/")+1, len(carouselImage))

	msg := ""
	switch actionType {
	case FORM_ACTION_TYPE.CREATE.Key:
		msg, err = NewsSave(news)
	case FORM_ACTION_TYPE.UPDATE.Key:
		msg, err = NewsUpdate(news)
	}

	// 图片文件处理
	if news.PreviewImage != "" {
		tmpFile, err := os.OpenFile("../tmp/"+news.PreviewImage, os.O_RDWR, 0777)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			lessgo.Log.Error(err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = os.Stat("../static/pic/news")

		if err != nil && os.IsNotExist(err) {
			lessgo.Log.Error("/static/pic/news", "文件夹不存在，创建")
			os.MkdirAll("../static/pic/news", 0777)
		}

		disFile, err := os.Create("../static/pic/news/" + news.PreviewImage)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			lessgo.Log.Error(err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		io.Copy(disFile, tmpFile)

		os.Remove("../tmp/" + news.PreviewImage)
	}
	// 图片文件处理
	if news.CarouselImage != "" {
		tmpFile, err := os.OpenFile("../tmp/"+news.CarouselImage, os.O_RDWR, 0777)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			lessgo.Log.Error(err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		_, err = os.Stat("../static/pic/news")

		if err != nil && os.IsNotExist(err) {
			lessgo.Log.Error("/static/pic/news", "文件夹不存在，创建")
			os.MkdirAll("../static/pic/news", 0777)
		}

		disFile, err := os.Create("../static/pic/news/" + news.CarouselImage)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			lessgo.Log.Error(err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		io.Copy(disFile, tmpFile)

		os.Remove("../tmp/" + news.CarouselImage)
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["sucess"] = false
		m["msg"] = msg
	} else {
		m["success"] = true
		m["datas"] = news
	}
	commonlib.OutputJson(w, m, "")

	return
}

func GetNewsById(id string) (*News, string, error) {
	sql := "select title, news_content, preview_image, preview_text, is_carousel, carousel_image, upload_time from news where id=?;"
	lessgo.Log.Debug(sql, id)
	rows, err := DBSelect(sql, id)

	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, "查询用户信息出错", err
	}

	news := new(News)
	news.Id, _ = strconv.Atoi(id)
	if rows.Next() {
		err = commonlib.PutRecord(rows, &news.Title, &news.NewsContent, &news.PreviewImage,
			&news.PreviewText, &news.IsCarousel, &news.CarouselImage, &news.UploadTime)
		if err != nil {
			lessgo.Log.Error(err.Error())
			return nil, "封装用户信息出错", err
		}
		news.UploadTime = TimeFormat(news.UploadTime)
	}
	lessgo.Log.Debug(news)

	return news, "处理成功", nil
}
