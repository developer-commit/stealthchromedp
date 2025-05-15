package listener

import "context"

type BrowserListener interface {
	ListenNewTab(ctx context.Context)   // 새 탭 감지
	ListenCloseTab(ctx context.Context) // 탭 닫힘 감지
	ListenAlert(ctx context.Context)    // JS alert 감지
}
