package proxy_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/developer-commit/stealthchromedp/browser"
	"github.com/developer-commit/stealthchromedp/proxy"
	"github.com/developer-commit/stealthchromedp/runner"
)

var code string = `
(() => {
  const { defineProperty } = Object;

  /* ── 마우스 디버거 ───────────────────────────── */
  const DOT_SIZE = 6;

  onReady(() => {
    const dot = document.createElement('div');
    Object.assign(dot.style, {
      position: 'fixed',
      width: "DOT_SIZE" + px,
      height: "DOT_SIZE" + px,
      borderRadius: '50%',
      background: 'red',
      pointerEvents: 'none',
      zIndex: 999999,
      transition: 'top .03s linear, left .03s linear',
    });

    const box = document.createElement('div');
    Object.assign(box.style, {
      position: 'fixed',
      top: '8px',
      left: '8px',
      padding: '2px 6px',
      fontSize: '11px',
      fontFamily: 'monospace',
      background: 'rgba(0,0,0,.7)',
      color: 'lime',
      pointerEvents: 'none',
      zIndex: 999999,
    });

    const keyBox = document.createElement('div');
    Object.assign(keyBox.style, {
      position: 'fixed',
      top: '30px',
      left: '8px',
      padding: '2px 6px',
      fontSize: '11px',
      fontFamily: 'monospace',
      background: 'rgba(0,0,0,.7)',
      color: 'orange',
      pointerEvents: 'none',
      zIndex: 999999,
    });
    document.body.appendChild(keyBox);

    document.body.append(dot, box);

    window.__updateMouseDebug = function(x, y) {
      dot.style.left  = (x - DOT_SIZE / 2) + "px";
      dot.style.top   = (y - DOT_SIZE / 2) + "px";
      box.textContent = "x:" + Math.round(x) + " y:" + Math.round(y);
    };

    document.addEventListener('mousemove', e =>
      window.__updateMouseDebug(e.clientX, e.clientY));

    document.addEventListener('keydown', e => {
      keyBox.textContent = "key: " + e.key;
    });
  });
})();
`

func TestProxy(t *testing.T) {
	relay, err := proxy.NewAuthProxyRelay(
		":8081",
		"http://brd-customer-hl_67f64aad-zone-mobile_proxy1:5w9mpdztf6cx@brd.superproxy.io:33335", // 사내 인증 프록시
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := relay.Run(ctx); err != nil && err != context.Canceled {
			log.Fatal(err)
		}
	}()

	run := runner.NewLowDataSteathRunner(
		code,
		"pt-PT",
		"Europe/Lisbon",
		device.IPhone13Pro,
	)
	opts := browser.ProxyBuildOption("http://127.0.0.1:8081")

	dw, err := browser.NewDriverWrapper(ctx, &run, opts)
	if err != nil {
		t.Log(err)
		return
	}

	if now, ok := dw.TabManager.Now(); ok {
		chromedp.Run(
			now.Ctx,
			chromedp.Navigate("https://nid.naver.com/nidlogin.login"),
		)
	}

	time.Sleep(120 * time.Second)

}
