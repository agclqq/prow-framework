package event

import "context"

type Eventer interface {
	GetName() string
	GetConcurrence() int32
	Handle(ctx context.Context, data []byte)
}
