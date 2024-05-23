package event

import "context"

type Eventer interface {
	ListenName() string
	Concurrence() int64
	Handle(ctx context.Context, data []byte)
}
