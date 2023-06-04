package main

import (
	"fmt"
	"log"

	"core/internal/app"
	"core/internal/config"
)

var version = "dev"

func main() {
	fmt.Printf("Version: %s\n", version)

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
