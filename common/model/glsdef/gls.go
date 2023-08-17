package glsdef

// GLSInterface is an interface that contains all GLS interface for eventkit SDK
type GLSInterface interface {
	Looklist(element *Element) ([]PrimarySu, uint, error)
	LookSuType(dim Dimension) (string, uint, error)
	Lookup(dim Dimension, element *Element) (*PrimarySu, uint, error)
	Lookups(dim Dimension, elements []Element) (*LookupsStruct, uint, error)
}

// Topic is a model of GLS API that contains topic type and topicID
type Topic struct {
	TopicType string `json:"topicType"`
	TopicID   string `json:"topicId"`
}

// Dimension is a model of GLS API that contains tenant\workspace\environment\topic\su type
type Dimension struct {
	Tenant      string `json:"tenant"`
	Workspace   string `json:"workspace"`
	Environment string `json:"environment"`
	Topic       Topic  `json:"topic"`
	SuType      string `json:"dcnType"`
}

// LookupArgs is a request model for GLS lookup API
type LookupArgs struct {
	Dimesion *Dimension `json:"dimesion"`
	Element  `json:"element"`
}

// LookupsArgs is a request model for GLS lookups API
type LookupsArgs struct {
	Dimesion *Dimension `json:"dimesion"`
	Elements []Element  `json:"elements"`
}

// LookupsStruct is a response model for GLS lookups API
type LookupsStruct struct {
	Total  uint        `json:"total"`
	Succ   uint        `json:"succ"`
	PrmSus []PrimarySu `json:"prmDcns"`
}

// UpdateElement is a request model for GLS element update API
type UpdateElement struct {
	Element     `json:"element"`
	SourceClass string `json:"sourceClass"`
	SourceID    string `json:"sourceId"`
	State       int    `json:"state"`
}

// Element is a model that contains all element info for GLS API
type Element struct {
	ElementType  string `json:"elementType"`
	ElementClass string `json:"elementClass"`
	ElementID    string `json:"elementId"`
}

// PrimarySu is a model that contains SU information
type PrimarySu struct {
	SuType string `json:"dcnType"`
	SuID   string `json:"dcnId"`
}

// RedisResponse is a response for GLS management API
type RedisResponse struct {
	ErrorCode int         `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Response  RedisConfig `json:"response"`
}

// RedisConfig is a model that contains all redis config information for SED server initialization
type RedisConfig struct {
	Type     string `json:"type"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Readonly bool   `json:"readonly"`
	Poolnum  int    `json:"poolnum"`
	Multi    bool   `json:"multi"`
}
