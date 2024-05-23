package queue

import (
	"context"
	"fmt"
	"time"
)

type Message struct {
	// Topic indicates which topic this message was consumed from via Reader.
	//
	// When being used with Writer, this can be used to configure the topic if
	// not already specified on the writer itself.
	Topic string

	// Partition is read-only and MUST NOT be set when writing messages
	Partition     int
	Offset        int64
	HighWaterMark int64
	Value         []byte

	// If not set at the creation, Time will be automatically set when
	// writing the message.
	Time time.Time
}

func (m Message) String() string {
	return fmt.Sprintf("Topic: %s, Partition: %d, Offset: %d, HighWaterMark: %d, Value: %s, Time: %v",
		m.Topic, m.Partition, m.Offset, m.HighWaterMark, string(m.Value), m.Time.Local())
}

// Queue 普通队列
type Queue interface {
	// ProduceWithCtx 生产数据，一次一条
	ProduceWithCtx(ctx context.Context, value string) error
	// ConsumeWithCtx 消费数据，一次一条
	ConsumeWithCtx(ctx context.Context) (string, error)
	// ConsumeFuncWithCtx 通过回调方法持续消
	ConsumeFuncWithCtx(ctx context.Context, f func(Message) bool) error
}

// PriorityQueue 优先级队列
type PriorityQueue interface {
	// ProduceWithCtx 生产数据，一次一条，带优先级
	ProduceWithCtx(ctx context.Context, value string, priority float64) error
	// ConsumeWithCtx 消费数据，一次一条
	ConsumeWithCtx(ctx context.Context) (string, error)
	// ConsumeFuncWithCtx 通过回调方法持续消
	ConsumeFuncWithCtx(ctx context.Context, f func(string) bool) error
}
