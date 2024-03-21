package event

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func tearDown() {
	std = &Event{
		eventMap:    make(map[string]*eventChan, 8),
		receiverMap: make(map[string][]*receiver, 8),
	}
}
func Test_InitEnvName(t *testing.T) {
	defer tearDown()
	err := InitEvent("", 100)
	if err == nil {
		t.Errorf("want error, got nil")
	}

	err = InitEvent("test", 100)
	if err != nil {
		t.Errorf("want nil, got %s", err.Error())
	}
	err = InitEvent("test", 100)
	if err == nil {
		t.Errorf("want error, got nil")
	}
}

func Test_longRunning(t *testing.T) {
	defer tearDown()
	err := InitEvent("test", 1)
	if err != nil {
		t.Error(err)
	}

	Register(&longRunning{})
	Run()
	for i := 0; i < 5; i++ {
		err := Fire("test", []byte(strconv.Itoa(i)))
		if err != nil {
			t.Error(err)
		}
	}
	time.Sleep(2 * time.Second)
}
func Test_NotExistChannel(t *testing.T) {
	defer tearDown()
	err := Fire("notExist", []byte("test"))
	if err == nil {
		t.Errorf("want error, got nil")
	}
	if err != ErrNotExistChannel {
		t.Errorf("want %s, got %s", ErrNotExistChannel.Error(), err.Error())
	}
}

type longRunning struct{}

func (l *longRunning) ListenName() string {
	return "test"
}

func (l *longRunning) Concurrence() int64 {
	return 1
}

func (l *longRunning) Handle(ctx context.Context, data []byte) {
	fmt.Println("longRunning", string(data))
	time.Sleep(1 * time.Second)
}

func TestRun(t *testing.T) {
	chanName := make(map[string]int)
	chanName["c1"] = 0
	chanName["c2"] = 100
	chanName["c3"] = 1000
	chanName["c4"] = 5000
	chanName["c5"] = 10000

	for k, v := range chanName {
		err := InitEvent(k, v)
		if err != nil {
			t.Error(err)
			return
		}
	}

	Register(&listen1{})
	Register(&listen2{})
	Register(&listen3{})
	Register(&listen4{})
	Register(&listen5{})
	Register(&listen11{})
	Register(&listen12{})
	Register(&listen13{})
	Register(&listen14{})
	Register(&listen15{})
	Register(&listen21{})
	Register(&listen22{})
	Register(&listen23{})
	Register(&listen24{})
	Register(&listen25{})
	Register(&listen31{})
	Register(&listen32{})
	Register(&listen33{})
	Register(&listen34{})
	Register(&listen35{})
	Run()

	for k := range chanName {
		for i := 0; i < 10; i++ {
			err := Fire(k, []byte(strconv.Itoa(i)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	time.Sleep(1 * time.Second)
}

var todoCtx = context.TODO()

type listen1 struct{}

func (*listen1) ListenName() string {
	return "c1"
}
func (*listen1) Concurrence() int64 {
	return 3
}
func (*listen1) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 1,data: %s\n", string(data))
}

type listen2 struct{}

func (*listen2) ListenName() string {
	return "c2"
}
func (*listen2) Concurrence() int64 {
	return 3
}
func (*listen2) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 2,data: %s\n", string(data))
}

type listen3 struct{}

func (*listen3) ListenName() string {
	return "c3"
}
func (*listen3) Concurrence() int64 {
	return 3
}
func (*listen3) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 3,data: %s\n", string(data))
}

type listen4 struct{}

func (*listen4) ListenName() string {
	return "c4"
}
func (*listen4) Concurrence() int64 {
	return 3
}
func (*listen4) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 4,data: %s\n", string(data))
}

type listen5 struct{}

func (*listen5) ListenName() string {
	return "event:listen5"
}
func (*listen5) Concurrence() int64 {
	return 3
}
func (*listen5) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 5,data: %s\n", string(data))
}

type listen11 struct{}

func (*listen11) ListenName() string {
	return "c1"
}
func (*listen11) Concurrence() int64 {
	return 3
}
func (*listen11) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 11,data: %s\n", string(data))
}

type listen12 struct{}

func (*listen12) ListenName() string {
	return "c2"
}
func (*listen12) Concurrence() int64 {
	return 3
}
func (*listen12) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 12,data: %s\n", string(data))
}

type listen13 struct{}

func (*listen13) ListenName() string {
	return "c3"
}
func (*listen13) Concurrence() int64 {
	return 3
}
func (*listen13) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 13,data: %s\n", string(data))
}

type listen14 struct{}

func (*listen14) ListenName() string {
	return "c4"
}
func (*listen14) Concurrence() int64 {
	return 3
}
func (*listen14) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 14,data: %s\n", string(data))
}

type listen15 struct{}

func (*listen15) ListenName() string {
	return "c5"
}
func (*listen15) Concurrence() int64 {
	return 3
}
func (*listen15) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 15,data: %s\n", string(data))
}

type listen21 struct{}

func (*listen21) ListenName() string {
	return "c1"
}
func (*listen21) Concurrence() int64 {
	return 3
}
func (*listen21) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 21,data: %s\n", string(data))
}

type listen22 struct{}

func (*listen22) ListenName() string {
	return "c2"
}
func (*listen22) Concurrence() int64 {
	return 3
}
func (*listen22) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 22,data: %s\n", string(data))
}

type listen23 struct{}

func (*listen23) ListenName() string {
	return "c3"
}
func (*listen23) Concurrence() int64 {
	return 3
}
func (*listen23) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 23,data: %s\n", string(data))
}

type listen24 struct{}

func (*listen24) ListenName() string {
	return "c4"
}
func (*listen24) Concurrence() int64 {
	return 3
}
func (*listen24) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 24,data: %s\n", string(data))
}

type listen25 struct{}

func (*listen25) ListenName() string {
	return "c5"
}
func (*listen25) Concurrence() int64 {
	return 3
}
func (*listen25) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 25,data: %s\n", string(data))
}

type listen31 struct{}

func (*listen31) ListenName() string {
	return "c1"
}
func (*listen31) Concurrence() int64 {
	return 3
}
func (*listen31) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 31,data: %s\n", string(data))
}

type listen32 struct{}

func (*listen32) ListenName() string {
	return "c2"
}
func (*listen32) Concurrence() int64 {
	return 3
}
func (*listen32) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 32,data: %s\n", string(data))
}

type listen33 struct{}

func (*listen33) ListenName() string {
	return "c3"
}
func (*listen33) Concurrence() int64 {
	return 3
}
func (*listen33) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 33,data: %s\n", string(data))
}

type listen34 struct{}

func (*listen34) ListenName() string {
	return "c4"
}
func (*listen34) Concurrence() int64 {
	return 3
}
func (*listen34) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 34,data: %s\n", string(data))
}

type listen35 struct{}

func (*listen35) ListenName() string {
	return "c5"
}
func (*listen35) Concurrence() int64 {
	return 3
}
func (*listen35) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 31,data: %s\n", string(data))
}
