package browser_test

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

//go test -v -count=1 ./browser/wrapper_test.go

func TestWrapper(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run := runner.LowDataSteathRunner{CustomDevice: device.IPhone14}
	// opts := browser.ProxyBuildOption("http://brd-customer-hl_67f64aad-zone-mobile_proxy1:5w9mpdztf6cx@brd.superproxy.io:33335")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("http://brd-customer-hl_67f64aad-zone-mobile_proxy1-country-us:5w9mpdztf6cx@brd.superproxy.io:33335"),
		chromedp.Flag("headless", false),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)

	defer allocCancel()
	dctx, _ := chromedp.NewContext(allocCtx)
	if err := chromedp.Run(dctx,
		chromedp.Navigate("http://example.com"),
	); err != nil {
		t.Log(err) // 인증 실패 시 에러 발생
	}
	// dw, err := browser.NewDriverWrapper(ctx, &run, opts)
	// if err != nil {
	// 	t.Log(err)
	// 	return
	// }

	// if now, ok := dw.TabManager.Now(); ok {
	// 	chromedp.Run(
	// 		now.Ctx,
	// 		chromedp.Navigate("https://www.naver.com/"),
	// 	)
	// }
	time.Sleep(30 * time.Second)

	// dw.Driver.Close()
}
