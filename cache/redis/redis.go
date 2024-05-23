package redis

import (
	"context"
	"time"

	"github.com/spf13/cast"

	"github.com/go-redis/redis/v8"
)

var cacheMap map[string]*Rds = make(map[string]*Rds)

type Rds struct {
	cli *redis.Client
}

type Config struct {
	Type      string
	RdsConfig *redis.Options
}

func New(conf *redis.Options) *Rds {
	connConfig := conf.Addr + conf.Password + cast.ToString(conf.DB)
	if v, ok := cacheMap[connConfig]; ok {
		_, err := v.cli.Ping(context.Background()).Result()
		if err == nil {
			return v
		}
	}
	rdb := redis.NewClient(conf)
	newRds := &Rds{
		cli: rdb,
	}
	cacheMap[connConfig] = newRds
	return newRds
}
func (r *Rds) Origin() interface{} {
	return r.cli
}
func (r *Rds) Get(key string) (string, error) {
	ctx := context.Background()
	return r.GetWithCtx(ctx, key)
}
func (r *Rds) GetWithCtx(ctx context.Context, key string) (string, error) {
	return r.cli.Get(ctx, key).Result()
}

func (r *Rds) Set(key, val string, dur time.Duration) error {
	ctx := context.Background()
	return r.SetWithCtx(ctx, key, val, dur)
}

func (r *Rds) SetWithCtx(ctx context.Context, key, val string, dur time.Duration) error {
	return r.cli.Set(ctx, key, val, dur).Err()
}

func (r *Rds) Forever(key, val string) error {
	ctx := context.Background()
	return r.ForeverWithCtx(ctx, key, val)
}

func (r *Rds) ForeverWithCtx(ctx context.Context, key, val string) error {
	return r.SetWithCtx(ctx, key, val, -1)
}

func (r *Rds) Forget(key string) error {
	ctx := context.Background()
	return r.ForgetWithCtx(ctx, key)
}

func (r *Rds) ForgetWithCtx(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()

}

func (r *Rds) Increment(key string, steps ...int64) (int64, error) {
	ctx := context.Background()
	step := int64(1)
	if len(steps) > 0 {
		step = steps[0]
	}
	return r.cli.IncrBy(ctx, key, step).Result()
}

func (r *Rds) Decrement(key string, steps ...int64) (int64, error) {
	ctx := context.Background()
	step := int64(1)
	if len(steps) > 0 {
		step = steps[0]
	}
	return r.cli.DecrBy(ctx, key, step).Result()
}
