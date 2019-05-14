package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/muraenateam/necrobrowser/core"
)

var (
	ErrWrongAuthentication = errors.New("wrong or missing authentication token")
)

func AuthMiddleware(t uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("authenticationToken", t)
		c.Next()
	}
}

func OptionsMiddleware(o *core.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("options", *o)
		c.Next()
	}
}

func Authenticate(c *gin.Context) bool {
	token := c.Params.ByName("token")
	sessionToken := c.MustGet("authenticationToken").(uuid.UUID).String()
	if token != sessionToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrWrongAuthentication})
		return false
	}

	return true
}

func GetOptions(c *gin.Context) core.Options {
	return c.MustGet("options").(core.Options)
}
