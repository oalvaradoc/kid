package gls

type TopicInfo struct {
	Type string
	ID   string
}

type ShardingData struct {
	Type  string
	Class string
	ID    string
}

type OptionalDimension struct {
	Organization string
	Workspace    string
	Environment  string
}

type SuInfo struct {
	Type string
	ID   string
}

type ReBindShardingData struct {
	Source ShardingData
	Target ShardingData
}

type SuTypeInfo struct {
	SuType    string
	TopicInfo TopicInfo
	SuList    []string
}

// ShardingDataExistCheckRequest Gls Create Api Entrance parameters
// swagger:parameters shardingDataExistCheckRequest
type ShardingDataExistCheckRequest struct {
	// Dimension Data dimension information, including Topic information and suType fields
	// in: body
	Dimension *DataDimension `json:"dimension"`

	// Element Topic Information Section locates the SuType Code of GLSData, which is generally
	// the topic information of the Event ID that the application service needs to call in the next step.
	// Create GLS Data List
	// in: body
	Element GlsData `json:"element"`
}

// CreateShardingDataRequest Gls Create Api Entrance parameters
// swagger:parameters createShardingDataRequest
type CreateShardingDataRequest struct {
	// Dimension Data dimension information, including Topic information and suType fields
	// in: body
	Dimension *DataDimension `json:"dimension"`

	// RefragmentShardingDataID for mark the current sharding data need to re-fragment
	RefragmentShardingDataID bool `json:"refragmentShardingDataID"`

	// Elements Topic Information Section locates the SuType Code of GLSData, which is generally
	// the topic information of the Event ID that the application service needs to call in the next step.
	// Create GLS Data List
	// in: body
	Elements []GlsData `json:"elements"`
}

// QuerySuListRequest Query Su List Api Entrance parameters
// swagger:parameters querySuListRequest
type QuerySuListRequest struct {
	// Dimension Data dimension information, including Topic information and suType fields
	// in: body
	Dimension *GlsDimension   `json:"dimension"`
	SuTypes   []SuTypeWrapper `json:"suTypes"`
}

// QuerySuListResponse Query Su List Api reqponse
// swagger:parameters querySuListResponse
type QuerySuListResponse struct {
	Dimension *GlsDimension   `json:"dimension"`
	SuTypes   []SuTypeWrapper `json:"suTypes"`
}

// SuTypeWrapper the wrapper for wrap topic type and su type
// swagger:parameters suTypeWrapper
type SuTypeWrapper struct {
	// Topic Data Topic information, including Topic Type and Topic Id.
	// in: body
	Topic GlsTopic `json:"topic"`

	// SuType if contain Topic information, suType is not required.
	//        if Topic information is empty, suType field must be entered.
	// !!!If the topic information is empty, this field must be lost
	SuType string `json:"suType"`

	SuList []string `json:"suList,omitempty"`
}

// DataDimension Data dimension information, including Topic information and suType fields
// swagger:parameters dataDimension
type DataDimension struct {
	// GlsDimension If you do not fill in the ORG/WKS/ENV information, you will default to the ORG/WKS/ENV information from the Mesh message
	// in: body
	GlsDimension

	// Topic Data Topic information, including Topic Type and Topic Id.
	// in: body
	Topic GlsTopic `json:"topic"`

	// SuType if contain Topic information, suType is not required.
	//        if Topic information is empty, suType field must be entered.
	// !!!If the topic information is empty, this field must be lost
	SuType string `json:"suType"`
}

// GlsDimension Dimension with Org/Wks/Env.
// swagger:parameters glsDimension
type GlsDimension struct {
	// Tenant Organization number, not required field
	Tenant string `json:"tenant"`

	// Workspace not required field
	Workspace string `json:"workspace"`

	// Environment not required field
	Environment string `json:"environment"`
}

// GlsTopic Data Topic information, including Topic Type and Topic Id.
// swagger:parameters glsTopic
type GlsTopic struct {
	// TopicType Topic Type, usually use TRN
	TopicType string `json:"topicType" validate:"omitempty"`

	// TopicID Topic Id, the topic that you call the next Event Id
	TopicID string `json:"topicId" validate:"omitempty"`
}

// GlsData GLS sharding element, including type, class, ID
// swagger:parameters glsData
type GlsData struct {
	// GlsType Shard element type, such as ID type = 'Cus'
	GlsType string `json:"glsType"`

	// GlsClass Shard element class, Not required field, if empty, the default value is GlsType,
	// such as ID class = 'IDN'(Id card)、'PAS'（passport)
	GlsClass string `json:"glsClass" validate:"omitempty"`

	// DataID The Id value of the sharding element to bind
	DataID string `json:"dataId"`
}

// GlsPrimarySu Gls Create Api Response parameter
// swagger:parameters glsPrimarySu
type GlsPrimarySu struct {
	// SuType The su type Code of the sharding element, such as The wallet SuType、The Custom SuType...
	SuType string `json:"suType"`

	// SuID The SU code that GLS assigns to the sharding element
	SuID string `json:"suId"`
}

// GlsResponse Gls Update/Remove Api Response parameter
// swagger:parameters glsResponse
type GlsResponse struct {
	// Result true is return successfully Gls Operate
	Result bool `json:"result"`
}

// GlsUpdatePrimaryParameter Gls Update Api Entrance parameters, use this parameters, you can add/update/remove
// the GlsData with The SU Code to which the Primary GlsData belongs
// swagger:parameters glsUpdatePrimaryParameter
type GlsUpdatePrimaryParameter struct {
	// Primary the Primary GlsData, use this info find Su Code
	// in: body
	Primary GlsData `json:"primary"`

	AppendToSU string `json:"appendToSu"`

	RefragmentShardingDataID bool `json:"refragmentShardingDataID"`

	// Dimension Data dimension information, including Topic information and suType fields
	// in: body
	Dimension *DataDimension `json:"dimension"`

	// Elements The Operate Gls Element List
	// in: body
	Elements []GlsUpdateElement `json:"elements"`
}

// GlsUpdateElement GLS Update sharding element
// swagger:parameters glsUpdateElement
type GlsUpdateElement struct {
	// GlsData GLS sharding element, including type, class, ID
	// in:body
	GlsData

	// SourceClass Source Class, Not mandatory field, if value is null, the value is filled with GlsData's GlsClass
	SourceClass string `json:"sourceClass" validate:"omitempty"`

	// SourceID Source Id, Not mandatory field, but if you want to remove the GlsData with the State = 9,
	// To prevent accidental deletion, you must populate the value as GlsData's DataID
	SourceID string `json:"sourceId" validate:"omitempty"`

	// State if State != 9, and the SourceID != DataID, Gls Operate will remove the SourceData(GlsData{GlsType, SourceClass, SourceID}) and Create GlsData
	// Operate State Value, Not mandatory field, default is 0, but if you want to delete the GlsData, you should set state = 9.
	State int `json:"state" validate:"omitempty"` /*  */
}
