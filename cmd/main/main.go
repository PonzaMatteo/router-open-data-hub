package main

import (
	"log"

	"github.com/PonzaMatteo/router-open-data-hub/pkg/app"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatalln("Application terminated with error: ", err)
	}
}
