package crons

import (
	"sync"
	"time"
)

//配合robfig/cron,实现只在将来的某时刻，只调度一次

type Once struct {
	sync.Mutex
	num     int
	EndTime time.Time
}

func (o *Once) Next(t time.Time) time.Time {

	if o.num > 0 {
		return time.Time{}
	}
	if t.Before(o.EndTime) {
		o.num = 1
		return o.EndTime
	}

	if t.After(o.EndTime) &&
		t.Year() == o.EndTime.Year() &&
		t.Month() == o.EndTime.Month() &&
		t.Day() == o.EndTime.Day() &&
		t.Hour() == o.EndTime.Hour() &&
		t.Minute() == o.EndTime.Minute() {
		o.num = 1
		return o.EndTime
	}
	return time.Time{}
}
