package client

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/util"
	"time"
)

// BackoffFunc When the transaction fails, the backoff logic will be triggered.
// If the return time of backoff is greater than 0, it will sleep before the next request is initiated.
// This function type is used to agree on the interface that needs to be implemented
// when implementing different backoff strategies.
type BackoffFunc func(ctx context.Context, request Request, attempts int) (time.Duration, error)

// FixedTimeBackoff is the fixed time strategy backoff
// This function will return fixed RetryWaitingTime every time exclude the first request
func FixedTimeBackoff(ctx context.Context, request Request, attempts int) (time.Duration, error) {
	if attempts == 0 || util.IsNil(request) || util.IsNil(request.RequestOptions()) {
		return time.Duration(0), nil
	}
	return request.RequestOptions().RetryWaitingTime, nil
}
