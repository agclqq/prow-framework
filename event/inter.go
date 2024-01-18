package event

import "context"

type Eventer interface {
	GetName() string
	Handle(ctx context.Context, data []byte)
}
