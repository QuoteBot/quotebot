package bot

import (
	"github.com/bwmarrin/discordgo"
)

type BotCommands struct {
	MessageCommands map[string]MessageCommand
	//ReactionCommands
}

//MessageCommand Commands triggers by message
type MessageCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageCreate)
