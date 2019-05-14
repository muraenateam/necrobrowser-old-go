package zombie

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/muraenateam/necrobrowser/log"
)

// Zombie is the interface
type Zombie interface {
	Instrument() (interface{}, error)
	Name() string
	SetLootPath(string)
	SetContext(ctx context.Context)
}

type Target struct {
	Context context.Context
	Cookies []SessionCookie
	Config
}

type Config struct {
	LootPath string
	Tag      string
}

type SessionCookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Expires  string `json:"expires"`
	Path     string `json:"path"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}

type Extrusion struct{}

// GetTag returns the slug equivalent of a given string
func GetTag(s string) string {
	var re = regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func formatTag(t string) string {
	return fmt.Sprintf("[%s] ", t)
}

func (m *Target) Debug(format string, args ...interface{}) {
	log.Debug(formatTag(m.Config.Tag)+format, args...)
}

func (m *Target) Info(format string, args ...interface{}) {
	log.Info(formatTag(m.Config.Tag)+format, args...)
}

func (m *Target) Important(format string, args ...interface{}) {
	log.Important(formatTag(m.Config.Tag)+format, args...)
}

func (m *Target) Warning(format string, args ...interface{}) {
	log.Warning(formatTag(m.Config.Tag)+format, args...)
}

func (m *Target) Error(format string, args ...interface{}) {
	log.Error(formatTag(m.Config.Tag)+format, args...)
}

func (m *Target) Fatal(format string, args ...interface{}) {
	log.Fatal(formatTag(m.Config.Tag)+format, args...)
}
