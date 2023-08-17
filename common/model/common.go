package model

// CommonResponse for wrap common response
type CommonResponse struct {
	ErrorCode int         `json:"errorCode"` //-1 indicat error ,0 indicat success
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data"`
}

// BuildErrorResponse is use for build a CommonResponse from error message
func BuildErrorResponse(errMsg string) *CommonResponse {
	response := CommonResponse{
		ErrorCode: -1,
		ErrorMsg:  errMsg,
	}
	return &response
}

// BuildErrorResponseWithErrorCode is use for build a CommonResponse from error code and error message
func BuildErrorResponseWithErrorCode(errCode int, errMsg string) *CommonResponse {
	response := CommonResponse{
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
	}
	return &response
}

// BuildResponse is use for build a CommonResponse from data
func BuildResponse(data interface{}) *CommonResponse {
	response := CommonResponse{
		ErrorCode: 0,
		Data:      data,
	}
	return &response
}

// DataSet is a common struct for organizing an interface data and a total number
type DataSet struct {
	Total int
	Datas interface{}
}

// BuildPagingResponse is use for build a common paging response
func BuildPagingResponse(total int, data interface{}) *CommonResponse {
	dataSet := DataSet{
		Total: total,
		Datas: data,
	}
	return BuildResponse(dataSet)
}
