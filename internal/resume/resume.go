package resume

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

const (
	resumePageUrl string = "https://hh.ru/applicant/resumes"
	upButton      string = `//button[@data-qa="resume-update-button_actions"]`
)

func Up(ctx context.Context) {
	err := chromedp.Run(ctx,
		// Open Resume Page
		chromedp.Navigate(resumePageUrl),

		// Search Up Button
		chromedp.WaitVisible(upButton),
		chromedp.Click(upButton),
	)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("[INFO] ResumeUp")
}
