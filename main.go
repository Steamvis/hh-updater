package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

const (
	loginPageUrl            string = "https://hh.ru/account/login"
	resumePageUrl           string = "https://hh.ru/applicant/resumes"
	loginButton             string = `//button[@data-qa="account-login-submit"]`
	inputEmail              string = `//input[@name="login"]`
	inputPassword           string = `//input[@data-qa="login-input-password"]`
	loginWithPasswordButton string = `//button[@data-qa="expand-login-by-password"]`
	upButtonWithText        string = `//button[@data-qa="resume-update-button_actions" and text()="Поднять в поиске"]`
)

var (
	email    string = os.Getenv("EMAIL")
	password string = os.Getenv("PASSWORD")
)

func init() {
	if len(email) == 0 || len(password) == 0 {
		fmt.Println("Length of password or email is 0 symbols")
		os.Exit(1)
	}
}

func main() {
	var clickUpResumeNodes []*cdp.Node
	ticker := time.NewTicker(time.Minute * 241)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		log.Println("[INFO] tick")
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			// chromedp.WithDebugf(log.Printf),
		)

		err := chromedp.Run(ctx,
			chromedp.Emulate(device.IPhone13ProMax),
			// Open Login Page
			chromedp.Navigate(loginPageUrl),
			chromedp.WaitVisible(inputEmail),

			// Fill Form
			chromedp.SendKeys(inputEmail, email),
			// Click Button Login With Password
			chromedp.Click(loginWithPasswordButton),

			chromedp.WaitVisible(inputPassword),
			chromedp.SendKeys(inputPassword, password),

			// Click Login Button
			chromedp.Click(loginButton),

			// Wait Load Page
			chromedp.Sleep(5*time.Second),

			// Open Resume Page
			chromedp.Navigate(resumePageUrl),
			// Search Up Buttons
			chromedp.Nodes(upButtonWithText, &clickUpResumeNodes, chromedp.AtLeast(0)),
		)

		if err != nil {
			log.Fatalln(err)
		}

		for _, cdp := range clickUpResumeNodes {
			err := chromedp.Run(ctx, chromedp.MouseClickNode(cdp))
			if err != nil {
				log.Fatalln(err)
			}
		}

		cancel()
	}
}
