package util

import "context"

func CheckCanceled(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// ctx가 취소된(또는 타임아웃된) 순간에만 실행
		return ctx.Err()
	default:
		// ctx.Done 채널에 읽을 값(=채널이 닫히는 신호)이 없으면
		// 바로 이곳이 실행되고 함수는 즉시 리턴
		return nil
	}
}
