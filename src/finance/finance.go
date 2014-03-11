package finance

import (
	"github.com/hjqhezgh/commonlib"
	"net/http"
	"strconv"
	"github.com/hjqhezgh/lessgo"
)

func UserInfoAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})
	userId := r.FormValue("userId")

	userInfo := new(UserInfo)

	if userId == "1" {
		userInfo.Id = "1"
		userInfo.Role = "CD"
	} else {
		userInfo.Id = "0"
		userInfo.Role = "C"
	}

	m["success"] = true
	m["code"] = "200"
	m["datas"] = userInfo
	commonlib.OutputJson(w, m, "")

	return
}

func HandleApplyAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})
	userId := r.FormValue("userId")

	role := GetUserRole(userId)

	handleApply := new(HandleApply)
	handleApply.PendingReceiptAmount = GetPendingReceiptAmount(role)
	handleApply.CompletedReceiptAmount = GetCompletedReceiptAmount(role)

	m["success"] = true
	m["code"] = "200"
	m["datas"] = handleApply
	commonlib.OutputJson(w, m, "")

	return
}

func ReceiptDetailsAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})
	userId := r.FormValue("userId")

	role := GetUserRole(userId)

	receiptList := GetReceiptList(role)

	m["success"] = true
	m["code"] = "200"
	m["datas"] = receiptList
	commonlib.OutputJson(w, m, "")

	return
}

func ClassifiedPendingReceiptListAction(w http.ResponseWriter, r *http.Request) {
	ClassifiedReceiptListAction(w, r)
}

func ClassifiedCompletedReceiptListAction(w http.ResponseWriter, r *http.Request) {
	ClassifiedReceiptListAction(w, r)
}

func ClassifiedReceiptListAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})
	userId := "1"//CD
	//	userId := "2"//管理
	role := GetUserRole(userId)

	receiptType := r.FormValue("status")
	week := r.FormValue("week")
	code := r.FormValue("code")
	lessgo.Log.Info("receiptType: ", receiptType, "  week: ", week, "  code: ", code)

	receiptList := getClassifiedReceiptList(role, receiptType)

	m["success"] = true
	m["code"] = "200"
	m["datas"] = receiptList
	commonlib.OutputJson(w, m, "")

	return
}

func getClassifiedReceiptList(role, receiptType string) []*Receipt {
	var receiptList []*Receipt
	if receiptType == "A" {
		if role == "CD" {
			// getCdData
		} else {
			// getOtherData
		}
	} else if receiptType == "P" {

	}

	return receiptList
}

func GetUserRole(userId string) string {
	var userRole string
	if userId == "1" {
		userRole = "CD"
	} else {
		userRole = "C"
	}

	return userRole
}

func GetPendingReceiptAmount(role string) string {
	var pendingReceiptAmount string
	if role == "CD" {
		pendingReceiptAmount = "2"
	} else {
		pendingReceiptAmount = "10"
	}

	return pendingReceiptAmount
}

func GetCompletedReceiptAmount(role string) string {
	var completeReceiptAmount string
	if role == "CD" {
		completeReceiptAmount = "20"
	} else {
		completeReceiptAmount = "100"
	}

	return completeReceiptAmount
}

func GetPendingReceiptDetails(role string) {
	if role == "CD" {
		//
	} else {
		//
	}
}

func GetCompletedReceiptDetails(role string) {
	if role == "CD" {
		//
	} else {
		//
	}
}

func GetReceiptList(role string) []*Receipt {
	var receiptList []*Receipt
	lessgo.Log.Debug(role)
	if role == "CD" {
		code := 100100
		curCode := code
		curAmount := 500
		for i := 0; i < 100; i++ {
			receipt := new(Receipt)
			receipt.Id = strconv.Itoa(i + 1)
			receipt.ApplyDate = "2014-03-02 12:00:00"
			if i%5 == 0 {
				curCode = code + i
			}
			receipt.AccountingSubjectCode = strconv.Itoa(curCode)
			receipt.AccountingSubjectName = "差旅费" + strconv.Itoa(curCode)
			receipt.ReceiptAmount = strconv.Itoa(curAmount + i*4)
			if i%3 == 0 {
				receipt.ApproveStatus = "A"// approved
			} else {
				receipt.ApproveStatus = "P"// pending
			}

			receiptList = append(receiptList, receipt)
		}
		for i := 0; i < 100; i++ {
			receipt := new(Receipt)
			receipt.Id = strconv.Itoa(i + 1)
			receipt.ApplyDate = "2014-03-5 12:00:00"
			if i%5 == 0 {
				curCode = code + i
			}
			receipt.AccountingSubjectCode = strconv.Itoa(curCode)
			receipt.AccountingSubjectName = "差旅费" + strconv.Itoa(curCode)
			receipt.ReceiptAmount = strconv.Itoa(curAmount + i*7)
			if i%3 == 0 {
				receipt.ApproveStatus = "A"// approved
			} else {
				receipt.ApproveStatus = "P"// pending
			}

			receiptList = append(receiptList, receipt)
		}
		for i := 0; i < 100; i++ {
			receipt := new(Receipt)
			receipt.Id = strconv.Itoa(i + 1)
			receipt.ApplyDate = "2014-03-10 12:00:00"
			if i%5 == 0 {
				curCode = code + i
			}
			receipt.AccountingSubjectCode = strconv.Itoa(curCode)
			receipt.AccountingSubjectName = "差旅费" + strconv.Itoa(curCode)
			receipt.ReceiptAmount = strconv.Itoa(curAmount + i*13)
			if i%3 == 0 {
				receipt.ApproveStatus = "A"// approved
			} else {
				receipt.ApproveStatus = "P"// pending
			}

			receiptList = append(receiptList, receipt)
		}
	} else {
		//
	}

	return receiptList
}
