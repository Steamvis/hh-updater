package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

const (
	loginPageUrl  = "https://hh.ru/account/login"
	resumePageUrl = "https://hh.ru/applicant/resumes"
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
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var response string
	var screenshot []byte
	ticker := time.NewTicker(time.Minute * 241)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		log.Println("[INFO] tick")
		err := chromedp.Run(ctx,
			chromedp.Emulate(device.IPhone13ProMax),
			tasks(&response, &screenshot),
		)
		if err != nil {
			log.Fatalln(err)
		}

		if err := ioutil.WriteFile("fullScreenshot.png", screenshot, 0o644); err != nil {
			log.Fatal(err)
		}
	}
}

func tasks(res *string, screenshotRes *[]byte) chromedp.Tasks {
	inputEmail := `//input[@name="login"]`
	inputPassword := `//input[@data-qa="login-input-password"]`
	loginWithPasswordButton := `//button[@data-qa="expand-login-by-password"]`
	return chromedp.Tasks{
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
		chromedp.Click(
			`button[data-qa="account-login-submit"]`,
			chromedp.ByQuery,
		),

		// Wait Load New Page
		chromedp.Sleep(2 * time.Second),

		// Open Resume Page
		chromedp.Navigate(resumePageUrl),

		// Up Resume
		chromedp.Click(`//button[@data-qa="resume-update-button_actions"]`),

		chromedp.FullScreenshot(screenshotRes, 90),
	}

}
