package browser

import (
	"github.com/chromedp/chromedp"
)

func PureBuildOption() []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("guest", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("useAutomationExtension", false),
	)
	return opts
}

func ProxyBuildOption(proxyURL string) []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("guest", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("useAutomationExtension", false),
		chromedp.Flag("ignore-certificate-errors", true),

		chromedp.ProxyServer(proxyURL),
	)
	return opts
}
