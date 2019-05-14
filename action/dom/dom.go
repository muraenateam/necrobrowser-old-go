package dom

import (
	"fmt"
	"time"

	"github.com/muraenateam/necrobrowser/action"

	"io/ioutil"

	"github.com/chromedp/chromedp"
)

const (
	// Name of this action
	Name = "DOM"

	// Description of this action
	Description = "DOM based actions"
)

// Screenshot is an action
type DOM struct {
	action.Action

	Selector string
	URL      string
}

// Name returns the action name
func (a *DOM) Name() string {
	return Name
}

// Description returns what the action does
func (a *DOM) Description() string {
	return Description
}

func (a *DOM) DumpSelector(html *string) (err error) {
	z := a.Target
	z.Info("Dumping raw html contents of %s (%s)", a.URL, a.Selector)
	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		chromedp.Sleep(5 * time.Second),
		chromedp.OuterHTML(a.Selector, html, chromedp.ByQueryAll),
	}

	if err = a.Run(t); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/dump_%s.html", z.LootPath, action.Now())
	z.Info("Dumping to %s", path)
	err = ioutil.WriteFile(path, []byte(*html), 0644)
	return
}

func (a *DOM) InjectJS(inject string, extrudedEmails *[]string) (err error) {

	z := a.Target
	z.Info("Injecting JS contents into %s", a.URL)

	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		chromedp.Evaluate(inject, extrudedEmails),
	}

	return a.Run(t)

}
