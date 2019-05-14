package action

import (
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/zombie"
)

// Action is the action structure
type Action struct {
	zombie.Target
}

// Run executes the action against the current target using the supplied context
func (a *Action) Run(t chromedp.Tasks) (err error) {

	z := a.Target
	z.Debug("Running")

	// Set default download path
	t = append(t, page.SetDownloadBehavior("allow").WithDownloadPath("/tmp"))

	err = chromedp.Run(a.Context, t)
	if err != nil {
		return err
	}

	z.Debug("Done")
	return nil
}

// Now returns time.Now() formatted for loot files
func Now() string {
	t := time.Now()
	return t.Format("20060102-150405")
}
