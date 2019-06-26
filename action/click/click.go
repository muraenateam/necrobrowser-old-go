package click

import (
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "Click"

	// Description of this action
	Description = "Click"
)

// Navigation is an action
type Click struct {
	action.Action
	Selector string
}

// Name returns the action name
func (a *Click) Name() string {
	return Name
}

// Description returns what the action does
func (a *Click) Description() string {
	return Description
}

// Click on an element
func (click *Click) Click() (err error) {
	z := click.Target
	z.Info("Clicking element %s", click.Selector)

	t := chromedp.Tasks{

		chromedp.WaitVisible(click.Selector, chromedp.ByQuery),
		chromedp.Click(click.Selector, chromedp.NodeVisible),
	}

	return click.Run(t)
}
