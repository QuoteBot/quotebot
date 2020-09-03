package bot

import (
	"github.com/bwmarrin/discordgo"
	"os"

	"github.com/QuoteBot/quotebot/pkg/pagination"

	"github.com/QuoteBot/quotebot/pkg/datastorage"
)

//Bot the state of the bot
type Bot struct {
	Sc          chan os.Signal
	Conf        *Config
	Commands    *Commands
	QuoteStore  datastorage.QuoteStore
	PageManager pagination.PageManager
	session     discordgo.Session
}

//NewBot build a bot given a config file and a set of commands
func NewBot(sc chan os.Signal, confFile string, dgs *discordgo.Session, commands *Commands, defaultConfig map[string]string) (*Bot, error) {
	conf, err := loadConfig(confFile, defaultConfig)
	if err != nil {
		return nil, err
	}
	store, err := datastorage.NewQuoteStore(conf.DataPath)
	if err != nil {
		return nil, err
	}

	//TODO Commands Blacklists?

	return &Bot{
		Sc:          sc,
		Conf:        conf,
		QuoteStore:  store,
		Commands:    commands,
		PageManager: pagination.NewPageManager(dgs),
		session:     *dgs,
	}, nil
}
