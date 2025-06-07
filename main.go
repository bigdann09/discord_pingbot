package main

import (
	"fmt"
	"log"
	"pingbot/config"
	"pingbot/database"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// connect to the database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	fmt.Println(db)
}
