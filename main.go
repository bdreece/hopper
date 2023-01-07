package main

import (
	"log"

	"github.com/bdreece/hopper/pkg/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
