package addressing

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/cache/v1"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model/glsdef"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Wrapper is an wrapper for addressing
type Wrapper struct{}


// CheckSuTypeTopic Check Su Type Topic
//
// @param *gls.GlsDimension
// @return gls.Code, error
func CheckSuTypeTopic(ctx context.Context, dim *glsdef.Dimension) *errors.Error {
	if nil == dim {
		return errors.Errorf(constant.InvalidParameters, "Topic && SuType Info is Nil")
	}

	dim.SuType = strings.Trim(dim.SuType, " ")
	if len(dim.SuType) > 0 {
		return nil
	}
	dim.Topic.TopicID = strings.Trim(dim.Topic.TopicID, " ")
	if len(dim.Topic.TopicID) == 0 {
		return errors.Errorf(constant.InvalidParameters, "Topic && SuType Info[%++v] is Null", dim)
	}

	if nil == config.GetConfigs() {
		return errors.Errorf(constant.SystemInternalError, "failed to get service configs")
	}
	var tp string
	var err error
	topicID := randomTopicIDIfNecessary(dim.Topic.TopicID)
	if !config.GetConfigs().Addressing.DisableGLSLookupOptimization {
		tp, err = cache.AddressingCacheOperator.Get(ctx, fmt.Sprintf("%s.%s.%s.%s.%s.%s",
			config.GetConfigs().Addressing.TopicSuTitle,
			dim.Tenant,
			dim.Workspace,
			dim.Environment,
			dim.Topic.TopicType,
			topicID))
	}
	if len(tp) == 0 {
		tp, err = cache.AddressingCacheOperator.HGet(ctx, config.GetConfigs().Addressing.TopicSuTitle, fmt.Sprintf("%s.%s.%s.%s.%s",
			dim.Tenant,
			dim.Workspace,
			dim.Environment,
			dim.Topic.TopicType, topicID))
		if nil != err {
			tp, err = cache.AddressingCacheOperator.HGet(ctx, config.GetConfigs().Addressing.TopicSuTitle, fmt.Sprintf("%s.%s", dim.Topic.TopicType, topicID))
		}
	}
	if err != nil {
		return errors.Errorf(constant.RecordsNotFound, "failed to get Sutype with Dimension[%++v] and topicID:%s: %v", dim, topicID, err)
	}
	dim.SuType = strings.Trim(tp, " ")
	if len(dim.SuType) == 0 {
		return errors.Errorf(constant.RecordsNotFound, "Gls SuType with Dimension[%++v] and topicID:%s is Empty !", dim, topicID)
	}
	return nil
}

func randomElementIDIfNecessary(elementID string) string {
	shardNumber := config.GetConfigs().Addressing.RandomElementIDMap[strings.ToLower(elementID)]
	if shardNumber > 0 {
		return elementID + strconv.Itoa(rand.Int() % shardNumber)
	}

	return elementID
}

func randomTopicIDIfNecessary(topicID string) string {
	shardNumber := config.GetConfigs().Addressing.RandomTopicIDMap[strings.ToLower(topicID)]
	if shardNumber > 0 {
		return topicID + strconv.Itoa(rand.Int() % shardNumber)
	}

	return topicID
}


// Lookup gls Element with dimension info
//
// @param dim *gls.GlsDimension
// @param element gls.Element
// @return pd gls.PrimarySu
// @return code gls.Code
// @return err error
func Lookup(ctx context.Context, dim *glsdef.Dimension, element glsdef.Element) (pd glsdef.PrimarySu, err *errors.Error) {
	pd = glsdef.PrimarySu{}

	element.ElementType = (element.ElementType + "   ")[0:3]
	if err = CheckSuTypeTopic(ctx, dim); err != nil {
		return pd, err
	}
	pd.SuType = dim.SuType
	elementID := randomElementIDIfNecessary(element.ElementID)
	v, gerr := cache.AddressingCacheOperator.HGet(ctx, fmt.Sprintf("CIF.%s.%s.%s.%s.%s.%s",
		dim.Tenant,
		dim.Workspace,
		dim.Environment,
		element.ElementType, element.ElementClass, elementID), dim.SuType)
	if gerr != nil {
		v, gerr = cache.AddressingCacheOperator.HGet(ctx, fmt.Sprintf("%s.%s.%s", element.ElementType, element.ElementClass, elementID), dim.SuType)
	}
	if gerr != nil {
		err = errors.Errorf(constant.RecordsNotFound, "Lookup element type[%s] element class[%s] element id[%s] Dim[%++v] From Redis Faild: %v", element.ElementType, element.ElementClass, elementID, dim, err)
		return pd, err
	}
	pd.SuID = v

	log.Debugsf("Lookup element type[%s] element class[%s] element id[%s]  Element[%++v] result[%++v]", element.ElementType, element.ElementClass, elementID, element, v)
	return pd, nil
}

// Before wraps the requests, will generate a new span ID each time.
func (t *Wrapper) Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		now := time.Now()
		if now.Sub(startTime).Milliseconds() > 10 {
			log.Infof(ctx, "######time cost in `Addressing` wrapper(Before) greater than 10 ms:%++v", now.Sub(startTime))
		}
	}()
	if util.IsNil(opts) {
		log.Info(ctx, "Addressing:request options is empty, skip addressing!")
		return ctx, nil
	}

	requestOptions := opts.(*client.RequestOptions)
	if !requestOptions.IsLocalCall {
		log.Debugf(ctx, "The request is not local call, skip addressing!")
		return ctx, nil
	}

	destinationSU := requestOptions.Su
	if "" != destinationSU {
		return ctx, nil
	}

	requestOptions.HeaderLock.RLock()
	defer requestOptions.HeaderLock.RUnlock()

	// check whether element type exists in topic attributes or not
	elementType := util.GetEither(requestOptions.Header, constant.GlsElementType, constant.GlsElementTypeOld)

	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil == handlerContexts {
		return ctx, errors.Errorf(constant.SystemInternalError, "Cannot found handler contexts in context")
	}

	if "" == elementType {
		requestOptions.Su = handlerContexts.CommonSu
	}

	// lookup cache to get target SU
	dimension := &glsdef.Dimension{
		Tenant:      handlerContexts.Org,
		Workspace:   handlerContexts.Wks,
		Environment: handlerContexts.Env,
		Topic: glsdef.Topic{
			TopicType: constant.TopicTypeBusiness,
			TopicID:   requestOptions.EventID,
		},
	}
	element := glsdef.Element{
		ElementType:  util.GetEither(requestOptions.Header, constant.GlsElementType, constant.GlsElementTypeOld),
		ElementClass: util.GetEither(requestOptions.Header, constant.GlsElementClass, constant.GlsElementClassOld),
		ElementID:    util.GetEither(requestOptions.Header, constant.GlsElementID, constant.GlsElementIDOld),
	}
	if element.ElementClass == "" {
		element.ElementClass = element.ElementType
	}

	if nil == cache.AddressingCacheOperator {
		return ctx, errors.Errorf(constant.SystemInternalError, "Cannot get cache operator")
	}

	// call gls lookup function to find the gls element with its dimension
	if pd, err := Lookup(ctx, dimension, element); nil != err {
		return ctx, err
	} else {
		log.Debugsf("Got SU ID = %s", pd.SuID)
		requestOptions.Su = pd.SuID
	}

	return ctx, nil
}

// After do nothing
func (t *Wrapper) After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error) {
	return ctx, nil
}

func (t Wrapper) String() string {
	return constant.WrapperAddressing
}
