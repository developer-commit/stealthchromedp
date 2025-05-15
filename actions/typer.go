package actions

import (
	"context"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

type Typer interface {
	InputText(ctx context.Context, sel, text string, opts ...chromedp.QueryOption) error
	Backspace(ctx context.Context, sel string, opts ...chromedp.QueryOption) error
}

type ImplTyper struct {
	IGesture   Gesture
	IExtractor Extractor
}

func (i *ImplTyper) InputText(ctx context.Context, sel, text string, opts ...chromedp.QueryOption) error {
	// 1) 먼저 클릭하여 입력 필드 포커스
	if err := i.IGesture.Tap(ctx, sel, 0, 0, opts...); err != nil {
		return err
	}
	// 2) 인간처럼 오타와 딜레이를 포함해 한 글자씩 입력
	for _, ch := range text {
		// 숫자와 알파벳 풀 결정
		pool := "abcdefghijklmnopqrstuvwxyz"
		if ch >= '0' && ch <= '9' {
			pool = "0123456789"
		}
		// 오타 확률 적용
		if rand.Float64() < 0.1 && strings.ContainsRune(pool, unicode.ToLower(ch)) {
			// 오타 입력
			typoRune := rune(pool[rand.Intn(len(pool))])
			_ = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				return chromedp.SendKeys(sel, string(typoRune)).Do(ctx)
				//return input.DispatchKeyEvent(input.KeyChar).WithText(string(typoRune)).Do(ctx)
			}))
			time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
			// 백스페이스
			_ = chromedp.Run(ctx, chromedp.SendKeys(sel, kb.Backspace, opts...))
			time.Sleep(time.Duration(rand.Intn(150)+50) * time.Millisecond)
		}
		// 실제 문자 입력
		_ = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.SendKeys(sel, string(ch)).Do(ctx)
			//return input.DispatchKeyEvent(input.KeyChar).WithText(string(ch)).Do(ctx)
		}))
		// 글자 간 기본 딜레이
		time.Sleep(time.Duration(rand.Intn(200)+350) * time.Millisecond)
	}
	return nil
}

// ActionBackspace clears the input field by sending backspaces equal to the current value length.
func (i *ImplTyper) Backspace(ctx context.Context, sel string, opts ...chromedp.QueryOption) error {
	// Focus the field
	if err := i.IGesture.Tap(ctx, sel, 0, 0, opts...); err != nil {
		return err
	}

	// Retrieve current text value and its length
	text, err := i.IExtractor.GetValue(ctx, sel, opts...)
	if err != nil {
		return err
	}
	count := len(text)

	// Build actions: send Backspace for each character with a small delay
	actions := make([]chromedp.Action, 0, count*2)
	for i := 0; i < count; i++ {
		actions = append(actions,
			chromedp.SendKeys(sel, kb.Backspace, opts...),
			chromedp.Sleep(time.Duration(rand.Intn(150)+200)*time.Millisecond),
		)
	}

	// Execute all actions in a single Run call
	return chromedp.Run(ctx, actions...)
}
