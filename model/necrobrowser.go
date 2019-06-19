package model

import (
	"github.com/jinzhu/gorm"
)

type NecroBrowser struct {
	gorm.Model

	Provider string `json:"provider" binding:"required"`

	DebuggingPort int `json:"debuggingPort" binding:"required"`

	EmulatedUseragent   string // which UA to spoof
	EmulatedFingerprint string // which fingerprint to spoof

	NecroTarget NecroTarget

	// each browser can be used for N extrusions
	Extrusions []Extrusion
}

type NecroTarget struct {
	gorm.Model

	NecroBrowserID uint
	Provider       string `json:"provider" binding:"required"`

	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`

	NecroCookies []NecroCookie
}

type NecroCookie struct {
	gorm.Model

	NecroTargetID uint

	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Expires  string `json:"expires"`
	Path     string `json:"path"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}
