package core

import (
	"fmt"
	"github.com/muraenateam/necrobrowser/zombie/o365"

	"github.com/muraenateam/necrobrowser/zombie"
	"github.com/muraenateam/necrobrowser/zombie/dropbox"
	"github.com/muraenateam/necrobrowser/zombie/github"
	"github.com/muraenateam/necrobrowser/zombie/gsuite"
	"github.com/muraenateam/necrobrowser/zombie/owa"
	"github.com/muraenateam/necrobrowser/zombie/slack"
	"github.com/muraenateam/necrobrowser/zombie/vultr"
)

var zombies = []string{"gsuite", "github", "owa2016", "dropbox", "atlassian", "vultr", "slack", "o365"}

func GetZombie(name string, target zombie.Target, options Options) (z zombie.Zombie, err error) {

	target.Config = zombie.Config{
		LootPath: options.LootPath,
	}

	switch name {
	case "github":
		z = github.NewGithub(target)
	case "dropbox":
		z = dropbox.NewDrobox(target)
	case "vultr":
		z = vultr.NewVultr(target)
	case "slack":
		z = slack.NewSlack(target)
	case "owa2016":
		z = owa.NewOWA(target)
	case "o365":
		z = o365.NewO365(target)
	case "gsuite":
		z = gsuite.NewGSuite(target)
	}

	// Update loot path
	lp := GetZombieLootPath(options.LootPath, zombie.GetTag(z.Name()))
	if _, err := CheckLoot(lp); err != nil {
		return nil, err
	}

	z.SetLootPath(lp)

	return
}

func GetZombieLootPath(loot string, tag string) string {
	return fmt.Sprintf("%s/%s/", loot, tag)
}

func IsValidZombie(name string) bool {
	for _, value := range zombies {
		if value == name {
			return true
		}
	}
	return false
}
