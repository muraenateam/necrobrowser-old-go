package github

import (
	"time"

	"github.com/muraenateam/necrobrowser/log"

	"github.com/chromedp/chromedp"
)

func (z *Github) addSSHKey() chromedp.Tasks {

	log.Info("Adding new key (%s):[%s]...", z.SSHKeyName, z.SSHKeyValue)
	return chromedp.Tasks{
		chromedp.Navigate(z.SSHURLPath),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(`form[class="new_public_key"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[id="public_key_title"]`, z.SSHKeyName),

		chromedp.Sleep(500 * time.Millisecond),
		chromedp.SendKeys(`textarea[id="public_key_key"]`, z.SSHKeyValue),

		chromedp.Sleep(5000 * time.Millisecond),
		chromedp.Click(`#new_key > p > button`),
		chromedp.Sleep(2000 * time.Millisecond),
	}
}
