package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/QuoteBot/quotebot/pkg/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	tokenFile := flag.String("token", "token", "path to the token file")

	flag.Parse()

	token, err := config.LoadToken(*tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	sc := make(chan os.Signal)

	h := handlerWithSession{
		Session: dg,
		Sc:      sc,
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(h.messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type handlerWithSession struct {
	Session *discordgo.Session
	Sc      chan os.Signal
}

func (h *handlerWithSession) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "shutdown" {
		h.Sc <- syscall.SIGINT
	}
}
