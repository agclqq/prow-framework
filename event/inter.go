package event

import "context"

type Eventer interface {
	GetName() string
	GetConcurrence() int64
	Handle(ctx context.Context, data []byte)
}
