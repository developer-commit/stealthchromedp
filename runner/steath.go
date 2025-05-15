package runner

import (
	"context"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	pageproto "github.com/chromedp/cdproto/page"
)

type SteathRunner struct {
	script       string
	locale       string
	timezone     string
	customDevice chromedp.Device
}

func NewSteathRunner(script, locale, timezone string, my_device chromedp.Device) SteathRunner {
	return SteathRunner{
		script:       script,
		locale:       locale,
		timezone:     timezone,
		customDevice: my_device,
	}
}

func (r *SteathRunner) Run(ctx context.Context) error {
	task := chromedp.Tasks{
		// ── 1) 도메인 활성화 ──────────────────────────
		network.Enable(), // must enable before setting headers
		network.SetExtraHTTPHeaders(network.Headers{
			"Accept-Language": r.locale,
		}),
		pageproto.Enable(),

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
