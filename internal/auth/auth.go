package auth

import (
	"context"
	"hhupdater/internal/screenshot"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

const (
	loginPageUrl            string = "https://hh.ru/account/login"
	loginButton             string = `//button[@data-qa="account-login-submit"]`
	inputEmail              string = `//input[@name="login"]`
	inputPassword           string = `//input[@data-qa="login-input-password"]`
	loginWithPasswordButton string = `//button[@data-qa="expand-login-by-password"]`
	logoutButton            string = `//button[@data-qa="mainmenu_applicantProfile"]`
)

func Login(ctx context.Context, email string, password string) error {
	return chromedp.Run(
		ctx,
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
	)
}

func IsAuthorized(ctx context.Context) bool {
	log.Println("[INFO] Check authorization")
	l := checkLogin(ctx)
	if l {
		log.Println("[INFO] Check authorization: already logged in")
		return true
	}
	log.Println("[INFO] Check authorization: unauthorized")
	return false
}

func checkLogin(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	var nodes []*cdp.Node
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(loginPageUrl),
		chromedp.Sleep(time.Second),
		chromedp.Nodes(logoutButton, &nodes, chromedp.AtLeast(0)),
	)

	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	screenshot.Make(ctx, "checkLogin")

	return len(nodes) >= 1
}
