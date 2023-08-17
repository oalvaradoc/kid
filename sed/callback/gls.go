package callback

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/json"
	"git.multiverse.io/eventkit/kit/common/model/glsdef"
	"git.multiverse.io/eventkit/kit/common/try"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"strings"
)

const (
	// LookListPath is an url path for gls looklist
	LookListPath = "/v1/gls/looklist"
	// LookSuTypePath is an url path for gls looksutype
	LookSuTypePath = "/v1/gls/looksutype"
	// LookupPath is an url path for gls lookup
	LookupPath = "/v1/gls/lookup"
	// LookupsPath is an url path for gls lookups
	LookupsPath = "/v1/gls/lookups"
)

//define Return Code
var (
	Successful uint = 0
	ParamNull  uint = 661400
	NotFound   uint = 661403
	SysErr     uint = 669999
)

// PrimarySuResult define PrimarySu Result Struct
type PrimarySuResult struct {
	Code uint
	*glsdef.PrimarySu
}

// ReplyLookupPrimaryList define Reply Lookup Primary List Struct
type ReplyLookupPrimaryList struct {
	Code    uint
	Message string
	Data    []glsdef.PrimarySu
}

// ReplyLookSutypePrimary define reply Look Sutype Primary Struct
type ReplyLookSutypePrimary struct {
	Code    uint
	Message string
	Data    string
}

// ReplyLookupPrimary define Reply Lookup Primary Struct
type ReplyLookupPrimary struct {
	Code    uint
	Message string
	Data    *PrimarySuResult
}

// ReplyLookupsPrimary define Reply Lookups Primary Struct
type ReplyLookupsPrimary struct {
	Code    uint
	Message string
	Data    *glsdef.LookupsStruct
}

// CreatePrimaryStruct define Create Primary Struct
type CreatePrimaryStruct struct {
	Option   int
	Dimesion *glsdef.Dimension
	Elements []glsdef.Element
}

// ExistStruct define check Exist Struct
type ExistStruct struct {
	Dimesion glsdef.Dimension
	Element  glsdef.Element
}

// GlsOperate a default implementer that implement the GLSInterface
type GlsOperate struct{}

// Looklist Search Su List by Element
func (op *GlsOperate) Looklist(element *glsdef.Element) ([]glsdef.PrimarySu, uint, error) {
	// 校验参数
	s, _ := json.Marshal(element)
	log.Debugsf("Looklist elements: %s", string(s))
	if len(strings.Trim(element.ElementID, " ")) == 0 || len(strings.Trim(element.ElementType, " ")) == 0 {
		return nil, ParamNull, errors.Errorf(constant.SystemInternalError, "Looklist Lookup Elements Paramer Error !")
	}

	rly := ReplyLookupPrimaryList{}
	if err := postGLSUserRequest(LookListPath, element, &rly); err != nil {
		log.Debugsf("Looklist Faild: err=%v", err)
		return nil, SysErr, err
	}
	return rly.Data, rly.Code, nil
}

// LookSuType lookup SuType by Topic Info
func (op *GlsOperate) LookSuType(dim glsdef.Dimension) (string, uint, error) {
	s, err := json.Marshal(dim)
	if err != nil {
		log.Debugsf("LookSuType Marshal[%v] Faild: err=%v", dim, err)
		fmt.Println(err)
		return "", SysErr, err
	}
	log.Debugsf("LookSuType: %s", string(s))
	//if 0 == len(strings.Trim(dim.Topic.TopicType, " ")) && 0 == len(strings.Trim(dim.Topic.TopicType, " ")) {
	if 0 == len(strings.Trim(dim.Topic.TopicType, " ")) && 0 == len(strings.Trim(dim.Topic.TopicID, " ")) {
		return "", ParamNull, err
	}

	rly := ReplyLookSutypePrimary{}
	if err := postGLSUserRequest(LookSuTypePath, dim, &rly); err != nil {
		log.Debugsf("LookSuType Faild: err=%v", err)
		return "", SysErr, err
	}
	if rly.Code == 0 {
		return rly.Data, rly.Code, nil
	}

	return rly.Data, rly.Code, errors.Errorf(constant.SystemInternalError, rly.Message)
}

func postGLSUserRequest(path string, req interface{}, resp interface{}) error {
	err := glsPostRequest(path, req, resp)
	switch err.(type) {
	case *errors.Error:
		{
			if err.(*errors.Error) != nil {
				err = err.(*errors.Error).Err
			} else {
				err = nil
			}
		}
	default:
		err = nil
	}
	return err
}

// Lookup Search Su by Element and Topic Info
func (op *GlsOperate) Lookup(dim glsdef.Dimension, element *glsdef.Element) (*glsdef.PrimarySu, uint, error) {
	// check Parameters
	s, _ := json.Marshal(element)
	log.Debugsf("Lookup elements: %s", string(s))
	if len(strings.Trim(element.ElementID, " ")) == 0 || len(strings.Trim(element.ElementType, " ")) == 0 {
		return nil, ParamNull, errors.Errorf(constant.SystemInternalError, "Lookup Elements Paramer Error !")
	}

	rly := ReplyLookupPrimary{}
	if err := postGLSUserRequest(LookupPath, glsdef.LookupArgs{
		Dimesion: &dim,
		Element:  *element,
	}, &rly); err != nil {
		log.Debugsf("Lookup Faild: err=%v", err)
		return nil, SysErr, err
	}
	if rly.Code == 0 {
		return rly.Data.PrimarySu, rly.Code, nil
	}

	return nil, rly.Code, errors.Errorf(constant.SystemInternalError, rly.Message)
}

// Lookups Search Su List by Element List and Topic Info
func (op *GlsOperate) Lookups(dim glsdef.Dimension, elements []glsdef.Element) (*glsdef.LookupsStruct, uint, error) {
	// check Parameters
	s, _ := json.Marshal(elements)
	log.Debugsf("Lookups elements: %v", string(s))
	for _, element := range elements {
		if len(strings.Trim(element.ElementID, " ")) == 0 || len(strings.Trim(element.ElementType, " ")) == 0 {
			return nil, ParamNull, errors.Errorf(constant.SystemInternalError, "Lookups, Lookup Elements Paramer Error !")
		}
	}

	rly := ReplyLookupsPrimary{}
	if err := postGLSUserRequest(LookupsPath, glsdef.LookupsArgs{
		Dimesion: &dim,
		Elements: elements,
	}, &rly); err != nil {
		log.Debugsf("Lookups Faild: err=%v", err)
		return nil, SysErr, err
	}
	if rly.Code == try.SuccCode {
		return rly.Data, rly.Code, nil
	}

	return rly.Data, rly.Code, errors.Errorf(constant.SystemInternalError, rly.Message)
}

// NewGlsOperate is used to create a new GLS operator
func NewGlsOperate() glsdef.GLSInterface {
	return &GlsOperate{}
}

// Looklist  is used to lookup multiple SU with an element
func Looklist(element *glsdef.Element) ([]glsdef.PrimarySu, uint, error) {
	operate := &GlsOperate{}
	return operate.Looklist(element)
}

// LookSuType is used to lookup su type
func LookSuType(dim glsdef.Dimension) (string, uint, error) {
	operate := &GlsOperate{}
	return operate.LookSuType(dim)
}

// Lookup is used to lookup SU with element information.
func Lookup(dim glsdef.Dimension, element *glsdef.Element) (*glsdef.PrimarySu, uint, error) {
	operate := &GlsOperate{}
	return operate.Lookup(dim, element)
}

// Lookups is used to lookup multiple SU with list of element
func Lookups(dim glsdef.Dimension, elements []glsdef.Element) (*glsdef.LookupsStruct, uint, error) {
	operate := &GlsOperate{}
	return operate.Lookups(dim, elements)
}
