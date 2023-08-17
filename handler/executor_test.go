package handler

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/handler/router"
	"testing"
	"time"
)

type Request struct {
	A string
}

type Response struct {
	B string
}

type SampleHandler2 struct {
	base.Handler
	test string
}

func (c *SampleHandler2) PreHandle(request interface{}) error {
	switch request.(type) {
	case *Request:
		{
			//can change request in PreHandle if necessary
			request.(*Request).A = "changed in pre handle"
			fmt.Printf("*Request:%++v\n", request)
		}
	case Request:
		{
			fmt.Printf("Request:%++v\n", request)
		}
	case *string:
		{
			fmt.Printf("*string:%++v\n", request)
		}
	case *[]byte:
		{
			fmt.Printf("*[]byte:%++v\n", request)
		}
	}
	c.test = "test from pre handle"
	return nil
}

func (c *SampleHandler2) EventHandleMethod1(request *Request) (*Response, *errors.Error) {
	fmt.Printf("---------EventHandleMethod1-------->,request header:%++v,request:%++v, current language:%s, test=%s\n", c.GetRequestHeader(), request, c.Lang, c.test)
	res := &Response{
		B: "this is a test response for format to JSON",
	}
	c.SetResponseHeader(map[string]string{"a": "a1"})
	return res, nil
}

func (c *SampleHandler2) EventHandleMethod2(request Request) (*Response, *errors.Error) {
	fmt.Printf("EventHandleMethod2,request:%++v, lang=%s, test=%s\n", request, c.Lang, c.test)
	res := &Response{
		B: "this is a test response for format to XML",
	}
	return res, nil
}

func (c *SampleHandler2) EventHandleMethod3(request *string) (*string, *errors.Error) {
	fmt.Printf("EventHandleMethod3,request:%s, test=%s\n", *request, c.test)
	str := "this is a test response for format to string"
	return &str, nil
}

func (c *SampleHandler2) EventHandleMethod4(request *[]byte) ([]byte, *errors.Error) {
	fmt.Printf("EventHandleMethod4,request:%++v, test=%s\n", string(*request), c.test)

	return []byte("this is a test response for format to byte[]"), nil
}

func (c *SampleHandler2) EventHandlerMethod5() *errors.Error {
	fmt.Println("EventHandlerMethod5, no request - no response")

	return errors.New("1", "")
}

func (c *SampleHandler2) EventHandlerMethod6() (*string, *errors.Error) {
	fmt.Println("EventHandlerMethod6, no request - nil response")

	return nil, nil
}

func (c *SampleHandler2) OtherMethod(request []byte) (*string, *errors.Error) {

	fmt.Printf("---------otherMethod-------->, request=[%s]\n", string(request))

	return nil, nil
}

// BuildBussinessTopicAttributes is build topic attributes of business service request events
func BuildBussinessTopicAttributes(dstOrg, dstWks, dstEnv, dstSu, dstVersion, dstTopicID string) (map[string]string, error) {
	attributeMap := make(map[string]string)
	if dstTopicID == "" {
		return attributeMap, fmt.Errorf("topic id is empty, please check")
	}

	attributeMap[constant.TopicType] = constant.TopicTypeBusiness
	attributeMap[constant.TopicDestinationORG] = dstOrg
	attributeMap[constant.TopicDestinationWorkspace] = dstWks
	attributeMap[constant.TopicDestinationEnvironment] = dstEnv
	attributeMap[constant.TopicDestinationSU] = dstSu
	attributeMap[constant.TopicDestinationVersion] = dstVersion
	attributeMap[constant.TopicID] = dstTopicID

	return attributeMap, nil
}

func TestInvokeHandler(t *testing.T) {
	callbackExecutor := NewCallbackExecutor()
	routerRegister := &router.HandlerRouter{}

	// JSON encoder/decoder example
	routerRegister.Router("TOPIC1", &SampleHandler2{},
		router.Method("EventHandleMethod1"),
		//router.MarkInvokePreHandle(),
	)
	routerRegister.RouterExpression(".*", &SampleHandler2{},
		router.Method("OtherMethod"),
		router.WithCodec(codec.BuildTextCodec()),
	)
	callbackExecutor.SetRouter(routerRegister)
	callbackExecutor.Init()
	topicAttribute, err := BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC1")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message := &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`{"A":"test1"}`),
	}

	startTime := time.Now()
	for i := 0; i < 1; i++ {
		callbackExecutor.Handle(context.Background(), message)
	}
	fmt.Println("=====>Time cost:", time.Now().Sub(startTime))

	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "XXXXX")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`{"A":"test1"}`),
	}

	startTime = time.Now()
	for i := 0; i < 1; i++ {
		callbackExecutor.Handle(context.Background(), message)
	}

	// XML encoder/decoder example
	routerRegister = &router.HandlerRouter{}
	routerRegister.Router("TOPIC12", &SampleHandler2{},
		router.Method("EventHandleMethod2"),
		router.WithCodec(codec.BuildXMLCodec()),
	)
	callbackExecutor.SetRouter(routerRegister)

	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC12")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`<Request><A>this is a test string2</A></Request>`),
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))

	routerRegister = &router.HandlerRouter{}
	// String encoder/decoder example
	routerRegister.Router("TOPIC13", &SampleHandler2{},
		router.Method("EventHandleMethod3"),
		router.WithCodec(codec.BuildTextCodec()),
	)
	callbackExecutor.SetRouter(routerRegister)
	callbackExecutor.Init()
	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC13")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`this is a test string3`),
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))

	routerRegister = &router.HandlerRouter{}
	// String encoder/decoder example
	routerRegister.Router("TOPIC15", &SampleHandler2{},
		router.Method("EventHandlerMethod5"),
	)
	routerRegister.Router("TOPIC16", &SampleHandler2{},
		router.Method("EventHandlerMethod6"),
	)
	callbackExecutor.SetRouter(routerRegister)
	callbackExecutor.Init()
	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC15")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           nil,
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))

	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC16")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           nil,
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))
	routerRegister = &router.HandlerRouter{}
	// bytes encoder/decoder example
	routerRegister.Router("TOPIC14", &SampleHandler2{},
		router.Method("EventHandleMethod4"),
		router.WithCodec(codec.BuildTextCodec()),
	)
	callbackExecutor.SetRouter(routerRegister)
	callbackExecutor.Init()
	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "TOPIC14")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`this is a test string4`),
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))

	routerRegister = &router.HandlerRouter{}
	routerRegister.RouterPrefix("TOPIC2", &SampleHandler2{}, router.Method("EventHandleMethod2"))
	routerRegister.RouterSuffix("TEST", &SampleHandler2{}, router.Method("EventHandleMethod2"))
	routerRegister.RouterExpression("EXP.*", &SampleHandler2{}, router.Method("EventHandleMethod2"))

	if hp := routerRegister.MatchHandler("TOPIC1"); nil == hp {
		t.Log("TestInvokeHandler failed, cannot found handler with event ID:[TOPIC1]")
	} else {
		if hp.HandlerOptions.HandlerMethodName != "EventHandleMethod1" {
			t.Log("TestInvokeService failed")
		}
	}

	if hp := routerRegister.MatchHandler("TOPIC2"); nil == hp {
		t.Log("TestInvokeHandler failed, cannot found handler with event ID:[TOPIC2]")
	} else {
		if hp.HandlerOptions.HandlerMethodName != "EventHandleMethod2" {
			t.Log("TestInvokeService failed")
		}
	}

	if hp := routerRegister.MatchHandler("TOPIC2_123"); nil == hp {
		t.Log("TestInvokeHandler failed, cannot found handler with event ID:[TOPIC2_123]")
	} else {
		if hp.HandlerOptions.HandlerMethodName != "EventHandleMethod2" {
			t.Log("TestInvokeService failed")
		}
	}

	if hp := routerRegister.MatchHandler("ANY_TOPIC_TEST"); nil == hp {
		t.Log("TestInvokeHandler failed, cannot found handler with event ID:[ANY_TOPIC_TEST]")
	} else {
		if hp.HandlerOptions.HandlerMethodName != "EventHandleMethod2" {
			t.Log("TestInvokeService failed")
		}
	}

	if hp := routerRegister.MatchHandler("EXP1"); nil == hp {
		t.Log("TestInvokeHandler failed, cannot found handler with event ID:[EXP1]")
	} else {
		if hp.HandlerOptions.HandlerMethodName != "EventHandleMethod2" {
			t.Log("TestInvokeService failed")
		}
	}
	callbackExecutor.SetRouter(routerRegister)
	callbackExecutor.Init()
	topicAttribute, err = BuildBussinessTopicAttributes("ORG001", "WKS1", "ENV1", "SU001", "V1", "AAATOPIC1")
	if nil != err {
		t.Errorf("TestInvokeHandler")
	}

	message = &msg.Message{
		ID:             1,
		TopicAttribute: topicAttribute,
		Body:           []byte(`{"A":"test"}`),
	}

	t.Log(callbackExecutor.Handle(context.Background(), message))
}
