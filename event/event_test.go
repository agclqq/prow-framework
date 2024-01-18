package event

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	chanName := make([]string, 0)
	chanName = append(chanName, "c1", "c2", "c3", "c4", "c5")

	InitChannel(chanName...)
	Register("c1", &listen1{}, 3)
	Register("c2", &listen2{}, 3)
	Register("c3", &listen3{}, 3)
	Register("c4", &listen4{}, 3)
	Register("c5", &listen5{}, 3)
	Register("c1", &listen11{}, 3)
	Register("c2", &listen12{}, 3)
	Register("c3", &listen13{}, 3)
	Register("c4", &listen14{}, 3)
	Register("c5", &listen15{}, 3)
	Register("c1", &listen21{}, 3)
	Register("c2", &listen22{}, 3)
	Register("c3", &listen23{}, 3)
	Register("c4", &listen24{}, 3)
	Register("c5", &listen25{}, 3)
	Register("c1", &listen31{}, 3)
	Register("c2", &listen32{}, 3)
	Register("c3", &listen33{}, 3)
	Register("c4", &listen34{}, 3)
	Register("c5", &listen35{}, 3)

	Run()

	for _, v := range chanName {
		for i := 0; i < 10; i++ {
			err := Fire(v, []byte(strconv.Itoa(i)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	time.Sleep(1 * time.Second)
}

var todoCtx = context.TODO()

type listen1 struct{}

func (*listen1) GetName() string {
	return "event:listen1"
}
func (*listen1) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 1,data: %s\n", string(data))
}

type listen2 struct{}

func (*listen2) GetName() string {
	return "event:listen2"
}
func (*listen2) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 2,data: %s\n", string(data))
}

type listen3 struct{}

func (*listen3) GetName() string {
	return "event:listen3"
}
func (*listen3) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 3,data: %s\n", string(data))
}

type listen4 struct{}

func (*listen4) GetName() string {
	return "event:listen4"
}
func (*listen4) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 4,data: %s\n", string(data))
}

type listen5 struct{}

func (*listen5) GetName() string {
	return "event:listen5"
}
func (*listen5) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 5,data: %s\n", string(data))
}

type listen11 struct{}

func (*listen11) GetName() string {
	return "event:listen11"
}
func (*listen11) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 11,data: %s\n", string(data))
}

type listen12 struct{}

func (*listen12) GetName() string {
	return "event:listen12"
}
func (*listen12) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 12,data: %s\n", string(data))
}

type listen13 struct{}

func (*listen13) GetName() string {
	return "event:listen13"
}
func (*listen13) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 13,data: %s\n", string(data))
}

type listen14 struct{}

func (*listen14) GetName() string {
	return "event:listen14"
}
func (*listen14) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 14,data: %s\n", string(data))
}

type listen15 struct{}

func (*listen15) GetName() string {
	return "event:listen15"
}
func (*listen15) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 15,data: %s\n", string(data))
}

type listen21 struct{}

func (*listen21) GetName() string {
	return "event:listen21"
}
func (*listen21) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 21,data: %s\n", string(data))
}

type listen22 struct{}

func (*listen22) GetName() string {
	return "event:listen22"
}
func (*listen22) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 22,data: %s\n", string(data))
}

type listen23 struct{}

func (*listen23) GetName() string {
	return "event:listen23"
}
func (*listen23) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 23,data: %s\n", string(data))
}

type listen24 struct{}

func (*listen24) GetName() string {
	return "event:listen24"
}
func (*listen24) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 24,data: %s\n", string(data))
}

type listen25 struct{}

func (*listen25) GetName() string {
	return "event:listen25"
}
func (*listen25) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 25,data: %s\n", string(data))
}

type listen31 struct{}

func (*listen31) GetName() string {
	return "event:listen31"
}
func (*listen31) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 31,data: %s\n", string(data))
}

type listen32 struct{}

func (*listen32) GetName() string {
	return "event:listen32"
}
func (*listen32) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 32,data: %s\n", string(data))
}

type listen33 struct{}

func (*listen33) GetName() string {
	return "event:listen33"
}
func (*listen33) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 33,data: %s\n", string(data))
}

type listen34 struct{}

func (*listen34) GetName() string {
	return "event:listen34"
}
func (*listen34) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 34,data: %s\n", string(data))
}

type listen35 struct{}

func (*listen35) GetName() string {
	return "event:listen35"
}
func (*listen35) Handle(todoctx context.Context, data []byte) {
	fmt.Printf("this is listen 31,data: %s\n", string(data))
}
