package interceptor

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/msg"
)

// Interceptor are component that intercept calls to handler methods,
// can be used for auditing and logging as and when handler are accessed
type Interceptor interface {
	PreHandle(ctx context.Context, request *msg.Message) error

	PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error
}
