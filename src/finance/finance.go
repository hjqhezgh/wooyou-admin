package finance

import (
	"github.com/hjqhezgh/commonlib"
	"github.com/hjqhezgh/lessgo"
	"net/http"
	"strconv"
	"tool"
	"strings"
)

func RoleCodesAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})

	roleCodes, err := tool.GetCurrentEmployeeRoles(r)
	if err != nil {
		lessgo.Log.Error(err)
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = "200"
	m["datas"] = roleCodes
	commonlib.OutputJson(w, m, "")

	return
}

func HandleApplyAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})

	roleCodes, err := tool.GetCurrentEmployeeRoles(r)
	if err != nil {
		lessgo.Log.Error(err)
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	handleApply := new(HandleApply)
	handleApply.PendingReceiptAmount = GetPendingReceiptAmount(roleCodes)
	handleApply.CompletedReceiptAmount = GetCompletedReceiptAmount(roleCodes)

	m["success"] = true
	m["code"] = "200"
	m["datas"] = handleApply
	commonlib.OutputJson(w, m, "")

	return
}

func ReceiptDetailsAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface {})

	roleCodes, err := tool.GetCurrentEmployeeRoles(r)
	if err != nil {
		lessgo.Log.Error(err)
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	receiptList := GetReceiptList(roleCodes)

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

func GetPendingReceiptAmount(roleCodes []string) string {
	var pendingReceiptAmount string
	if isStringExistInStringArray(roleCodes, "cd") {
		pendingReceiptAmount = "2"
	} else if isStringExistInStringArray(roleCodes, "yyzj") {
		pendingReceiptAmount = "10"
	}

	return pendingReceiptAmount
}

func GetCompletedReceiptAmount(roleCodes []string) string {
	var completeReceiptAmount string
	if isStringExistInStringArray(roleCodes, "cd") {
		completeReceiptAmount = "20"
	} else if isStringExistInStringArray(roleCodes, "yyzj") {
		completeReceiptAmount = "100"
	}

	return completeReceiptAmount
}

func GetPendingReceiptDetails(roleCodes []string) {
	if isStringExistInStringArray(roleCodes, "cd") {
		//
	} else if isStringExistInStringArray(roleCodes, "yyzj") {
		//
	}
}

func GetCompletedReceiptDetails(roleCodes []string) {
	if isStringExistInStringArray(roleCodes, "cd") {
		//
	} else if isStringExistInStringArray(roleCodes, "yyzj") {
		//
	}
}

func GetReceiptList(roleCodes []string) []*Receipt {
	var receiptList []*Receipt
	if isStringExistInStringArray(roleCodes, "cd") {
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
	} else if isStringExistInStringArray(roleCodes, "yyzj") {
		//
	}

	return receiptList
}

func isStringExistInStringArray (strArray []string, str string) bool {
	isExist := false
	for _, v := range strArray {
		if strings.ToLower(str) == strings.ToLower(v) {
			isExist = true
			break
		}
	}

	return isExist
}
