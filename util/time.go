package util

import (
	"context"
	"math/rand"
	"time"
)

func SleepContextMs(ctx context.Context, minMs, maxMs int) {
	delay := rand.Intn(maxMs-minMs+1) + minMs
	SleepContext(ctx, time.Duration(delay)*time.Millisecond)
}

func SleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err() // 취소(또는 타임아웃) 즉시 리턴
	case <-timer.C:
		return nil // 지정된 시간만큼 정상 대기
	}
}

func RandomSleep(minMs, maxMs int) {
	delay := rand.Intn(maxMs-minMs+1) + minMs
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
