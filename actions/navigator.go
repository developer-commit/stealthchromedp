package actions

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Navigator interface {
	GoTo(ctx context.Context, url string) error
	Back(ctx context.Context) error
	Forward(ctx context.Context) error
	Refresh(ctx context.Context) error
}

type ImplNavigator struct{}

func (d *ImplNavigator) GoTo(ctx context.Context, url string) error {
	return chromedp.Run(ctx, chromedp.Navigate(url))
}

func (d *ImplNavigator) Back(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.NavigateBack())
}

func (d *ImplNavigator) Forward(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.NavigateForward())
}

func (d *ImplNavigator) Refresh(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Reload())
}
