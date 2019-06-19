package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Extrusion struct {
	gorm.Model

	StartedAt  time.Time
	FinishedAt time.Time

	NecroBrowserID uint

	Status string `sql:"default:new"` // processing, completed, error

	Provider string `json:"provider" binding:"required"`

	EmailExtrusions  []EmailExtrusion
	FileExtrusions   []FileExtrusion
	GithubExtrusions []GithubExtrusion
}

// N types depending on the custom fields we need here - they need to be linked to the Extrusion
type EmailExtrusion struct {
	gorm.Model

	ExtrusionID uint

	MessageId string `json:"messageId"`

	// base64 of the HTML content of the email
	RawHtml string `json:"rawHtml" sql:"type:bytea"`

	// if a keyword was specified, and was found in the email, it's specified here
	MatchedOn string

	// if there is any open/click warning the message is stored here
	OpenWarning string `json:"openWarning"`

	ClickWarning string `json:"clickWarning"`
}

type FileExtrusion struct {
	gorm.Model

	ExtrusionID uint

	Name string

	Type string // pdf/doc/...

	// if a keyword was specified, and was found in the email, it's specified here
	MatchedOn string

	EncodedContent string `json:"content" sql:"type:bytea"` // raw base64 encoded file content

}

type GithubExtrusion struct {
	gorm.Model

	ExtrusionID uint

	Keys     string // base64 encoded /settings/keys HTML page content
	Security string // base64 encoded /settings/security HTML page content
	Profile  string
	Account  string
	Emails   string

	Repositories []GithubRepository
}

type GithubRepository struct {
	gorm.Model

	GithubExtrusionID uint
	Name              string
	Path              string

	EncodedContent string `json:"content" sql:"type:bytea"` // raw base64 encoded repo-master.zip content
}
