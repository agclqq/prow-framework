package retry

import (
	"github.com/agclqq/retry"
)

var std = &retry.Retry{
	InitialBackoff:    1,
	MaxBackoff:        30,
	BackoffMultiplier: 1.5,
	MaxAttempts:       30,
}

func NewRetry(init, max, multi float32, maxStep uint) *retry.Retry {
	return &retry.Retry{
		InitialBackoff:    init,
		MaxBackoff:        max,
		BackoffMultiplier: multi,
		MaxAttempts:       maxStep,
	}
}
func Run(f func(step uint)) {
	std.Run(f)
}
func Reset() {
	std.Reset()
}
func Cancel() {
	std.Cancel()
}
