package file

import (
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
)

const (
	// Name of this action
	Name = "File"

	// Description of this action
	Description = "File based actions"
)

// Screenshot is an action
type File struct {
	action.Action

	Selector string
	URL      string
	FileURL  string
}

// Name returns the action name
func (a *File) Name() string {
	return Name
}

// Description returns what the action does
func (a *File) Description() string {
	return Description
}

// Download downloads a file
func (a *File) Download() (err error) {
	z := a.Target
	z.Info("Downloading a file from %s", a.FileURL)

	//   page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorAllow).WithDownloadPath("."),

	/*
	   allocContext, _ := chromedp.NewExecAllocator(context.Background())
	   ctx, cancel := chromedp.NewContext(allocContext)
	   defer cancel()

	   err = chromedp.Run(cxt, page.SetDownloadBehavior("allow").WithDownloadPath("/home/user/myDownloadFolder"))
	   if err != nil {
	       // Handle error
	   }
	   err = chromedp.Run(cxt, Navigate("www.google.com"))

	*/

	t := chromedp.Tasks{
		chromedp.Navigate(a.URL),
		chromedp.WaitVisible(a.Selector, chromedp.ByQuery),
	}

	return a.Run(t)
}
