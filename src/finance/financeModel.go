package finance

type UserInfo struct {
	Id	                          		string			`json:"id"`
	Role			               		string			`json:"role"`
}

type HandleApply struct {
	PendingReceiptAmount				string			`json:"pendingReceiptAmount"`
	CompletedReceiptAmount				string			`json:"completedReceiptAmount"`
}

type Receipt struct {
	Id									string			`json:"id"`
	ApplyDate							string			`json:"applyDate"`
	AccountingSubjectCode				string			`json:"accountingSubjectCode"`
	AccountingSubjectName				string			`json:"accountingSubjectName"`
	ReceiptAmount						string			`json:"receiptAmount"`
	ApproveStatus						string			`json:"approveStatus"`
}
