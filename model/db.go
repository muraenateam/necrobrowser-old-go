package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

var (
	Session *gorm.DB
	user    = "necrobrowser"
	passwd  = "necromancing-your-way-through"
	db      = "necrobrowser"
	host    = "localhost"
)

func Init() {
	Connect()
	Migrate()
}

func Connect() {

	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, passwd, db)
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		log.Fatalf("There was an error connecting to Postgres:\n %s", err)
		panic(err.Error())
	}

	Session = db
}

func Migrate() {
	Session.AutoMigrate(&Instrumentation{})
	Session.AutoMigrate(&NecroBrowser{})
	Session.AutoMigrate(&NecroTarget{})
	Session.AutoMigrate(&NecroCookie{})
	Session.AutoMigrate(&Extrusion{})
	Session.AutoMigrate(&EmailExtrusion{})
	Session.AutoMigrate(&FileExtrusion{})
	Session.AutoMigrate(&GithubExtrusion{})
	Session.AutoMigrate(&GithubRepository{})
}
