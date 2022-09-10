package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"

	"hhupdater/internal/auth"
	"hhupdater/internal/resume"
	"hhupdater/internal/screenshot"
)

const (
	scheduleTimeInMinutes time.Duration = 60
)

var (
	email     string = os.Getenv("EMAIL")
	password  string = os.Getenv("PASSWORD")
	debugMode string = os.Getenv("DEBUG")
)

func init() {
	if len(email) == 0 || len(password) == 0 {
		fmt.Println("Length of password or email is 0 symbols")
		os.Exit(1)
	}
}

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ticker := time.NewTicker(time.Minute * scheduleTimeInMinutes)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		log.Println("[INFO] Tick")

		if isAuthorized := auth.IsAuthorized(ctx); !isAuthorized {
			if ctx != nil {
				cancel()
			}

			ctx, cancel = newCtx()
			err := auth.Login(ctx, email, password)
			if err != nil {
				screenshot.Make(ctx, "errorLogin")
				log.Fatal("[ERROR] ", err)
			}
			log.Println("[INFO] SuccessLogin")
			screenshot.Make(ctx, "successLogin")
		}

		resume.Up(ctx)
		screenshot.Make(ctx, "upResume")
	}
}

func newCtx() (context.Context, context.CancelFunc) {
	if debugMode == "chromedp" {
		return chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	} else {
		return chromedp.NewContext(context.Background())
	}
}
