package alert

import (
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type AlertStream struct {
	userAction chan func(ctx context.Context, dialog *page.EventJavascriptDialogOpening)
}

func NewAlertStream() *AlertStream {
	return &AlertStream{
		userAction: make(chan func(ctx context.Context, dialog *page.EventJavascriptDialogOpening), 2), // 버퍼 2
	}
}

func (s *AlertStream) AddListenAlert(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if dialog, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			select {
			case act := <-s.userAction:
				go act(ctx, dialog)
			default:
			}
		}
	})
}

func (s *AlertStream) PushHandler(f func(ctx context.Context, dialog *page.EventJavascriptDialogOpening)) error {
	// 1초 타임아웃 타이머
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case s.userAction <- f:
		// 정상적으로 핸들러를 푸시함
		return nil
	case <-timer.C:
		// 타임아웃
		return errors.New("AlertStream PushHandler timeout")
	}
}
