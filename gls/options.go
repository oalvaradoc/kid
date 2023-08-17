package gls

import "time"

type Options struct {
	RefragmentShardingDataID                bool
	OptionalDimension                       *OptionalDimension
	Timeout                                 time.Duration
	MaxWaitingTime                          time.Duration
	RetryWaitingMilliseconds                time.Duration
	Version                                 string
	MaxRetryTimes                           int
	DeleteTransactionPropagationInformation bool
	TargetEventID                           string
}

type Option func(*Options)

func MarkRefragmentShardingDataID() Option {
	return func(options *Options) {
		options.RefragmentShardingDataID = true
	}
}

func WithOptionalDimension(optionalDimension *OptionalDimension) Option {
	return func(options *Options) {
		options.OptionalDimension = optionalDimension
	}
}

func WithSyncCallTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithSyncCallMaxWaitingTime(maxWaitingTime time.Duration) Option {
	return func(options *Options) {
		options.MaxWaitingTime = maxWaitingTime
	}
}

func WithSyncCallRetryWaitingMilliseconds(retryWaitingMilliseconds time.Duration) Option {
	return func(options *Options) {
		options.RetryWaitingMilliseconds = retryWaitingMilliseconds
	}
}

func WithSyncCallVersion(version string) Option {
	return func(options *Options) {
		options.Version = version
	}
}

func WithSyncCallMaxRetryTimes(maxRetryTimes int) Option {
	return func(options *Options) {
		options.MaxRetryTimes = maxRetryTimes
	}
}

func WithSyncCallDeleteTransactionPropagationInformation(deleteTransactionPropagationInformation bool) Option {
	return func(options *Options) {
		options.DeleteTransactionPropagationInformation = deleteTransactionPropagationInformation
	}
}

func WithSyncCallTargetEventID(targetEventID string) Option {
	return func(options *Options) {
		options.TargetEventID = targetEventID
	}
}
