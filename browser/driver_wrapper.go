// browser/driver_wrapper.go
package browser

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/developer-commit/stealthchromedp/browser/alert"
	"github.com/developer-commit/stealthchromedp/browser/driver"
	"github.com/developer-commit/stealthchromedp/browser/listener"
	"github.com/developer-commit/stealthchromedp/browser/tabmanager"
	"github.com/developer-commit/stealthchromedp/runner"
)

// DriverWrapper는 CDP 드라이버, 탭매니저, 러너, AlertStream 등을 조립합니다.
type DriverWrapper struct {
	Driver      *driver.ChromeDpDriver
	TabManager  *tabmanager.TabManager
	TaskRunner  runner.TaskRunner
	AlertStream *alert.AlertStream
	Listener    listener.BrowserListener
}

func NewDriverWrapper(
	ctx context.Context,
	runner runner.TaskRunner,
	opts []chromedp.ExecAllocatorOption,
) (*DriverWrapper, error) {
	drv := driver.NewChromedpDriver(ctx, opts)
	tabCtx, tabCancel := drv.BuildTab()
	tid := chromedp.FromContext(tabCtx).Target.TargetID

	// 초기 stealth 스크립트 실행
	if err := runner.Run(tabCtx); err != nil {
		tabCancel()
		return nil, err
	}

	as := alert.NewAlertStream()
	as.AddListenAlert(tabCtx)

	tm := tabmanager.NewTabManager(tid, tabCtx, tabCancel)

	wrapper := &DriverWrapper{
		Driver:      drv,
		TabManager:  tm,
		TaskRunner:  runner,
		AlertStream: as,
	}

	wrapper.Listener = listener.NewDefaultListener(drv, tm, runner, as)
	wrapper.Listener.ListenNewTab(tabCtx)
	wrapper.Listener.ListenCloseTab(tabCtx)

	return wrapper, nil
}

func (w *DriverWrapper) Close() {
	w.Driver.Close()
}
