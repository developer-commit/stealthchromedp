package runner

import (
	"context"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	pageproto "github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type LowDataSteathRunner struct {
	script       string
	locale       string
	timezone     string
	customDevice chromedp.Device
}

func NewLowDataSteathRunner(script, locale, timezone string, my_device chromedp.Device) LowDataSteathRunner {
	return LowDataSteathRunner{
		script:       script,
		locale:       locale,
		timezone:     timezone,
		customDevice: my_device,
	}
}

func (r *LowDataSteathRunner) Run(ctx context.Context) error {
	task := chromedp.Tasks{
		// ── 1) 도메인 활성화 ──────────────────────────
		network.Enable(), // must enable before setting headers
		network.SetExtraHTTPHeaders(network.Headers{
			"Accept-Language": r.locale,
		}),
		pageproto.Enable(),
		// ── 저데이터 모드 활성화
		network.SetBlockedURLs([]string{
			"*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.svg", "*.woff", "*.woff2", "*.ttf", "*.eot", "*.mp4", "*.webm",
		}),
		// ── 2) CSP 우회 ──────────────────────────────
		pageproto.SetBypassCSP(true),

		// ── 3) 헤더·환경 설정 ────────────────────────
		emulation.SetTimezoneOverride(r.timezone),
		emulation.SetLocaleOverride().WithLocale(r.locale),

		chromedp.Emulate(r.customDevice),
		// ❶ 먼저 스크립트 훅을 걸고
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := pageproto.
				AddScriptToEvaluateOnNewDocument(r.script).
				WithRunImmediately(true).
				Do(ctx)
			return err
		}),
		// ❷ 그다음 리로드해서 주입된 상태로 로드하게 만들기
		// chromedp.Reload(),
	}

	return chromedp.Run(
		ctx,
		task,
	)
}
