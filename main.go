package main

import (
	"log"
	"pingbot/bot"
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

	// connect to discord bot
	discordBot, err := bot.New(cfg.DiscordToken, cfg.DiscordChannelID, db)
	if err != nil {
		panic(err)
	}

	// start the bot
	log.Println("Starting the Discord bot...")
	discordBot.Monitor()
	discordBot.Run()
}
