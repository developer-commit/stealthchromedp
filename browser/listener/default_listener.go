// browser/default_listener.go
package listener

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/developer-commit/stealthchromedp/browser/alert"
	"github.com/developer-commit/stealthchromedp/browser/driver"
	"github.com/developer-commit/stealthchromedp/browser/tabmanager"
	"github.com/developer-commit/stealthchromedp/runner"
)

// DefaultListener는 BrowserListener 인터페이스의 기본 구현체입니다.
type DefaultListener struct {
	drv    *driver.ChromeDpDriver
	tm     *tabmanager.TabManager
	runner runner.TaskRunner
	alert  *alert.AlertStream
}

func NewDefaultListener(
	drv *driver.ChromeDpDriver,
	tm *tabmanager.TabManager,
	runner runner.TaskRunner,
	alertStream *alert.AlertStream,
) BrowserListener {
	return &DefaultListener{drv: drv, tm: tm, runner: runner, alert: alertStream}
}

func (l *DefaultListener) ListenNewTab(ctx context.Context) {
	chromedp.ListenBrowser(ctx, func(ev interface{}) {
		e, ok := ev.(*target.EventTargetCreated)
		if !ok || e.TargetInfo.OpenerID == "" || e.TargetInfo.Type != "page" {
			return
		}
		go func(tid target.ID) {
			tab, ok := l.tm.Now()
			if !ok {
				return
			}
			popupCtx, popupCancel := l.drv.BuildTabByID(tab.Ctx, tid)
			l.tm.AddNewTab(tid, popupCtx, popupCancel)
			if err := l.runner.Run(popupCtx); err != nil {
				popupCancel()
				fmt.Println(err)
			}
			l.ListenNewTab(popupCtx)
			l.ListenCloseTab(popupCtx)
			l.alert.AddListenAlert(popupCtx)
		}(e.TargetInfo.TargetID)
	})
}

func (l *DefaultListener) ListenCloseTab(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		e, ok := ev.(*target.EventTargetDestroyed)
		if !ok || e.TargetID == "" {
			return
		}
		go l.tm.CloseTab(e.TargetID)
	})
}

func (l *DefaultListener) ListenAlert(ctx context.Context) {
	l.alert.AddListenAlert(ctx)
}
