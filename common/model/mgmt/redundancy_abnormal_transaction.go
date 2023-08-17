package mgmt

import "time"

// DoRedundancyAbnormalTxnRecordRequest is a request model for redundancy abnormal transaction record.
type DoRedundancyAbnormalTxnRecordRequest struct {
	NodeID           string    `json:"nodeId"`
	RootXId          string    `json:"rootXid"`
	MethodName       string    `json:"methodName"`
	MethodFieldsJSON string    `json:"methodFieldsJson"`
	GlobalTxnJSON    string    `json:"globalTxnJson"`
	CreateTime       time.Time `json:"createTime"`
}

// PageQueryRedundancyAbnormalTxnRecordRequest is a request model for redundancy abnormal transaction record paging query
type PageQueryRedundancyAbnormalTxnRecordRequest struct {
	NodeID   string `json:"nodeId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// PageQueryRedundancyAbnormalTxnRecordResult is a response model for redundancy abnormal transaction record paging query
type PageQueryRedundancyAbnormalTxnRecordResult struct {
	NodeID           string    `json:"nodeId"`
	RootXId          string    `json:"rootXid"`
	MethodName       string    `json:"methodName"`
	MethodFieldsJSON string    `json:"methodFieldsJson"`
	GlobalTxnJSON    string    `json:"globalTxnJson"`
	CreateTime       time.Time `json:"createTime"`
}
