package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	DiscordToken     string
	DiscordChannelID string
}

func Load() (*Config, error) {
	// load environment configs
	err := godotenv.Load()

	return &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
	}, err
}
