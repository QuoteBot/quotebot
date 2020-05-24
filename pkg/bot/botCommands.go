package bot

import (
	"github.com/bwmarrin/discordgo"
)

//BotCommands commands struct containing the maps of executable commands
type BotCommands struct {
	MessageCommands        map[string]MessageCommand
	ReractionAddCommands   map[string]ReractionAddCommand
	ReactionRemoveCommands map[string]ReractionRemoveCommand
}

//MessageCommand Commands triggers by message
type MessageCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageCreate)

//ReractionAddCommand Commands triggers by new reaction
type ReractionAddCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd)

//ReractionRemoveCommand Commands triggers by new reaction
type ReractionRemoveCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageReactionRemove)
