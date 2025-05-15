package driver

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type ChromeDpDriver struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc
}

func NewChromedpDriver(ctx context.Context, opts []chromedp.ExecAllocatorOption) *ChromeDpDriver {
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	return &ChromeDpDriver{
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
	}
}

func (d *ChromeDpDriver) WaitForDriverClose() {
	<-d.allocCtx.Done()
}

func (d *ChromeDpDriver) BuildTab() (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(d.allocCtx)

	// 반드시 최소 한 번 실행해야 target이 설정됨
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		panic(fmt.Errorf("failed to initialize tab context: %w", err))
	}

	return ctx, cancel
}

func (d *ChromeDpDriver) BuildTabByID(ctx context.Context, tid target.ID) (context.Context, context.CancelFunc) {
	return chromedp.NewContext(ctx, chromedp.WithTargetID(tid))
}

func (d *ChromeDpDriver) Close() {
	d.allocCancel()
}
