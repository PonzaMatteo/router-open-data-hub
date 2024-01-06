package main

import (
	"log"
	"opendatahubchallenge/pkg/app"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatalln("Application terminated with error: ", err)
	}
}
