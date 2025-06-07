package bot

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pingbot/database"
	"strings"
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

		} else if strings.Contains(m.Content, "!add") {

			// Extract the hostname from the message
			content := strings.Split(m.Content, " ")
			if len(content) < 2 {
				bot.session.ChannelMessageSend(m.ChannelID, "Please provide a hostname to add. Usage: `!add <hostname>`")
				return
			}

			hostname := content[1]
			fmt.Println("Adding server:", hostname)
			if err := bot.add(hostname); err != nil {
				log.Println("Failed to add server:", err)
				bot.session.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			log.Printf("Server %s added successfully.\n", hostname)
			bot.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server %s added successfully.", "hostname"))

		} else if strings.Contains(m.Content, "!status") {

			// Extract the hostname from the message
			content := strings.Split(m.Content, " ")
			if len(content) < 2 {
				bot.session.ChannelMessageSend(m.ChannelID, "Please provide a hostname to check status. Usage: `!status <hostname>`")
				return
			}

			hostname := content[1]
			fmt.Println("Checking status for server:", hostname)
			if err := bot.monitor(hostname); err != nil {
				log.Printf("Server %s is down or invalid host: %v\n", hostname, err)
				bot.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("⚠️ %s is down or invalid host: %v", hostname, err))
				return
			}

			log.Printf("Server %s is reachable.\n", hostname)
			bot.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ %s is up 🌐\n", hostname))

		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// monitor hostname
func (bot *Bot) monitor(hostname string) error {
	if _, err := net.LookupHost(hostname); err != nil {
		fmt.Println("DNS lookup failed:", err)
		return fmt.Errorf("dns server is not reachable for %s. please check your dns settings. %v", hostname, err)
	}

	// Check if the database is reachable
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", hostname), 5*time.Second)
	if err != nil {
		return fmt.Errorf("server %s is not reachable: %v", hostname, err)
	}
	return err
}

// add hostname to database
func (bot *Bot) add(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("please provide a hostname to add")
	}

	// lookup host
	if _, err := net.LookupHost(hostname); err != nil {
		return fmt.Errorf("failed to resolve hostname %s: %v", hostname, err)
	}

	// Add the server to the database
	if err := bot.db.AddServer(hostname); err != nil {
		return fmt.Errorf("failed to add server: %v", err)
	}

	return nil
}
