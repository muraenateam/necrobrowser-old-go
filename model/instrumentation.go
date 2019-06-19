package model

import (
	"github.com/jinzhu/gorm"
)

type Instrumentation struct {
	gorm.Model

	ClientIP string
	ClientUA string

	Provider string `json:"provider" binding:"required"`
}
