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
//	"process"
	"strconv"
	_ "tool"
	"web"
	"wooyousite"
	"finance"
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

	//	tool.SendMsg()

	lessgo.Log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

//URL映射列表
var handlers = map[string]func(http.ResponseWriter, *http.Request){

	"/login": server.LoginAction,

	//根据角色ID获取员工列表
	"/employeeListByRoleId.json" : server.EmployeeListByRoleIdAction,
	"/employeeByCenterId.json"   : server.EmployeeListByCenterIdAction,
	"/courseInCenter.json"       : server.CourseInCenterAction,
	"/employeeInCenter.json"     : server.EmployeeListInCenterAction,
	"/getCurrentEmployee.json"   : server.GetCurrentEmployeeAction,

	//音频相关服务
	"/consultant_phone_list.json":        server.ConsultantPhoneListAction,
	"/consultant_phone_detail_list.json": server.ConsultantPhoneDetailListAction,
	"/queryVedio.json":                   server.VideoListAction,
	"/downloadAudio":                     server.DownloadAudioAction,
	"/audioNoteLoad.json":                server.AudioNoteLoadAction,
	"/audioNoteSave.json":                server.AudioNoteSaveAction,

	//客户相关服务
	"/consumer.json":                       web.ConsumerListAction,
	"/consumerSave.json":                   web.ConsumerSaveAction,
	"/consumerLoad.json":                   server.ConsumerLoadAction,
	"/consumer/contact_record.json":        server.ConsumerContactRecordListAction,
	"/web/consumer/backToAllConsumer.json": server.BackToAllConsumerAction,
	"/contacts/page":                       server.ContactsListAction,
	"/contacts/save.json":                  server.ContactsSaveAction,
	"/contacts/delete.json":                server.ContactsDeleteAction,
	"/web/consumerContactsLog/page.json":   server.ConsumerContactLogAction,
	"/web/consumerContactsLog/save.json":   server.ConsumerContactLogSaveAction,
	"/web/consumerContactsLog/load.json":   server.ConsumerContactLogLoadAction,
	"/web/contacts/contactsLoad.json":      server.ContactsLoadAction,
	"/web/child/data.json":                 web.ChildListAction,
	"/web/contract/page.json":              web.ContractListAction,
	"/web/contract/save.json":              web.ContractSaveAction,
	"/web/contract/load.json":              web.ContractLoadAction,
	"/web/contract/contractOfChild.json":   server.ContractOfChildAction,
	"/web/child/signInData.json":           server.ChildSignInLogListAction,
	"/web/child/addCardToSignIn.json":      server.AddCardToSignInAction,
	"/web/child/addContractToSignIn.json":  server.AddContractToSignInAction,

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
	//tmk的名单总库
	"/consumer/tmkAllConsumer.json": server.TmkAllConsumerListAction,
	//开始邀约按钮
	"/consumer/tmkInvite.json": server.TmkInviteAction,
	//tmk自己的名单库
	"/consumer/tmk_consumer.json": server.TmkConsumerSelfListAction,
	//缴费
	"/web/consumer/consumerPay.json": web.ConsumerPayAction,
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
	"/web/gallery/list.json":           wooyousite.GalleryListAction,
	"/web/news/delete.json":            wooyousite.NewsDeleteAction,
	"/web/news/load.json":              wooyousite.NewsLoadAction,
	"/web/news/save.json":              wooyousite.NewsSaveAction,
	"/web/news/update.json":            wooyousite.NewsUpdateAction,
	"/newsImageUplaod":                 wooyousite.NewsImageUplaodAction,

	//app开户数据读取
	"/addAppAccountLoad.json": server.AddAppAccountLoadAction,
	//app开户数据保存
	"/addAppAccountSave.json": server.AddAppAccountSaveAction,

	//课程信息相关服务
	"/web/course/page.json":       web.CourseListAction,
	"/web/course/save.json":       web.CourseSaveAction,
	"/web/course/load.json":       server.CourseLoadAction,
	"/courseByCenterId.json":      server.CourseByCenterIdListAction,
	"/web/time_section/page.json": web.TimeSectionListAction,
	"/web/time_section/save.json": web.TimeSectionSaveAction,
	"/class_schedule_detail.json": server.ClassScheduleDetailListAction,
	"/lessonByClassId.json":       server.LessonByClassIdAction,
	"/timeSectionByCenterId.json": server.TimeSectionByCenterIdAction,
	"/web/room/page.json": web.RoomListAction,
	"/web/room/save.json": web.RoomSaveAction,

	//课表相关
	"/web/class_schedule_detail/data.json":                       server.ClassScheduleDetailListAction,
	"/web/class_schedule_detail/quick_data.json":                 server.ClassScheduleDetailListQuickAction,
	"/web/class_schedule_detail/temp_data.json":                  server.ClassScheduleDetailListTempAction,
	"/web/class_schedule_detail/load.json":                       server.ClassScheduleDetailLoadAction,
	"/web/class_schedule_detail/add.json":                        server.ClassScheduleDetailAddAction,
	"/web/class_schedule_detail/modify.json":                     server.ClassScheduleDetailModifyAction,
	"/web/class_schedule_detail/createWeekSchedule.json":         web.CreateWeekScheduleAction,
	"/web/class_schedule_detail/leave.json":                      web.ClassScheduleDetailLeaveAction,
	"/web/class_schedule_detail/truant.json":                     web.ClassScheduleDetailTruantAction,
	"/web/class_schedule_detail/addChild.json":                   web.AddChildToClassAction,
	"/web/class_schedule_detail/signIn.json":                     web.ClassScheduleDetailSignInAction,
	"/web/class_schedule_detail/addChildForNormalTempelate.json": web.AddChildForNormalTempelateAction,
	"/web/class_schedule_detail/addChildForNormalOnce.json":      web.AddChildForNormalOnceAction,
	"/web/class_schedule_detail/removeChild.json":                web.RemoveChildFromScheduleAction,
	"/web/class_schedule_detail/removeChildForNormal.json":       web.RemoveChildFromScheduleForNormalAction,
	"/web/class_schedule_detail/pay.json":                        web.ChildPayAction,
	"/web/class_schedule_detail/changeClass.json":                web.ChangeClassScheduleAction,
	"/web/class_schedule_detail/page.json":                       web.ClassScheduleDetailPageAction,
	"/web/schedule_detail/deleteSingle.json":                     server.DeleteSingleScheduleAction,
	"/web/child/childInCenter.json":                              web.ChildInCenterAction,
	"/web/child/childInClass.json":                               web.ChildInClassListAction,
	"/web/child/childInNormalSchedule.json":                      web.ChildInNormalScheduleAction,
	"/web/child/childInParent.json":                              web.ChildInParentAction,
	"/web/child/potential.json":                      			  web.PotentialChildListAction,
	"/web/member/data.json":                      			      web.MemberListAction,
	"/web/wyclass/signInWithoutClass.json":                       web.ChildSignInWithoutClassAction,
	"/web/wyclass/sendSMS/save.json":                             server.WyClassSendSMSSaveAction,
	"/web/wyclass/addChildQuick.json":                            web.AddChildToClassQuickAction,
	"/web/wyclass/sendSMS/load.json":                             server.WyClassSendSMSLoadAction,
	"/web/wyclass/contractCheckInSave.json":                      web.ContractCheckInSaveAction,
	"/web/class_schedule_attach/data.json":                       server.ClassScheduleAttachListAction,
	"/web/class_schedule_attach/load.json":                       server.ClassScheduleAttachLoadAction,
	"/web/class_schedule_attach/save.json":                       server.ClassScheduleAttachSaveAction,
	"/web/class_schedule_attach/videoplay.html":                  server.ClassScheduleAttachVideoPlayAction,

	//班级相关服务
	"/wyclass.json":                  server.WyClassListAction,
	"/web/wyclass/save.json":         server.WyClassSaveAction,
	"/web/wyclass/load.json":         server.WyClassLoadAction,
	"/wyclass/loadChildInClass.json": server.LoadChildInClassAction,
	"/wyclass/saveChildToClass.json": server.SaveChildToClassAction,
	"/classByCenterId.json":          server.ClassByCenterIdAction,

	"/roomByCenterId.json": server.RoomByCenterIdListAction,

	//员工签到表
	"/employee_sign_in.json": server.EmployeeSignInListAction,
	//试听课程，添加至客户库中
	"/web/apply_log/addToConsumer.json": server.ApplyLogAddToConsumerAction,

	//课件管理
	"/web/courseware/page.json":      web.CoursewareListAction,
	"/web/courseware/save.json":      web.CoursewareSaveAction,
	"/web/courseware/load.json":      web.CoursewareLoadAction,
	"/web/courseware/uploadCallBack": web.CoursewareUploadCallBack,

	"/getRoleCodes.json"		  : finance.RoleCodesAction,
	"/getHandleApplyInfo.json"	  : finance.HandleApplyAction,
	"/getReceiptDetails.json"	  : finance.ReceiptDetailsAction,
	"/classifiedPendingReceiptList"	  : finance.ClassifiedPendingReceiptListAction,
	"/classifiedCompletedReceiptList" : finance.ClassifiedCompletedReceiptListAction,
}
