package main

import (
	"log"

	"github.com/bdreece/hopper/pkg/app"
	"github.com/bdreece/hopper/pkg/config"
)

func main() {
	cfg := config.NewConfigBuilder().
		AddCredentials().
		AddLogger().
		Build()

	db, err := app.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}
}
