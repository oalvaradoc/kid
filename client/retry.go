package client

import "context"

// RetryFunc wraps Retry logic when the request fails
type RetryFunc func(ctx context.Context, request Request, retryCount int, err error) (bool, error)

// Always Retry unconditionally
func Always(ctx context.Context, request Request, retryCount int, err error) (bool, error) {
	return true, nil
}

// Never never retry
func Never(ctx context.Context, request Request, retryCount int, err error) (bool, error) {
	return false, nil
}