package screenshot

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "Screenshot"

	// Description of this action
	Description = "Screenshot takes a picture of the screen at the given URL"
)

// Screenshot is an action
type Screenshot struct {
	action.Action

	URL      string
	Selector string
	SavePath string
}

// Name returns the action name
func (a *Screenshot) Name() string {
	return Name
}

// Description returns what the action does
func (a *Screenshot) Description() string {
	return Description
}

// Take performs the action to take a screenshot
func (a *Screenshot) Take() (err error) {

	var (
		buf []byte
		z   = a.Target
	)

	// Define save path if not defined
	if a.SavePath == "" {
		a.SavePath = fmt.Sprintf("%s/%s.png", a.LootPath, action.Now())
	}

	z.Info("Taking screenshot of page %s", a.URL)
	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		chromedp.Sleep(2 * time.Second),
	}

	if a.Selector != "" {
		act := chromedp.WaitVisible(a.Selector, chromedp.ByQuery)
		t = append(t, act)
	}

	act := chromedp.CaptureScreenshot(&buf)
	t = append(t, act)
	if err = a.Run(t); err != nil {
		return
	}

	err = ioutil.WriteFile(a.SavePath, buf, 0644)
	if err != nil {
		return
	}

	z.Warning("Screenshot saved to: %s", a.SavePath)
	return nil
}
