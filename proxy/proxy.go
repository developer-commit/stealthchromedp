package proxy

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/elazarl/goproxy"
)

// ────── 인터페이스 정의 ──────
type ProxyRelay interface {
	// ctx가 Done() 되면 서버를 자동으로 종료
	Run(ctx context.Context) error
	// 외부에서 즉시 종료하고 싶을 때 (예: SIGINT 핸들러)
	Shutdown(ctx context.Context) error
}

// ────── 구조체 구현 ──────
type AuthProxyRelay struct {
	ListenAddr string
	Upstream   *url.URL

	server *http.Server
	errCh  chan error
}

// ────── 생성자 ──────
func NewAuthProxyRelay(listen, upstream string) (*AuthProxyRelay, error) {
	u, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}
	return &AuthProxyRelay{
		ListenAddr: listen,
		Upstream:   u,
		errCh:      make(chan error, 1),
	}, nil
}

// ────── Run: 컨텍스트 기반 메인 루프 ──────
func (p *AuthProxyRelay) Run(ctx context.Context) error {
	// goproxy 설정
	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true

	authHeader := "Basic " + basicAuth(p.Upstream)

	// 1) 평문 HTTP 트랜스포트 + CONNECT 헤더
	proxy.Tr = &http.Transport{
		Proxy: http.ProxyURL(p.Upstream),
		ProxyConnectHeader: http.Header{
			"Proxy-Authorization": []string{authHeader},
		},
	}

	// 2) CONNECT 터널에도 인증 강제 주입
	proxy.ConnectDial = proxy.NewConnectDialToProxyWithHandler(
		p.Upstream.String(),
		func(r *http.Request) { r.Header.Set("Proxy-Authorization", authHeader) },
	)

	// 3) 모든 요청 헤더 클린업 + 인증 주입
	proxy.OnRequest().DoFunc(func(r *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		r.Header.Set("Proxy-Authorization", authHeader)
		r.Header.Del("X-Forwarded-For")
		r.Header.Del("Via")
		r.Header.Del("Forwarded")
		return r, nil
	})

	// HTTP 서버 생성
	p.server = &http.Server{
		Addr:              p.ListenAddr,
		Handler:           proxy,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 별도 고루틴에서 ListenAndServe
	go func() {
		log.Printf("AuthProxyRelay listening at %s\n", p.ListenAddr)
		if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
			p.errCh <- err
		}
		close(p.errCh)
	}()

	// 컨텍스트 또는 서버 오류 대기
	select {
	case <-ctx.Done():
		// 외부에서 취소 신호
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.Shutdown(shutdownCtx)
	case err := <-p.errCh:
		// ListenAndServe 오류
		return err
	}
}

// ────── Shutdown: 그레이스풀 셧다운 ──────
func (p *AuthProxyRelay) Shutdown(ctx context.Context) error {
	if p.server != nil {
		return p.server.Shutdown(ctx)
	}
	return nil
}

// ────── 헬퍼: 기본 인증 문자열 생성 ──────
func basicAuth(u *url.URL) string {
	if u == nil || u.User == nil {
		return ""
	}
	pass, _ := u.User.Password()
	creds := u.User.Username() + ":" + pass
	return base64.StdEncoding.EncodeToString([]byte(creds))
}
