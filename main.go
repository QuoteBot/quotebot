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

	/*conf*/
	_, err = bot.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// declare the shutdown channel
	sc := make(chan os.Signal)

	h := bot.Bot{
		Sc: sc,
	}

	// Register the messageReceived func as a callback for MessageCreate events.
	dg.AddHandler(h.MessageReceived)
	dg.AddHandler(h.GuildJoined)
	dg.AddHandler(h.ReactionAdd)

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

func startBot() {

}
