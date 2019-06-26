package login

import (
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "Login"

	// Description of this action
	Description = "Login"
)

// Navigation is an action
type Login struct {
	action.Action

	URL string

	Username         string
	UsernameSelector string
	Password         string
	PasswordSelector string
}

// Name returns the action name
func (a *Login) Name() string {
	return Name
}

// Description returns what the action does
func (a *Login) Description() string {
	return Description
}

// Navigate changes the browser URL
func (a *Login) Do() (err error) {
	z := a.Target
	z.Info("Logging into %s", a.URL)
	z.Info("Username form selector [%s] - Passwd form selector [%s]", a.UsernameSelector, a.PasswordSelector)

	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		//chromedp.Sleep(2 * time.Second),

		chromedp.WaitVisible(a.UsernameSelector, chromedp.ByQuery),
		chromedp.SendKeys(a.UsernameSelector, a.Username+"\n"),

		chromedp.WaitVisible(a.PasswordSelector, chromedp.ByQuery),
		chromedp.SendKeys(a.PasswordSelector, a.Password+"\n"),
	}

	// TODO wait for some visibile elemtn to see if we are logged in, return error otherwise

	return a.Run(t)
}
