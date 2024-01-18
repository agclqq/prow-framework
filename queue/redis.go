package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
type RedisQueue struct {
	cli       *redis.Client
	queueName string
}

func NewRedisQueue(ctx context.Context, config *RedisConfig, queueName string) (*RedisQueue, error) {
	//假设常规连接数在10个
	connNum := 10
	cli := redis.NewClient(&redis.Options{
		Addr:            config.Addr,
		Password:        config.Password,
		DB:              config.DB,
		MinIdleConns:    0,
		MaxIdleConns:    2 * connNum,
		ConnMaxIdleTime: 3 * time.Second,
		ConnMaxLifetime: 30 * time.Second,
	})
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &RedisQueue{cli: cli, queueName: queueName}, nil
}

func (r *RedisQueue) ProduceWithCtx(ctx context.Context, value string) error {
	_, err := r.cli.RPush(ctx, r.queueName, value).Result()
	if err != nil {
		return err
	}

	fmt.Printf("Producing message %s to Redis queue %s\n", value, r.queueName)
	return nil
}

func (r *RedisQueue) ConsumeWithCtx(ctx context.Context) (string, error) {
	value, err := r.cli.BLPop(ctx, 0, r.queueName).Result()
	if err != nil {
		return "", err
	}
	fmt.Printf("Consuming message %s from Redis queue %s\n", value[1], r.queueName)
	return value[1], nil
}

func (r *RedisQueue) ConsumeFuncWithCtx(ctx context.Context, f func(any) bool) error {
	finish := false
	for {
		if finish {
			break
		}
		value, err := r.cli.BLPop(ctx, 0, r.queueName).Result()
		if err != nil {
			return err
		}
		if !f(value[1]) {
			finish = true
		}
	}
	return nil
}
