package file

import (
	"context"
	"sync"
	"time"
)

func NewDebouncer(delay time.Duration) func(handler func(ctx context.Context)) func(ctx context.Context) {
	var mu sync.Mutex
	var timer *time.Timer

	return func(handler func(ctx context.Context)) func(ctx context.Context) {
		return func(ctx context.Context) {
			mu.Lock()
			defer mu.Unlock()

			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(delay, func() {
				if ctx.Err() != nil {
					return
				}
				handler(ctx)
			})
		}
	}
}
