package event

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Handle func(context.Context, []byte)
type receiver struct {
	mu             sync.Mutex
	name           string
	handler        Handle //回调方法
	concurrencyNum int64  //当前异步回调并行数量
	maxConcurrency int64  //最大异步回调并行数量
}
type eventMsg struct {
	ctx  context.Context
	data []byte
}
type eventChan struct {
	ch chan *eventMsg
}
type Event struct {
	mu          sync.Mutex
	eventMap    map[string]*eventChan
	receiverMap map[string][]*receiver
}

var ErrNotExistChannel = errors.New("channel does not exist")

var std = &Event{
	eventMap:    make(map[string]*eventChan, 8),
	receiverMap: make(map[string][]*receiver, 8),
}

// InitEvent 初始化事件名称
func InitEvent(name string, capacity int) error {
	std.mu.Lock()
	defer std.mu.Unlock()
	if name == "" {
		return errors.New("the event name cannot be empty")
	}
	if capacity <= 0 {
		capacity = 1000
	}
	if _, ok := std.eventMap[name]; !ok {
		std.eventMap[name] = &eventChan{
			ch: make(chan *eventMsg, capacity),
		}
	} else {
		return errors.New("name is exist")
	}
	return nil
}

// Register 监听者注册
func Register(event Eventer) {
	var concurrentNum int64 = 1
	if event.Concurrence() > 1 {
		concurrentNum = event.Concurrence()
	}
	r := &receiver{
		name:           event.ListenName(),
		handler:        event.Handle,
		concurrencyNum: 0,
		maxConcurrency: concurrentNum,
	}
	std.mu.Lock()
	defer std.mu.Unlock()
	std.receiverMap[event.ListenName()] = append(std.receiverMap[event.ListenName()], r)
}

// Run event开始
func Run() {
	for k, v := range std.eventMap {
		go consumeEventChan(k, v)
	}
}

// 消费单个chan的事件
func consumeEventChan(k string, ec *eventChan) {
	for {
		data := <-ec.ch
		if revs, ok := std.receiverMap[k]; ok {
			for _, rev := range revs { //单个的receiver
				for {
					if rev.concurrencyNum < rev.maxConcurrency {
						go func(v *receiver) {
							atomic.AddInt64(&v.concurrencyNum, 1)
							v.handler(data.ctx, data.data)
							atomic.AddInt64(&v.concurrencyNum, -1)
						}(rev)
						break
					} else {
						time.Sleep(10 * time.Microsecond)
					}
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func Fire(channelName string, data []byte) error {
	ctx := context.Background()
	return FireWithCtx(ctx, channelName, data)
}
func FireWithCtx(ctx context.Context, channelName string, data []byte) error {
	if e, ok := std.eventMap[channelName]; ok {
		e.ch <- &eventMsg{ctx, data}
		return nil
	}
	return ErrNotExistChannel
}
