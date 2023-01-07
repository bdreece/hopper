package main

import (
	"log"

	"github.com/bdreece/hopper/pkg/app"
)

func main() {
	db, err := app.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}
}
