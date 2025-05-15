package actions

import (
	"context"
	"math/rand"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/developer-commit/stealthchromedp/util"
)

type Gesture interface {
	Tap(ctx context.Context, sel string, ux, uy float64, opts ...chromedp.QueryOption) error
	Select(ctx context.Context, sel, val string, opts ...chromedp.QueryOption) error
	// Swipe(ctx context.Context, up bool) error
}

type ImplGesture struct {
	Noise float64
}

func (i *ImplGesture) Tap(ctx context.Context, sel string, ux, uy float64, opts ...chromedp.QueryOption) error {
	// 1) d.ctx 를 기반으로 10초짜리 서브 컨텍스트 생성
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var boxModel *dom.BoxModel
	var nodes []*cdp.Node

	// XPath를 사용하여 요소의 박스 모델 정보 가져오기
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(sel, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 노드 ID 가져오기

			err := chromedp.Nodes(sel, &nodes, chromedp.BySearch).Do(ctx)
			if err != nil {
				return err
			}
			// 박스 모델 정보 가져오기
			boxModel, err = dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err
	}

	// 중심 좌표 계산
	content := boxModel.Content
	x := (content[0]+content[4])/2 + ux
	y := (content[1]+content[5])/2 + uy

	util.NoiseGenerator(x, i.Noise)
	util.NoiseGenerator(y, i.Noise)

	// 해당 좌표 클릭
	return chromedp.Run(ctx,
		input.DispatchMouseEvent(input.MouseMoved, x, y),
		// 그 다음 마우스 다운
		input.DispatchMouseEvent(input.MousePressed, x, y).
			WithButton("left").
			WithClickCount(1),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond) // 50~150ms
			return nil
		}),
		// 그리고 마우스 업
		input.DispatchMouseEvent(input.MouseReleased, x, y).
			WithButton("left").
			WithClickCount(1),
	)
}

func (d *ImplGesture) Select(ctx context.Context, sel, val string, opts ...chromedp.QueryOption) error {
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	return chromedp.Run(tctx,
		chromedp.SetValue(sel, val, opts...),
	)
}

// func (i *ImplGesture) Swipe(ctx context.Context, up bool) error {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	//layoutViewport, visualViewport, contentSize, cssLayoutViewport, cssVisualViewport, cssContentSize, err := page.GetLayoutMetrics().Do(ctx)
// 	layoutViewport, _, _, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	width := layoutViewport.ClientWidth
// 	height := layoutViewport.ClientHeight

// 	startX := float64(width) / 2
// 	startY := float64(height) / 2
// 	endY := startY + rand.Float64()*150 // 0~150px

// 	if up {
// 		endY = startY - rand.Float64()*150
// 	}

// 	// 랜덤 X 범위 설정 (약간만 흔들림)
// 	startX = util.NoiseGenerator(startX, i.Noise)
// 	endX := util.NoiseGenerator(startX, i.Noise/2)

// 	steps := 5

// 	touchTasks := chromedp.Tasks{}

// 	touchpoint := input.TouchPoint{
// 		X: startX,
// 		Y: startY,
// 	}
// 	touchTasks = append(touchTasks, input.DispatchTouchEvent(input.TouchStart, []*input.TouchPoint{&touchpoint}))

// 	// TouchMove 단계별로 추가
// 	for i := 1; i <= steps; i++ {
// 		interX := startX + (endX-startX)*float64(i)/float64(steps)
// 		interY := startY + (endY-startY)*float64(i)/float64(steps)

// 		touchpoint := input.TouchPoint{
// 			X: interX,
// 			Y: interY,
// 		}
// 		touchTasks = append(touchTasks, input.DispatchTouchEvent(input.TouchMove, []*input.TouchPoint{&touchpoint}))

// 	}
// 	// TouchEnd
// 	touchpoint = input.TouchPoint{
// 		X: endX,
// 		Y: endY,
// 	}
// 	touchTasks = append(touchTasks, input.DispatchTouchEvent(input.TouchEnd, []*input.TouchPoint{&touchpoint}))

// 	// 실행
// 	return chromedp.Run(ctx, touchTasks)
// }
