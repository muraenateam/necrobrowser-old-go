package core

import (
	"bufio"
	"os"

	"github.com/muraenateam/necrobrowser/log"
)

func Debug(params ...string) {
	message := "Press 'Enter' to continue..."
	if len(params) > 0 && params[0] != "" {
		message = params[0]
	}
	log.Debug(message)
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
