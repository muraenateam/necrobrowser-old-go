package navigation

import (
	"time"

	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "Navigation"

	// Description of this action
	Description = "Navigation allows to surf web pages"
)

// Navigation is an action
type Navigation struct {
	action.Action

	URL string
}

// Name returns the action name
func (a *Navigation) Name() string {
	return Name
}

// Description returns what the action does
func (a *Navigation) Description() string {
	return Description
}

// Navigate changes the browser URL
func (a *Navigation) Navigate() (err error) {
	z := a.Target
	z.Info("Navigating to %s", a.URL)

	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		chromedp.Sleep(2 * time.Second),
	}

	return a.Run(t)
}
