package screenshot

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

var (
	debug bool = os.Getenv("DEBUG") != ""
)

func Make(ctx context.Context, screenshotName string) {
	if debug {
		var screenshot []byte
		chromedp.Run(ctx, chromedp.FullScreenshot(&screenshot, 90))
		if err := ioutil.WriteFile(screenshotName+".png", screenshot, 0644); err != nil {
			log.Fatalln(err)
		}
		log.Println("[INFO] Make screenshot", screenshotName)
	}
}
