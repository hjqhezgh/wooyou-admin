// Title：应用入口
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 13:04
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 13:04 black 创建文档
package main

import (
	"fmt"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"server"
	"wooyousite"
	"strconv"
)


func main() {

	r := lessgo.ConfigLessgo()

	portString, _ := lessgo.Config.GetValue("lessgo", "port")

	port, _ := strconv.Atoi(portString)

	for url, handler := range handlers {
		r.HandleFunc(url, handler)
	}

//	go server.ReportSend()

	http.Handle("/", r)

	http.Handle("/js/", http.FileServer(http.Dir("../static")))
	http.Handle("/img/", http.FileServer(http.Dir("../static")))
	http.Handle("/css/", http.FileServer(http.Dir("../static")))
	http.Handle("/json/", http.FileServer(http.Dir("../static")))
	http.Handle("/newsimg/", http.FileServer(http.Dir("../")))
	http.Handle("/artimg/", http.FileServer(http.Dir("../")))

	fmt.Println("服务器监听", portString, "端口")
//	go server.UpdateVideoStatus()

	lessgo.Log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

//URL映射列表
var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/login": server.LoginAction,

	//根据角色ID获取员工列表
	"/employeeListByRoleId.json": server.EmployeeListByRoleIdAction,
	"/employeeByCenterId.json": server.EmployeeListByCenterIdAction,

	//音频相关服务
	"/consultant_phone_list.json":        server.ConsultantPhoneListAction,
	"/consultant_phone_detail_list.json": server.ConsultantPhoneDetailListAction,
	"/queryVedio.json":                   server.VideoListAction,
	"/downloadAudio":                     server.DownloadAudioAction,
	"/audioNoteLoad.json":                server.AudioNoteLoadAction,
	"/audioNoteSave.json":                server.AudioNoteSaveAction,

	//客户相关服务
	"/consumer.json":     server.ConsumerListAction,
	"/consumerSave.json": server.ConsumerSaveAction,
	"/consumerLoad.json": server.ConsumerLoadAction,

	//Call Center统计
	"/callCenterStatistics.json": server.CallCenterStatisticsAction,
	//CD获取可以分配给CallCenter名单
	"/validForCallCenterList.json": server.ValidForCallCenterListAction,
	//中心的CallCenter详情
	"/centerCallCenterDetail.json": server.CenterCallCenterDetailAction,
	//分配给CallCenter
	"/web/consumer/sendToCallCenter": server.SendToCallCenterAction,
	//全部分配给CallCenter
	"/web/consumer/allSendToCallCenter": server.AllSendToCallCenter,
	//tmk运营报表
	"/tmkStatistics.json": server.TmkStatisticsAction,
	//tmk运营报表详情
	"/tmk_statistics_detail.json": server.TmkStatisticsDetailAction,
	//分配tmk页面的表单加载
	"/sendToTmkLoad.json": server.SendToTmkLoadAction,
	//分配tmk页面的表单保存
	"/sendToTmkSave.json": server.SendToTmkSaveAction,
	//客户状态变更
	"/consumerStatusChange": server.ConsumerStatusChangeAction,

	//网站内容管理
	"/web/gallery/delete.json":         wooyousite.GalleryDeleteAction,
	"/web/gallery/image_category.json": wooyousite.GalleryImageCategoryAction,
	"/web/gallery/load.json":           wooyousite.GalleryLoadAction,
	"/web/gallery/save.json":           wooyousite.GallerySaveAction,
	"/web/gallery/update.json":         wooyousite.GalleryUpdateAction,
	"/web/gallery/list.json":			wooyousite.GalleryListAction,
	"/web/news/delete.json":            wooyousite.NewsDeleteAction,
	"/web/news/load.json":              wooyousite.NewsLoadAction,
	"/web/news/save.json":              wooyousite.NewsSaveAction,
	"/web/news/update.json":            wooyousite.NewsUpdateAction,
	"/newsImageUplaod":            		wooyousite.NewsImageUplaodAction,

	//app开户数据读取
	"/addAppAccountLoad.json": server.AddAppAccountLoadAction,
	//app开户数据保存
	"/addAppAccountSave.json": server.AddAppAccountSaveAction,

	//课程信息相关服务
	"/course.json": server.CourseListAction,
	"/web/course/save.json": server.CourseSaveAction,
	"/web/course/load.json": server.CourseLoadAction,
	"/courseByCenterId.json": server.CourseByCenterIdListAction,
	"/time_section.json": server.TimeSectionListAction,
	"/class_schedule_detail.json": server.ClassScheduleDetailListAction,
	"/class_schedule_detail/save.json": server.ClassScheduleDetailSaveAction,
	"/class_schedule_detail/load.json": server.ClassScheduleDetailLoadAction,
	"/lessonByClassId.json": server.LessonByClassIdAction,
	"/timeSectionByCenterId.json": server.TimeSectionByCenterIdAction,


	//班级相关服务
	"/wyclass.json": server.WyClassListAction,
	"/web/wyclass/update.json": server.WyClassUpdateAction,
	"/web/wyclass/load.json": server.WyClassLoadAction,
	"/wyclass/loadChildInClass.json": server.LoadChildInClassAction,
	"/wyclass/saveChildToClass.json": server.SaveChildToClassAction,
	"/classByCenterId.json": server.ClassByCenterIdAction,

	"/roomByCenterId.json": server.RoomByCenterIdListAction,
}
