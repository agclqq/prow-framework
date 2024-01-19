package crons

import (
	"fmt"
	"testing"
	"time"

	"github.com/agclqq/prow-framework/times"

	"github.com/robfig/cron/v3"
)

func TestOncePastTime(t *testing.T) {
	num := 0
	now := time.Now()
	duration, err := time.ParseDuration("-1m")
	if err != nil {
		t.Errorf(err.Error())
	}
	nowTime := now.Add(duration)
	c := cron.New()
	once := &Once{EndTime: nowTime}
	c.Schedule(once, cron.FuncJob(func() { num++; fmt.Printf("我只执行一次at%s", time.Now().Format(times.FormatDatetimeMicro)) }))
	c.Start()
	fmt.Printf("now time:\t%s,\nexec time:\t%s\n", now.Format(times.FormatDatetimeMicro), nowTime.Format(times.FormatDatetimeMicro))
	time.Sleep(60 * time.Second)
	if num > 0 {
		t.Errorf("want:0 got:%d", num)
	}
}
func TestOnceFutureTime(t *testing.T) {
	num := 0
	now := time.Now()
	duration, err := time.ParseDuration("1m")
	if err != nil {
		t.Errorf(err.Error())
	}
	nowTime := now.Add(duration)
	c := cron.New()
	once := &Once{EndTime: nowTime}
	c.Schedule(once, cron.FuncJob(func() { num++; fmt.Printf("我只执行一次at%s", time.Now().Format(times.FormatDatetimeMicro)) }))
	c.Start()
	fmt.Printf("now time:\t%s,\nexec time:\t%s\n", now.Format(times.FormatDatetimeMicro), nowTime.Format(times.FormatDatetimeMicro))
	time.Sleep(duration + 10*time.Second)
	if num > 1 {
		t.Errorf("want:1 got:%d", num)
	}
}
