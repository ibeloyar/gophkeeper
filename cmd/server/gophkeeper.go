package main

import (
	"log"

	"github.com/ibeloyar/gophkeeper/internal/app"
	"github.com/ibeloyar/gophkeeper/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
