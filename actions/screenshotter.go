package actions

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
)

type Screenshotter interface {
	FullPage(ctx context.Context, quality int) ([]byte, error)
}

type ImplScreenshotter struct{}

// Screenshot captures full page screenshot
func (d *ImplScreenshotter) Screenshot(ctx context.Context, path string) error {
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90)); err != nil {
		return err
	}
	return os.WriteFile(path, buf, 0644)
}
