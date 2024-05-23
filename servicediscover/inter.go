package servicediscover

import "context"

type ServiceDiscovery interface {
	Register
}
type Register interface {
	Register(ctx context.Context, key, val string, ttl int64) error
	UnRegister(ctx context.Context, key string) error
	Renew(ctx context.Context) error
}
type Discover interface {
}
