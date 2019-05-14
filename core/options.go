package core

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/muraenateam/necrobrowser/log"
)

const (
	AuthToken = "ada9f7b8-6e6c-4884-b2a3-ea757c1eb617"
)

type Options struct {
	Debug    *bool
	Headless *bool

	DockerImage *string
	UserAgent   *string
	AuthToken   *string

	ListeningAddress *string
	ListeningPort    *string

	LootPath string
}

func getHerePwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	return dir
}

func ParseOptions() (Options, error) {
	o := Options{
		Debug:    flag.Bool("debug", false, "Print debug messages."),
		Headless: flag.Bool("headless", false, "If headless is true, expects an headless Zombie instance with remote debugging on localhost:9222"),

		DockerImage: flag.String("docker", "registry.hub.docker.com/zenika/alpine-chrome:latest", "Docker image"),
		UserAgent:   flag.String("useragent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Zombie/73.0.3683.86 Safari/537.36", "UserAgent string"),
		AuthToken:   flag.String("token", "", "Authentication token"),

		ListeningAddress: flag.String("laddress", "0.0.0.0", "Listening address where expose the Necrobrowser APIs"),
		ListeningPort:    flag.String("lport", "8080", "TCP Port to bind the listening server"),

		//		LootPath: flag.String("loot", "./loot/", "Where to store collected resources and logs"),
	}

	flag.Parse()

	if *o.AuthToken == "" {
		*o.AuthToken = AuthToken
	}

	o.LootPath = filepath.Join(getHerePwd(), "loot")
	return o, nil
}
