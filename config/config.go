package config

import (
	"fmt"
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

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")

	// check if variables are empty
	if user == "" || pass == "" || host == "" || port == "" || db == "" {
		return &Config{}, fmt.Errorf("database credentials are required")
	}

	// built dsn from variables
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		pass,
		host,
		port,
		db,
	)

	return &Config{
		DatabaseURL:      dsn,
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
	}, err
}
