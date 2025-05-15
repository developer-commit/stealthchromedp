package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type StorageState struct {
	Cookies []*network.Cookie `json:"cookies"`
}

// SaveSession persists cookies and localStorage
func SaveSession(ctx context.Context, path string) error {
	var cookies []*network.Cookie
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	})); err != nil {
		return err
	}

	state := StorageState{Cookies: cookies}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadSession restores cookies and localStorage
func LoadSession(ctx context.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var state StorageState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	// Set cookies
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		for _, c := range state.Cookies {
			timeEpoch := cdp.TimeSinceEpoch(time.Unix(int64(c.Expires), 0))
			if err := network.SetCookie(c.Name, c.Value).
				WithDomain(c.Domain).
				WithPath(c.Path).
				WithExpires(&timeEpoch).
				WithHTTPOnly(c.HTTPOnly).
				WithSecure(c.Secure).
				Do(ctx); err != nil {
				return fmt.Errorf("cookie set failed: %v", err)
			}
		}
		return nil
	})); err != nil {
		return err
	}
	return nil
}
