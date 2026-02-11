package main

import (
	"flag"
	"log"
	"wyvern/server/internal/app"
	"wyvern/server/internal/config"
)

func main() {
	configPath := flag.String("config", "config.json", "server config path")
	flag.Parse()

	cfg := config.Load(*configPath)

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
