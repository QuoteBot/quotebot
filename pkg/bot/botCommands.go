package bot

import (
	"github.com/bwmarrin/discordgo"
)

type BotCommands struct {
	MessageCommands      map[string]MessageCommand
	ReractionAddCommands map[string]ReractionAddCommand
}

//MessageCommand Commands triggers by message
type MessageCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageCreate)

//ReractionAddCommand Commands triggers by new reaction
type ReractionAddCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd)
