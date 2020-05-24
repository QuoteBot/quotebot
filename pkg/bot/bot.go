package bot

import (
	"os"

	"github.com/QuoteBot/quotebot/pkg/datastorage"
)

//Bot the state of the bot
type Bot struct {
	Sc         chan os.Signal
	Conf       *BotConfig
	Commands   *BotCommands
	QuoteStore datastorage.QuoteStore
}

func NewBot(sc chan os.Signal, confFile string, commands *BotCommands) (*Bot, error) {
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
		Sc:         sc,
		Conf:       conf,
		QuoteStore: store,
		Commands:   commands,
	}, nil
}
