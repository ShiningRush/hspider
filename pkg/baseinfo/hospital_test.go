package baseinfo

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"testing"
	"time"
)

func TestGetHospitals(t *testing.T) {
	// create context
	opts := append([]chromedp.ExecAllocatorOption(nil), chromedp.ExecPath(`C:\Users\vincixu.TENCENT\AppData\Local\Google\Chrome\Application\chrome.exe`))
	//opts = append(opts, chromedp.Headless)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		opts...,)
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// run task list
	for {
		var res string
		err := chromedp.Run(ctx,
			chromedp.Navigate(`https://sz.91160.com/`),
			chromedp.WaitReady(`ul.combo-dropdown li[data-index="1"]`, chromedp.ByQuery),
			chromedp.InnerHTML(`ul.combo-dropdown`, &res,chromedp.ByQuery),
			chromedp.Sleep(time.Second * 3),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(strings.TrimSpace(res))
	}
}
