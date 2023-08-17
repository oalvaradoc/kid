package client

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/apm"
	"git.multiverse.io/eventkit/kit/wrapper"
	"reflect"
	"testing"
)

const MyWrapperName = "MY WRAPPER"

type MyWrapper struct{}

// Before wraps the requests, will generate a new span ID each time.
func (t *MyWrapper) Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error) {
	return ctx, nil
}

// After do nothing
func (t *MyWrapper) After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error) {
	return ctx, nil
}

func (m MyWrapper) String() string {
	return MyWrapperName
}

func TestWithCallWrappers(t *testing.T) {
	c := CallOptions{}

	opt := WithCallWrappers([]wrapper.Wrapper{
		&MyWrapper{},
	})

	opt(&c)

	assert.Equal(t, len(c.CallWrappers), 1)
	nameOfWrapper := fmt.Sprintf("%s", c.CallWrappers[0])
	t.Logf("c.CallWrappers[0]=%s", nameOfWrapper)
	assert.True(t, reflect.DeepEqual(nameOfWrapper, MyWrapperName))
}

func TestDefaultWrapperCall(t *testing.T) {
	c := Options{}
	opt := DefaultWrapperCall([]wrapper.Wrapper{
		&MyWrapper{},
	}...)

	opt(&c)

	assert.Equal(t, len(c.CallOptions.CallWrappers), 1)
	nameOfWrapper := fmt.Sprintf("%s", c.CallOptions.CallWrappers[0])
	t.Logf("c.CallOptions.CallWrappers[0]=%s", nameOfWrapper)
	assert.True(t, reflect.DeepEqual(nameOfWrapper, MyWrapperName))
}

func TestWithCallbackExecutor(t *testing.T) {
	c := CallOptions{}
	opt := WithCallbackExecutor(nil)
	opt(&c)

	assert.Nil(t, c.CallbackExecutor)
}

func TestWithWrapper(t *testing.T) {
	c := CallOptions{}
	opt := WithWrapper([]wrapper.Wrapper{&MyWrapper{}}...)
	opt(&c)

	assert.Equal(t, len(c.CallWrappers), 1)
	nameOfWrapper := fmt.Sprintf("%s", c.CallWrappers[0])
	t.Logf("c.CallWrappers[0]=%s", nameOfWrapper)
	assert.True(t, reflect.DeepEqual(nameOfWrapper, MyWrapperName))
}

func TestWithCallInterceptors(t *testing.T) {
	c := CallOptions{}
	opt := WithCallInterceptors([]interceptor.Interceptor{
		&apm.Interceptor{},
	})

	opt(&c)
	assert.Equal(t, len(c.CallInterceptors), 1)

	nameOfInterceptor := fmt.Sprintf("%s", c.CallInterceptors[0])
	t.Logf("c.CallInterceptors[0]=%s", nameOfInterceptor)
	assert.True(t, reflect.DeepEqual(nameOfInterceptor, constant.InterceptorApm))
}

func TestDefaultCallInterceptors(t *testing.T) {
	c := Options{}
	opt := DefaultCallInterceptors([]interceptor.Interceptor{
		&apm.Interceptor{},
	})

	opt(&c)
	assert.Equal(t, len(c.CallOptions.CallInterceptors), 1)

}

func TestWithOptionFromCallOption(t *testing.T) {
	c := Options{}

	innerOpt := WithCallInterceptors([]interceptor.Interceptor{
		&apm.Interceptor{},
	})

	opt := WithOptionFromCallOption(innerOpt)
	opt(&c)

	assert.Equal(t, len(c.CallOptions.CallInterceptors), 1)
	nameOfInterceptor := fmt.Sprintf("%s", c.CallOptions.CallInterceptors[0])
	t.Logf("c.CallInterceptors[0]=%s", nameOfInterceptor)
	assert.True(t, reflect.DeepEqual(nameOfInterceptor, constant.InterceptorApm))
}

func TestNewOptions(t *testing.T) {
	o := NewOptions()

	assert.NotNil(t, o)
	assert.Equal(t, len(o.CallOptions.CallInterceptors), 0)

	opt := DefaultWrapperCall([]wrapper.Wrapper{
		&MyWrapper{},
	}...)

	o = NewOptions(opt)

	assert.Equal(t, len(o.CallOptions.CallWrappers), 1)
	nameOfWrapper := fmt.Sprintf("%s", o.CallOptions.CallWrappers[0])
	t.Logf("o.CallOptions.CallWrappers[0]=%s", nameOfWrapper)
	assert.True(t, reflect.DeepEqual(nameOfWrapper, MyWrapperName))
}
