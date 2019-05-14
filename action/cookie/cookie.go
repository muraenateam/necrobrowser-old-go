package cookie

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "Cookie"

	// Description of this action
	Description = "Cookie performs operation on cookies"
)

// Screenshot is an action
type Cookie struct {
	action.Action
}

// Name returns the action name
func (a *Cookie) Name() string {
	return Name
}

// Description returns what the action does
func (a *Cookie) Description() string {
	return Description
}

// SetSessionCookies performs the action
func (a *Cookie) SetSessionCookies() (err error) {

	z := a.Target
	z.Debug("Setting cookies %s", a.Tag)

	t := chromedp.Tasks{
		chromedp.ActionFunc(func(ctxt context.Context) error {

			expr := cdp.TimeSinceEpoch(time.Now().Add(14 * 24 * time.Hour)) // 2 weeks
			for _, c := range a.Cookies {
				z.Info("Setting %d Sessions Cookies for %s%s", len(a.Cookies), c.Domain, c.Path)

				_, err := network.SetCookie(c.Name, c.Value).
					WithDomain(c.Domain).WithExpires(&expr).WithPath(c.Path).WithSecure(c.Secure).WithHTTPOnly(c.HttpOnly).
					Do(ctxt)
				if err != nil {
					z.Error("Error setting cookie %s", c.Name)
					return err
				}
				z.Info("Setting cookie %s : %s", c.Name, c.Value)
			}
			return nil
		}),
	}

	return a.Run(t)
}
