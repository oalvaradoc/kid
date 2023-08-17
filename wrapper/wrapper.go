package wrapper

import (
	"context"
)

// Wrapper defines the wrapper operation interface when executing the service or requesting downstream
type Wrapper interface {
	Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error)
	After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error)
}
