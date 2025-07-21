package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/thoj/go-ircevent"

	"github.com/allthingslinux/ircBot/commands"
)

// Holds all configuration values for the bot
type Config struct {
	DiscordToken string
	IRCOperPass  string
	IRCUser      string
	IRCPassword  string
	BotPrefix    string
	IRCServer    string
	IRCPort      string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Discord bot
	discord, err := initDiscordBot(config.DiscordToken)
	if err != nil {
		log.Fatalf("Failed to initialize Discord bot: %v", err)
	}
	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("Error closing Discord session: %v", err)
		}
	}()

	// Initialize IRC connection
	ircCon, err := initIRCConnection(config)
	if err != nil {
		log.Fatalf("Failed to initialize IRC connection: %v", err)
	}

	// Register command handlers
	discord.AddHandler(commands.CommandMapper(ircCon, config.BotPrefix))

	log.Println("Bot is running. Press CTRL+C to exit.")

	waitForShutdown(ctx, cancel)
}

// loadConfig loads and validates configuration from environment variables
func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	config := &Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		IRCOperPass:  os.Getenv("IRC_OPER_PASS"),
		IRCUser:      os.Getenv("IRC_BOT_USER"),
		IRCPassword:  os.Getenv("IRC_BOT_PASSWORD"),
		BotPrefix:    os.Getenv("BOT_PREFIX"),
		IRCServer:    os.Getenv("IRC_SERVER"),
		IRCPort:      os.Getenv("IRC_PORT"),
	}

	// Validate variables
	if config.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}
	if config.IRCUser == "" {
		return nil, fmt.Errorf("IRC_BOT_USER is required")
	}
	if config.BotPrefix == "" {
		return nil, fmt.Errorf("BOT_PREFIX is required")
	}
	if config.IRCServer == "" {
		return nil, fmt.Errorf("IRC_SERVER is required")
	}
	if config.IRCPort == "" {
		return nil, fmt.Errorf("IRC_PORT is required")
	}

	return config, nil
}

// Initialises and opens a Discord session
func initDiscordBot(token string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	if err := discord.Open(); err != nil {
		return nil, fmt.Errorf("failed to open Discord session: %w", err)
	}

	return discord, nil
}

// Connects to IRC server
func initIRCConnection(config *Config) (*irc.Connection, error) {
	ircCon := irc.IRC(config.IRCUser, config.IRCUser)
	ircCon.VerboseCallbackHandler = false
	ircCon.UseTLS = true
	ircCon.TLSConfig = &tls.Config{
		ServerName: config.IRCServer,
	}
	ircCon.Password = config.IRCPassword

	// Set up connection callback for authentication
	ircCon.AddCallback("001", func(e *irc.Event) {
		if config.IRCPassword != "" && config.IRCUser != "" {
			ircCon.Privmsg("NickServ", fmt.Sprintf("IDENTIFY %s %s", config.IRCUser, config.IRCPassword))
			if config.IRCOperPass != "" {
				ircCon.SendRaw(fmt.Sprintf("OPER atl-mod %s", config.IRCOperPass))
			}
		}
	})

	serverAddr := config.IRCServer + ":" + config.IRCPort
	if err := ircCon.Connect(serverAddr); err != nil {
		return nil, fmt.Errorf("failed to connect to IRC server %s: %w", serverAddr, err)
	}

	go ircCon.Loop()
	return ircCon, nil
}

// Graceful shutdown handling
func waitForShutdown(ctx context.Context, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		cancel()
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}
}
