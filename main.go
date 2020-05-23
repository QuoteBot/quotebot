package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/config"
	"github.com/QuoteBot/quotebot/pkg/bot/command"

	"github.com/bwmarrin/discordgo"
)

func main() {
	tokenFile := flag.String("token", "token", "path to the token file")
	configFile := flag.String("config", "config.json", "path to the config file (json)")

	flag.Parse()

	token, err := config.LoadToken(*tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := bot.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	startBot(token, conf)

}

func startBot(token string, conf *bot.BotConfig) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// declare the shutdown channel
	sc := make(chan os.Signal)

	b := bot.Bot{
		Sc:   sc,
		Conf: conf,
		Commands: &bot.BotCommands{
			MessageCommands: command.MessageCommands(),
		},
	}

	// Register the messageReceived func as a callback for MessageCreate events.
	dg.AddHandler(b.MessageReceived)
	dg.AddHandler(b.GuildJoined)
	dg.AddHandler(b.ReactionAdd)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc // wait until a Signal is put in the channel

	// Cleanly close down the Discord session.
	dg.Close()
}
