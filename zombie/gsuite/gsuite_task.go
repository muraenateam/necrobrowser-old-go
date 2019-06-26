package gsuite

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/evilsocket/islazy/tui"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/log"
)

func (z *GSuite) searchGSuite(keyword string) chromedp.Tasks {

	log.Info("Searching for '%s'...", keyword)

	url := fmt.Sprint("%s/%s", z.GMailSearch, keyword)

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),

		// #aso_search_form_anchor > div > input
		chromedp.WaitVisible(`#table`, chromedp.ByQueryAll),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`#aso_search_form_anchor`),
		chromedp.SendKeys(`#aso_search_form_anchor > div > input`, keyword+"\n"),
		chromedp.Sleep(4 * time.Second),
	}
}

func (z *GSuite) dumpEmailByMessageId(messageId string) (html string) {

	log.Info("Searching for messageId %s", tui.Bold(tui.Red(messageId)))

	a := action.Action{Target: z.Target}

	// Dump Settings
	d := &dom.DOM{
		Action:   a,
		URL:      "https://mail.google.com/mail/u/0/#inbox",
		Selector: `div[role='main']`,
	}
	if err := d.DumpGsuiteEmailByMessageId(1, messageId, &html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		log.Debug(html)
	}

	return
}

func (z *GSuite) dumpEmailByKeyword(keyword string) (html string) {

	log.Info("Searching in gmail for %s", tui.Bold(tui.Red(keyword)))

	a := action.Action{Target: z.Target}
	url := fmt.Sprintf("%s%s", z.GMailSearch, keyword)

	// Dump Settings
	d := &dom.DOM{
		Action:   a,
		URL:      url,
		Selector: `html`,
	}
	if err := d.DumpSelector(&html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		log.Debug(html)
	}

	return
}
