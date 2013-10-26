/**
 * Title：
 *
 * Description:
 *
 * Author: Ivan
 *
 * Create Time: 2013-09-27 09:21
 *
 * Version: 1.0
 *
 * 修改历史: 版本号 修改日期 修改人 修改说明
 *   1.0 2013-09-27 Ivan 创建文件
 */
package wooyousite

type Category struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Ename       string      `json:"ename"`
	Code        string      `json:"code"`
	ParentCode  string      `json:"parentCode"`
	Rank        int         `json:"rank"`
	SubCategory []*Category `json:"subCategory"`
}

type Employee struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	ReallyName   string `json:"reallyName"`
	DepartmentId string `json:"departmentId"`
	Role         string `json:"role"`
}

type Gallery struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	ImageName     string `json:"imageName"`
	ImageSuffix   string `json:"imageSuffix"`
	Comment       string `json:"comment"`
	ImageCategory string `json:"imageCategory"`
}

type News struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	NewsContent   string `json:"newsContent"`
	UploadTime    string `json:"uploadTime"`
	PreviewImage  string `json:"previewImage"`
	PreviewText   string `json:"previewText"`
	IsCarousel    string `json:"isCarousel"`
	CarouselImage string `json:"carouselImage"`
}

type SelectOption struct {
	Description string `json:"desc"`
	Value       string `json:"value"`
}
