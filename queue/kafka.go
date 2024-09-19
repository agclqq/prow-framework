package queue

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	Address string
	Brokers string
	GroupID string
	Topic   string
}

func (k *Kafka) Conn() (*kafka.Conn, error) {
	return kafka.DialContext(context.Background(), "tcp", k.Address)
}
func (k *Kafka) Produce(data []byte) error {
	w := &kafka.Writer{
		Addr:  kafka.TCP(strings.Split(k.Address, ",")...),
		Topic: k.Topic,
	}

	return w.WriteMessages(context.Background(),
		kafka.Message{
			Value: data,
		},
	)
}
func (k *Kafka) ProduceBatch(data [][]byte) error {
	w := &kafka.Writer{
		Addr:  kafka.TCP(strings.Split(k.Address, ",")...),
		Topic: k.Topic,
	}
	d := make([]kafka.Message, 0)
	for _, v := range data {
		d = append(d, kafka.Message{
			Value: v,
		})
	}
	return w.WriteMessages(context.Background(), d...)
}
func (k *Kafka) Consume(f func([]byte) bool) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(k.Brokers, ","),
		GroupID:  k.GroupID,
		Topic:    k.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		rs := f(m.Value)
		if rs {
			r.CommitMessages(ctx, m) // #nosec G104
		}
	}
}

func (k *Kafka) GetPartitions() ([]kafka.Partition, error) {
	conn, err := k.Conn()
	if err != nil {
		return nil, nil
	}
	return conn.ReadPartitions(k.Topic)

}

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}
type KafkaQueue struct {
	conn *kafka.Conn
	r    *kafka.Reader
	w    *kafka.Writer
}

func NewKafkaQueue(ctx context.Context, config *KafkaConfig) (*KafkaQueue, error) {
	conn, err := kafka.DialContext(ctx, "tcp", strings.Join(config.Brokers, ","))
	if err != nil {
		return nil, err
	}
	return &KafkaQueue{
		conn: conn,
		w: &kafka.Writer{
			Addr:  kafka.TCP(config.Brokers...),
			Topic: config.Topic,
		},
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  config.Brokers,
			GroupID:  config.GroupID,
			Topic:    config.Topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}, nil
}

func (kq *KafkaQueue) ProduceWithCtx(ctx context.Context, data string) error {
	return kq.w.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(data),
		},
	)
}

func (kq *KafkaQueue) ConsumeWithCtx(ctx context.Context) (string, error) {
	m, err := kq.r.FetchMessage(ctx)
	if err != nil {
		return "", err
	}
	return string(m.Value), nil
}

func (kq *KafkaQueue) ConsumeFuncWithCtx(ctx context.Context, f func(Message) bool) error {
	i := 0
	for {
		m, err := kq.r.FetchMessage(ctx)
		if err != nil {
			return err
		}

		data := Message{
			Topic:         m.Topic,
			Partition:     m.Partition,
			Offset:        m.Offset,
			HighWaterMark: m.HighWaterMark,
			Value:         m.Value,
			Time:          m.Time,
		}
		if f(data) {
			i++
			if i == 100 {
				err := kq.r.CommitMessages(ctx, m)
				if err != nil {
					return err
				}
				i = 0
			}
		}
	}
}

func (kq *KafkaQueue) GetLag(ctx context.Context) int64 {
	return kq.r.Lag()
}

func (kq *KafkaQueue) GetPartitions(ctx context.Context, topic string) ([]kafka.Partition, error) {
	return kq.conn.ReadPartitions(topic)
}
