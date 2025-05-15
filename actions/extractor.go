package actions

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Extractor interface {
	GetValue(ctx context.Context, sel string, opts ...chromedp.QueryOption) (string, error)
	GetText(ctx context.Context, sel string, opts ...chromedp.QueryOption) (string, error)
	GetAttr(ctx context.Context, sel, attr string, opts ...chromedp.QueryOption) (string, error)
}

type ImplExtractor struct{}

// GetText retrieves text content
func (d *ImplExtractor) GetText(ctx context.Context, sel string, opts ...chromedp.QueryOption) (string, error) {
	var text string
	if err := chromedp.Run(ctx, chromedp.Text(sel, &text, opts...)); err != nil {
		return "", err
	}
	return text, nil
}

// GetAttr retrieves attribute value
func (d *ImplExtractor) GetAttr(ctx context.Context, sel, attr string, opts ...chromedp.QueryOption) (string, error) {
	var val string
	if err := chromedp.Run(ctx, chromedp.AttributeValue(sel, attr, &val, nil, opts...)); err != nil {
		return "", err
	}
	return val, nil
}

// GetValue retrieves the "value" attribute of an input or textarea element.
func (d *ImplExtractor) GetValue(ctx context.Context, sel string, opts ...chromedp.QueryOption) (string, error) {
	var value string
	// Use the built-in Value action if available
	if err := chromedp.Run(ctx, chromedp.Value(sel, &value, opts...)); err != nil {
		return "", err
	}
	return value, nil
}
