package bot

import (
	"os"
)

const floppy = "💾" //can use "U+1F4BE" but it should be cast in rune then in string (this is a single character)

//Bot the state of the bot
type Bot struct {
	Sc   chan os.Signal
	Conf *BotConfig
}
