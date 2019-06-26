package o365

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/evilsocket/islazy/tui"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/log"
)

func (z *O365) searchO365(keyword string) chromedp.Tasks {

	log.Info("Searching for '%s'...", keyword)

	url := fmt.Sprint("%s/%s", z.O365Search, keyword)

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

func (z *O365) dumpEmailByKeyword(keyword string) (html string) {

	log.Info("Searching in office365 for %s", tui.Bold(tui.Red(keyword)))

	a := action.Action{Target: z.Target}
	url := fmt.Sprintf("%s%s", z.O365Search, keyword)

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
