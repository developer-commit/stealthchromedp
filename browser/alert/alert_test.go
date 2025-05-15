package alert_test

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/developer-commit/stealthchromedp/browser/alert"
)

// go test -v -count=1 ./browser/alert/
func TestAlert(t *testing.T) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.UserDataDir(fmt.Sprintf("data/profile/%s", profileName)),
		chromedp.Flag("guest", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("useAutomationExtension", false),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	defer allocCancel()
	ctx, _ := chromedp.NewContext(allocCtx)

	as := alert.NewAlertStream()
	as.AddListenAlert(ctx)

	as.PushHandler(func(ctx context.Context, dialog *page.EventJavascriptDialogOpening) {
		t.Log("alert handle")
		chromedp.Run(ctx, page.HandleJavaScriptDialog(true))
	})

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://www.selenium.dev/selenium/web/alerts.html#"),
		chromedp.DoubleClick(`//*[@id="alert"]`),
	)

	if err != nil {
		t.Log(err)
	}

	time.Sleep(5 * time.Second)

}
