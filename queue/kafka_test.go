package queue

import (
	"context"
	"testing"
)

func TestKafkaQueue(t *testing.T) {
	t.Skip("此测试函数不会被执行")
	ctx := context.Background()
	kafkaQueue, _ := NewKafkaQueue(ctx, &KafkaConfig{
		Brokers: []string{"10.134.192.6:9092"},
		GroupID: "test-2",
		Topic:   "workflow_step_status_changed_dev",
	})

	// Produce two messages to Redis priority queue
	err := kafkaQueue.ProduceWithCtx(ctx, "hello world")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	err = kafkaQueue.ProduceWithCtx(ctx, "hi there")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	err = kafkaQueue.ProduceWithCtx(ctx, "it's over")
	if err != nil {
		t.Fatalf("Error producing to Redis queue: %v", err)
	}
	data, err := kafkaQueue.ConsumeWithCtx(ctx)
	if err != nil {
		t.Fatalf("Error consumeing from Redis queue: %v", err)
		return
	}
	if data != "hello world" {
		t.Fatalf("Expected message 'hello world' from Redis queue, but got '%s'", data)
	}
	rs := make([]string, 0)
	// Consume a message from Redis priority queue
	err = kafkaQueue.ConsumeFuncWithCtx(ctx, func(data Message) bool {
		rs = append(rs, string(data.Value))
		if len(rs) == 2 {
			if rs[0] != "hi there" || rs[1] != "it's over" {
				t.Fatalf("Expected message list 'hi there,it's over' from Redis queue, but got '%s,%s'", rs[0], rs[1])
			}
		}
		return true
	})
	if err != nil {
		t.Fatalf("Error consuming from Redis queue: %v", err)
	}

}
