package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pingbot/database"
	"time"

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

func (bot *Bot) Run() {
	if bot.session == nil {
		log.Fatal("discord session is not initialized")
	}

	err := bot.session.Open()
	if err != nil {
		log.Fatal("error opening discord session: %w", err)
	}
	defer bot.session.Close()

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	bot.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return // Ignore messages from the bot itself
		}

		if m.Content == "!ping" {

			response := "Pong! The bot is alive."
			bot.session.ChannelMessageSend(m.ChannelID, response)

		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
