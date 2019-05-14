package dropbox

import (
	"github.com/muraenateam/necrobrowser/log"

	"github.com/chromedp/chromedp"
)

// UploadFile uploads a file in a Dropbox defined location
func (z *Dropbox) UploadFile() chromedp.Tasks {
	log.Info("Uploading file %s ...", z.InjectionFile)

	t := chromedp.Tasks{}
	/*  WIP
	// This used to work but now gives issues with:
	// error: unhandled page event *page.EventNavigatedWithinDocument
	// which is EventNavigatedWithinDocument fired when same-document navigation happens,
	//  e.g. due to history API usage or anchor navigation.
	t = chromedp.Tasks{
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`div.uee-AppActionsView-SecondaryActionMenu-text-upload-file`, z.InjectionFile, chromedp.NodeVisible),
		chromedp.Sleep(2 * time.Second),
		//chromedp.Click(`input[name="submit"]`),
	}
	*/

	return t
}
