package actions_test

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/developer-commit/stealthchromedp/actions"
)

//go test -v -count=1 ./actions

func TestGesture(t *testing.T) {
	mCtx, mCancel := context.WithCancel(context.Background())
	defer mCancel()

	allocCtx, allocCancel := chromedp.NewExecAllocator(mCtx)
	defer allocCancel()
	tapCtx, _ := chromedp.NewContext(allocCtx)

	implGesture := actions.ImplGesture{Noise: 3.0}
	implNavigator := actions.ImplNavigator{}
	chromedp.Run(tapCtx, chromedp.Emulate(device.IPhone12Mini))

	func(gesture actions.Gesture, navi actions.Navigator) {
		navi.GoTo(tapCtx, "https://www.selenium.dev/selenium/web/web-form.html")
		gesture.Tap(tapCtx, "/html/body/main/div/form/div/div[2]/label[1]/select", 0, 0)
		gesture.Select(tapCtx, "/html/body/main/div/form/div/div[2]/label[1]/select", "2")
		time.Sleep(1 * time.Second)
		gesture.Tap(tapCtx, "/html/body/main/div/form/div/div[2]/button", 0, 0)
		time.Sleep(1 * time.Second)

	}(&implGesture, &implNavigator)
}

func TestExtractor(t *testing.T) {
	mCtx, mCancel := context.WithCancel(context.Background())
	defer mCancel()

	allocCtx, allocCancel := chromedp.NewExecAllocator(mCtx)
	defer allocCancel()
	exCtx, _ := chromedp.NewContext(allocCtx)

	implExtractor := actions.ImplExtractor{}
	implNavigator := actions.ImplNavigator{}
	chromedp.Run(exCtx, chromedp.Emulate(device.IPhone12Mini))

	func(ext actions.Extractor, navi actions.Navigator) {
		navi.GoTo(exCtx, "https://www.selenium.dev/selenium/web/web-form.html")

		if text, err := ext.GetText(exCtx, `/html/body/main/div/div/div/h1`); err == nil {
			t.Log(text)
		}
		time.Sleep(1 * time.Second)

		if link, err := ext.GetAttr(exCtx, `/html/body/main/div/form/div/div[1]/div/a`, "href"); err == nil {
			t.Log(link)
		}
		time.Sleep(1 * time.Second)

		if v, err := ext.GetAttr(exCtx, `/html/body/main/div/form/div/div[3]/label[1]/input`, "value"); err == nil {
			t.Log(v)
		}
		time.Sleep(1 * time.Second)

		if v, err := ext.GetValue(exCtx, `/html/body/main/div/form/div/div[3]/label[1]/input`); err == nil {
			t.Log(v)
		}
		time.Sleep(1 * time.Second)

	}(&implExtractor, &implNavigator)
}

func TestTyper(t *testing.T) {
	mCtx, mCancel := context.WithCancel(context.Background())
	defer mCancel()

	allocCtx, allocCancel := chromedp.NewExecAllocator(mCtx)
	defer allocCancel()
	typeCtx, _ := chromedp.NewContext(allocCtx)

	implTyper := actions.ImplTyper{
		IGesture:   &actions.ImplGesture{},
		IExtractor: &actions.ImplExtractor{},
	}

	implNavigator := actions.ImplNavigator{}
	chromedp.Run(typeCtx, chromedp.Emulate(device.IPhone12Mini))

	func(typer actions.Typer, navi actions.Navigator) {
		navi.GoTo(typeCtx, "https://www.selenium.dev/selenium/web/web-form.html")

		if err := typer.InputText(typeCtx, `//*[@id="my-text-id"]`, "thisismyid"); err != nil {
			t.Log(err)
		}

		if err := typer.InputText(typeCtx, `/html/body/main/div/form/div/div[1]/label[2]/input`, "testsdfaklj"); err != nil {
			t.Log(err)
		}

		if err := typer.Backspace(typeCtx, `/html/body/main/div/form/div/div[1]/label[2]/input`); err != nil {
			t.Log(err)
		}

	}(&implTyper, &implNavigator)
}
