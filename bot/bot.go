package bot

import (
	"fmt"
	"pingbot/database"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session   *discordgo.Session
	db        *database.Database
	channelID string
}

func New(botToken, channelID string, db *database.Database) (*Bot, error) {
	if botToken == "" || channelID == "" {
		return &Bot{}, fmt.Errorf("discord token is not set in the environment variables")
	}

	botToken = fmt.Sprintf("Bot %s", botToken)
	session, err := discordgo.New(botToken)
	return &Bot{session, db, channelID}, err
}
