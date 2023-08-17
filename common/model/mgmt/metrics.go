package mgmt

// MetricInfo is a model that contains the TPS\AvgTimeCost\SuccessRate information for metric logic.
type MetricInfo struct {
	TPS         float64 `json:"tps"`
	AvgTimeCost float64 `json:"avgTimeCost"`
	SuccessRate float64 `json:"successRate"`
}

// Metrics is a model that contains all the metric information of DXC server
type Metrics struct {
	TxnBegin        MetricInfo `json:"txnBegin"`
	TxnJoin         MetricInfo `json:"txnJoin"`
	TxnEnd          MetricInfo `json:"txnEnd"`
	CallbackConfirm MetricInfo `json:"callbackConfirm"`
	CallbackCancel  MetricInfo `json:"callbackCancel"`
	AvgTxnNum       float64    `json:"avgTxnNum"`
	MetricsTime     string     `json:"metricsTime"`
}

// DoMetricsRequest is a request model for DXC server report metric information.
type DoMetricsRequest struct {
	Metrics     []Metrics `json:"metrics"`
	RequestTime string    `json:"requestTime"`
}
