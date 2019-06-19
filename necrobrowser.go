package main

import (
	"fmt"
	"github.com/muraenateam/necrobrowser/model"

	ll "github.com/evilsocket/islazy/log"

	"github.com/muraenateam/necrobrowser/core"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/server"
)

func main() {

	// Logging
	log.SetLevel(ll.DEBUG)

	// init gorm db
	model.Init()

	// Load configuration options
	options, err := core.ParseOptions()
	if err != nil {
		log.Fatal("Error parsing options: %v", err)
	}

	// Load Docker image before starting
	err = core.InitDocker(*options.DockerImage)
	if err != nil {
		log.Fatal("Error pulling the %s image: %v", options.DockerImage, err)
	}

	// Init server
	router := server.SetupRouter(&options)

	// Run!
	address := fmt.Sprintf("%s:%s", *options.ListeningAddress, *options.ListeningPort)
	log.Info("NecroBrowser - by antisnatchor & ohpe\nWwaiting for commands on %s \\m/\nAuth Token: %s", address, *options.AuthToken)

	err = router.Run(address)
	if err != nil {
		log.Fatal("Error binding NecroBrowser: %v", err)
	}
}
