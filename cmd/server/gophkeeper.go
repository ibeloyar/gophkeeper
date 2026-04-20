package main

import (
	"log"
	"path/filepath"

	"github.com/ibeloyar/gophkeeper/internal/app"
	"github.com/ibeloyar/gophkeeper/internal/config"
)

func main() {
	configPath := filepath.Join("config", "server", "config.yaml")

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
