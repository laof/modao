package data

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func txt(urlstr string, str *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitNotPresent("script"),
		chromedp.InnerHTML("body", str),
	}
}

func Body(url string) string {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 32*time.Second)
	defer cancel()

	var str = ""

	chromedp.Run(ctx, txt(url, &str))

	return str

}
