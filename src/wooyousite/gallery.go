/**
 * Title：
 *
 * Description:
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-28 10:05
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-28 Ivan 创建文件
 */
package wooyousite

import (
	"fmt"
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func GalleryLoad(id string) (*Gallery, string, error) {
	return GetGalleryById(id)
}

func GalleryLoadAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	id := r.FormValue("id")
	gallery, msg, err := GalleryLoad(id)
	if err != nil {
		m["success"] = false
		m["msg"] = msg
		commonlib.OutputJson(w, m, "")
		return
	}

	datas := []lessgo.LoadFormObject{}
	jsonId := lessgo.LoadFormObject{"id", id}
	datas = append(datas, jsonId)
	jsonTitle := lessgo.LoadFormObject{"title", gallery.Title}
	datas = append(datas, jsonTitle)
	jsonComment := lessgo.LoadFormObject{"comment", gallery.Comment}
	datas = append(datas, jsonComment)
	jsonImageCategory := lessgo.LoadFormObject{"image_category", gallery.ImageCategory}
	datas = append(datas, jsonImageCategory)
	jsonImageName := lessgo.LoadFormObject{"image_name", gallery.ImageName}
	datas = append(datas, jsonImageName)
	jsonImageSuffix := lessgo.LoadFormObject{"image_suffix", gallery.ImageSuffix}
	datas = append(datas, jsonImageSuffix)
	imagePath := "/artimg/"
	jsonImage1 := lessgo.LoadFormObject{"image1", imagePath + gallery.ImageName + "_" + IMAGE_SIZE_1 + "." + gallery.ImageSuffix}
	datas = append(datas, jsonImage1)
	jsonImage2 := lessgo.LoadFormObject{"image2", imagePath + gallery.ImageName + "_" + IMAGE_SIZE_2 + "." + gallery.ImageSuffix}
	datas = append(datas, jsonImage2)
	jsonImage3 := lessgo.LoadFormObject{"image3", imagePath + gallery.ImageName + "_" + IMAGE_SIZE_3 + "." + gallery.ImageSuffix}
	datas = append(datas, jsonImage3)

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

func GallerySave(gallery *Gallery) (string, error) {
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
	// 生成图片信息
	sqlInsGal := "insert into gallery (title, comment, image_name, image_suffix) values (?,?,?,?);"
	res, err := TxInsert(tx, sqlInsGal, gallery.Title, gallery.Comment, gallery.ImageName, gallery.ImageSuffix)
	lessgo.Log.Debug(sqlInsGal, gallery.Title, gallery.Comment, gallery.ImageName, gallery.ImageSuffix)
	if err != nil {
		return "插入图片信息，数据库异常", err
	}
	// 获取图片信息ID
	id, err := res.LastInsertId()
	if err != nil {
		return "获取生成图片信息ID失败!", err
	}
	// 保存图片关联类别信息
	category := gallery.ImageCategory
	sqlInsCG := "insert into category_gallery(category_id, gallery_id) values (?,?);"
	if category != "" {
		categoryArr := strings.Split(category, ",")
		for _, categoryId := range categoryArr {
			_, err = TxInsert(tx, sqlInsCG, categoryId, id)
			if err != nil {
				return "插入图片关联类别信息，数据库异常", err
			}
		}
	}
	// 返回当前插入图片信息
	gallery.Id = int(id)
	lessgo.Log.Info("保存时间: ", time.Now().Sub(t))

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

func GallerySaveAction(w http.ResponseWriter, r *http.Request) {
	GalleryFormAction(w, r, FORM_ACTION_TYPE.CREATE.Key)
}

func GalleryUpdate(gallery *Gallery) (string, error) {
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
	// 更新商品分类信息
	sqlUpdGal := "update gallery set title=?, comment=?, image_suffix=? where id=?;"
	_, err = TxUpdate(tx, sqlUpdGal, gallery.Title, gallery.Comment, gallery.ImageSuffix, gallery.Id)
	lessgo.Log.Debug(sqlUpdGal, gallery.Title, gallery.Comment, gallery.Id)
	if err != nil {
		return "更新商品分类信息，数据库异常", err
	}
	// 删除图片关联类别信息
	sqlDelCG := "delete from category_gallery where gallery_id=?;"
	_, err = TxDelete(tx, sqlDelCG, gallery.Id)
	if err != nil {
		return "删除图片关联类别信息，数据库异常", err
	}
	// 保存新图片关联类别信息
	category := gallery.ImageCategory
	sqlInsCG := "insert into category_gallery(category_id, gallery_id) values (?,?);"
	if category != "" {
		categoryArr := strings.Split(category, ",")
		for _, categoryId := range categoryArr {
			_, err = TxInsert(tx, sqlInsCG, categoryId, gallery.Id)
			if err != nil {
				return "插入图片关联类别信息，数据库异常", err
			}
		}
	}
	lessgo.Log.Info("更新商品分类信息时间: ", time.Now().Sub(t))
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

func GalleryUpdateAction(w http.ResponseWriter, r *http.Request) {
	GalleryFormAction(w, r, FORM_ACTION_TYPE.UPDATE.Key)
}

func GalleryDelete(id string) (string, error) {

	db := lessgo.GetMySQL()
	defer db.Close()

	sqlDelGal := "delete from gallery where id=?;"
	lessgo.Log.Debug(sqlDelGal)

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return "删除商品分类信息，数据库异常", err
	}

	stmt, err := tx.Prepare(sqlDelGal)
	if err != nil {
		tx.Rollback()
		return "删除商品分类信息，数据库异常", err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		tx.Rollback()
		return "删除商品分类信息，数据库异常", err
	}

	sqlDelGal1 := "delete from category_gallery where gallery_id=?;"
	lessgo.Log.Debug(sqlDelGal1)

	stmt1, err := tx.Prepare(sqlDelGal1)
	if err != nil {
		tx.Rollback()
		return "删除商品分类信息，数据库异常", err
	}

	_, err = stmt1.Exec(id)

	if err != nil {
		tx.Rollback()
		return "删除商品分类信息，数据库异常", err
	}

	tx.Commit()

	t := time.Now()
	lessgo.Log.Info("删除商品分类信息时间: ", time.Now().Sub(t))

	return "删除成功", nil
}

func GalleryDeleteAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	err := r.ParseForm()
	if err != nil {
		lessgo.Log.Error(err.Error())
		return
	}
	id := r.FormValue("id")

	msg, err := GalleryDelete(id)
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["sucess"] = false
	} else {
		m["success"] = true
	}
	m["msg"] = msg
	commonlib.RenderTemplate(w, r, "notify_delete.html", m, nil, "../template/notify_delete.html")

	return
}

func GalleryFormAction(w http.ResponseWriter, r *http.Request, actionType string) {
	m := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		lessgo.Log.Error(err.Error())
		m["success"] = false
		m["msg"] = "解析表单信息出错"
		commonlib.OutputJson(w, m, "")
		return
		return
	}

	artImgDir, _ := lessgo.Config.GetValue("wooyou", "artImgDir")

	gallery := new(Gallery)
	gallery.Id, _ = strconv.Atoi(r.FormValue("id"))
	gallery.Title = r.FormValue("title")
	gallery.Comment = r.FormValue("comment")
	gallery.ImageCategory = r.FormValue("image_category")
	image1 := GetImageFileName(r.FormValue("image1"))
	image2 := GetImageFileName(r.FormValue("image2"))
	image3 := GetImageFileName(r.FormValue("image3"))
	imageName := r.FormValue("image_name")
	if imageName == "" {
		gallery.ImageName = GetTimeString()
	} else {
		gallery.ImageName = imageName
	}
	imageSuffix := r.FormValue("image_suffix")

	msg := ""
	switch actionType {
	case FORM_ACTION_TYPE.CREATE.Key:
		if image1 == "" || image2 == "" || image3 == "" {
			m["success"] = false
			m["msg"] = "必须上传三张图片"
			lessgo.Log.Error("必须上传三张图片")
			commonlib.OutputJson(w, m, " ")
			return
		} else {
			if !(GetImageFileSuffix(image1) == GetImageFileSuffix(image2) && GetImageFileSuffix(image1) == GetImageFileSuffix(image3)) {
				m["success"] = false
				m["msg"] = "请确保三张图片的文件格式一致"
				lessgo.Log.Error("请确保三张图片的文件格式一致")
				commonlib.OutputJson(w, m, " ")
				return
			}
			gallery.ImageSuffix = GetImageFileSuffix(image1)
		}
		msg, err = GallerySave(gallery)
	case FORM_ACTION_TYPE.UPDATE.Key:
		if image1 == "" || image2 == "" || image3 == "" {
			if (image1 != "" && GetImageFileSuffix(image1) != imageSuffix) || (image2 != "" && GetImageFileSuffix(image2) != imageSuffix) || (image3 != "" && GetImageFileSuffix(image3) != imageSuffix) {
				m["success"] = false
				m["msg"] = "请确保三张图片的文件格式一致"
				lessgo.Log.Error("请确保三张图片的文件格式一致")
				commonlib.OutputJson(w, m, " ")
				return
			}
			gallery.ImageSuffix = imageSuffix
		} else {
			if !(GetImageFileSuffix(image1) == GetImageFileSuffix(image2) && GetImageFileSuffix(image1) == GetImageFileSuffix(image3)) {
				m["success"] = false
				m["msg"] = "请确保三张图片的文件格式一致"
				lessgo.Log.Error("请确保三张图片的文件格式一致")
				commonlib.OutputJson(w, m, " ")
				return
			}
			gallery.ImageSuffix = GetImageFileSuffix(image1)
		}
		msg, err = GalleryUpdate(gallery)
	}
	// 图片处理
	// 图片1
	if UploadImage(image1, gallery.ImageName+"_"+IMAGE_SIZE_1+"."+gallery.ImageSuffix, artImgDir) != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		lessgo.Log.Error(err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}
	// 图片2
	if UploadImage(image2, gallery.ImageName+"_"+IMAGE_SIZE_2+"."+gallery.ImageSuffix, artImgDir) != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		lessgo.Log.Error(err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}
	// 图片3
	if UploadImage(image3, gallery.ImageName+"_"+IMAGE_SIZE_3+"."+gallery.ImageSuffix, artImgDir) != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		lessgo.Log.Error(err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	if err != nil {
		lessgo.Log.Error(err.Error())
		m["sucess"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		m["msg"] = msg
	} else {
		m["success"] = true
		m["datas"] = gallery
	}
	commonlib.OutputJson(w, m, "")

	return
}

func GalleryImageCategoryAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	options, err := GetSubcategories()
	if err != nil {
		m["success"] = false
		m["code"] = "100"
		m["msg"] = "获取类别信息出错" + err.Error()
		commonlib.OutputJson(w, m, "")
		return
	}

	m["success"] = true
	m["code"] = "200"
	m["datas"] = options
	commonlib.OutputJson(w, m, "")

	return
}

func GetSubcategories() (options []*SelectOption, err error) {
	sql := "select id, name from category where parent_code is not null and parent_code <> '';"
	rows, err := DBSelect(sql)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		option := new(SelectOption)

		err = commonlib.PutRecord(rows, &option.Value, &option.Description)
		if err != nil {
			lessgo.Log.Error("封装对象数据时发生错误")
			return nil, err
		}
		options = append(options, option)
	}

	return options, nil
}

func GetGalleryById(id string) (*Gallery, string, error) {
	sql := "select gal.id, gal.title, gal.comment, (select GROUP_CONCAT(category_id) from category_gallery where gallery_id=gal.id) as image_category, gal.image_name, gal.image_suffix from gallery gal where gal.id=?;"
	rows, err := DBSelect(sql, id)
	lessgo.Log.Debug(sql, id)
	if err != nil {
		lessgo.Log.Error(err.Error())
		return nil, "查询信息出错", err
	}
	gallery := new(Gallery)
	if rows.Next() {
		err = commonlib.PutRecord(rows, &gallery.Id, &gallery.Title, &gallery.Comment, &gallery.ImageCategory, &gallery.ImageName, &gallery.ImageSuffix)
		if err != nil {
			lessgo.Log.Error("封装对象数据时发生错误")
			return nil, "封装信息出错", err
		}
	}

	return gallery, "处理成功", nil
}

//商品分类分页数据服务
func GalleryListAction(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	pageNoString := r.FormValue("page")
	pageNo := 1
	if pageNoString != "" {
		pageNo, err = strconv.Atoi(pageNoString)
		if err != nil {
			pageNo = 1
			lessgo.Log.Warn("错误的pageNo:", pageNo)
		}
	}

	pageSizeString := r.FormValue("rows")
	pageSize := 10
	if pageSizeString != "" {
		pageSize, err = strconv.Atoi(pageSizeString)
		if err != nil {
			lessgo.Log.Warn("错误的pageSize:", pageSize)
		}
	}

	title := r.FormValue("title-like")
	comment := r.FormValue("comment-like")

	params := []interface{}{}

	sql := "select gal.id, gal.title, gal.comment, gal.image_name, gal.image_suffix, "
	sql += "    (select GROUP_CONCAT(c.name) "
	sql += "        from category c left join category_gallery cg ON c.id = cg.category_id "
	sql += "        where cg.gallery_id = gal.id) as image_category "
	sql += "from gallery gal where 1=1"
	if strings.TrimSpace(title) != "" {
		params = append(params, "%"+title+"%")
		sql += " and gal.title like ? "
	}
	if strings.TrimSpace(comment) != "" {
		params = append(params, "%"+comment+"%")
		sql += " and gal.comment like ? "
	}

	countSql := "select count(1) from (" + sql + ") num"

	lessgo.Log.Debug(countSql)

	db := lessgo.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	totalNum := 0

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门，错误信息：" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	currPageNo := pageNo

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	sql += " order by gal.id desc limit ?,?"

	lessgo.Log.Debug(sql)

	params = append(params, (currPageNo-1)*pageSize)
	params = append(params, pageSize)

	rows, err = db.Query(sql, params...)

	if err != nil {
		lessgo.Log.Warn(err.Error())
		m["success"] = false
		m["code"] = 100
		m["msg"] = "系统发生错误，请联系IT部门"
		commonlib.OutputJson(w, m, " ")
		return
	}

	objects := []interface{}{}

	for rows.Next() {

		model := new(lessgo.Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for i := 0; i < 5; i++ {
			prop := new(lessgo.Prop)
			prop.Name = fmt.Sprint(i)
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			lessgo.Log.Warn(err.Error())
			m["success"] = false
			m["code"] = 100
			m["msg"] = "系统发生错误，请联系IT部门"
			commonlib.OutputJson(w, m, " ")
			return
		}

		objects = append(objects, model)
	}

	pageData := commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects)

	m["PageData"] = pageData
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*lessgo.Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": lessgo.GetPropValue, "compareInt": lessgo.CompareInt, "dealJsonString": lessgo.DealJsonString}, "../lessgo/template/entity_page.json")

}
