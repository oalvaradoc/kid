package router

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/handler/base"
	"testing"
)

type MyHandler struct {
	base.Handler
}

func (m MyHandler) PreHandle() {}

func (m MyHandler) HandlerName() string {
	return "MyHandler"
}

func (m *MyHandler) Method1() {}

func (m MyHandler) Method2() error { return nil }

func (m *MyHandler) Method3() error { return nil }

func TestHandlerRouter_RouterWithEventKey(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.Router("eventID1", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey2(t *testing.T) {
	router := &HandlerRouter{}
	router.SetEventKeyMap(map[string]string{
		"eventKey1": "eventID11",
		"eventKey2": "eventID12",
		"eventKey3": "eventID13",
	})

	router.SetDefaultInterceptors(defaultInterceptors)
	router.DefaultEnableValidation(true, nil)
	router.Router("eventID1", &MyHandler{}, Method("Method1"), HandlePost("/v1/method1"))
	router.Router("eventID2", &MyHandler{}, Method("Method2"), HandlePost("/v1/method2"))
	router.Router("eventID3", &MyHandler{}, Method("Method3"), HandlePost("/v1/method3"))

	router.RouterWithEventKey("eventKey1", &MyHandler{}, Method("Method1"), HandlePost("/v1/method11"))
	router.RouterWithEventKey("eventKey2", &MyHandler{}, Method("Method2"), HandlePost("/v1/method12"))
	router.RouterWithEventKey("eventKey3", &MyHandler{}, Method("Method3"), HandlePost("/v1/method13"))

	router.RouterPrefix("Prefix1", &MyHandler{}, Method("Method1"), HandlePost("/v1/method21"))
	router.RouterPrefix("Prefix2", &MyHandler{}, Method("Method2"), HandlePost("/v1/method22"))
	router.RouterPrefix("Prefix3", &MyHandler{}, Method("Method3"), HandlePost("/v1/method23"))

	router.RouterSuffix("Suffix1", &MyHandler{}, Method("Method1"), HandlePost("/v1/method31"))
	router.RouterSuffix("Suffix2", &MyHandler{}, Method("Method2"), HandlePost("/v1/method32"))
	router.RouterSuffix("Suffix3", &MyHandler{}, Method("Method3"), HandlePost("/v1/method33"))

	router.RouterExpression("Test1*", &MyHandler{}, Method("Method1"), HandlePost("/v1/method41"))
	router.RouterExpression("My2*", &MyHandler{}, Method("Method2"), HandlePost("/v1/method42"))
	router.RouterExpression("Expression3*", &MyHandler{}, Method("Method3"), HandlePost("/v1/method43"))
	o := router.MatchHandler("does not exist")
	assert.Nil(t, o)

	o = router.MatchHandler("eventID1")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method1")

	o = router.MatchHandler("eventID11")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method1")

	o = router.MatchHandler("Prefix3_123")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method3")

	o = router.MatchHandler("testSuffix2")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method2")

	o = router.MatchHandler("Expression3")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method3")

	o = router.MatchHandlerWithURLPath("/v1/path_does_not_exists")
	assert.Nil(t, o)

	o = router.MatchHandlerWithURLPath("/v1/method1")
	assert.NotNil(t, o)
	t.Logf("The option of router:%++v", o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "Method1")
}

func TestHandlerRouter_RouterWithEventKey3(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.Router("eventID1", MyHandler{}, Method("Method1"))
}

func TestHandlerRouter_RouterWithEventKey4(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterPrefix("event", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey5(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterPrefix("event", MyHandler{}, Method("Method1"))
}

func TestHandlerRouter_RouterWithEventKey6(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterSuffix("event", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey7(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterSuffix("event", MyHandler{}, Method("Method1"))
}

func TestHandlerRouter_RouterWithEventKey8(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterExpression("event*", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey9(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterExpression("event*", MyHandler{}, Method("Method1"))
}

func TestHandlerRouter_RouterWithEventKey10(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.RouterWithEventKey("eventKey", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey11(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.EventKeyMap = map[string]string{"eventKey": "testEventID"}
	router.RouterWithEventKey("eventKey", &MyHandler{}, Method("does not exist!"))
}

func TestHandlerRouter_RouterWithEventKey12(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Logf("failed to register router, error:%++v\n", e)
		}
	}()
	router := &HandlerRouter{}
	router.EventKeyMap = map[string]string{"eventKey": "testEventID"}
	router.RouterWithEventKey("eventKey", MyHandler{}, Method("Method1"))
}
