package bot

import (
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
	PageHandler pagination.PageHandler
}

//NewBot build a bot given a config file and a set of commands
func NewBot(sc chan os.Signal, confFile string, commands *Commands) (*Bot, error) {
	conf, err := loadConfig(confFile)
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
		PageHandler: pagination.NewPageHandler(),
	}, nil
}
