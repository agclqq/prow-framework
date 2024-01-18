package queue

import (
	"context"
	"testing"
)

func TestRedisQueue(t *testing.T) {
	ctx := context.Background()
	redisQueue, err := NewRedisQueue(ctx, &RedisConfig{
		Addr:     "10.134.116.173:6379",
		Password: "5|Y;SDC:hKXs_EgN>4",
		DB:       1,
	}, "test_queue")
	if err != nil {
		t.Fatalf("Error creating Redis priority queue: %v", err)
	}

	// Produce two messages to Redis priority queue
	err = redisQueue.ProduceWithCtx(ctx, "hello world")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	err = redisQueue.ProduceWithCtx(ctx, "hi there")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	err = redisQueue.ProduceWithCtx(ctx, "it's over")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	data, err := redisQueue.ConsumeWithCtx(ctx)
	if err != nil {
		t.Fatalf("Error consumeing from Redis queue: %v", err)
		return
	}
	if data != "hello world" {
		t.Fatalf("Expected message 'hello world' from Redis queue, but got '%s'", data)
	}
	rs := make([]string, 0)
	// Consume a message from Redis priority queue
	err = redisQueue.ConsumeFuncWithCtx(ctx, func(data any) bool {
		if d, ok := data.(string); ok {
			rs = append(rs, d)
		}
		if len(rs) == 2 {
			return false
		}
		return true
	})
	if err != nil {
		t.Fatalf("Error consuming from Redis queue: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("Expected message num 2 from Redis queue, but got '%d'", len(rs))
	}
	if rs[0] != "hi there" || rs[1] != "it's over" {
		t.Fatalf("Expected message list 'hi there,it's over' from Redis queue, but got '%s,%s'", rs[0], rs[1])
	}
}
