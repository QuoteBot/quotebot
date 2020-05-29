package bot

import (
	"github.com/bwmarrin/discordgo"
)

//Commands commands struct containing the maps of executable commands
type Commands struct {
	MessageCommands        map[string]MessageCommand
	ReactionAddCommands    map[string]ReactionAddCommand
	ReactionRemoveCommands map[string]ReactionRemoveCommand
}

//MessageCommand Commands triggers by message
type MessageCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageCreate)

//ReactionAddCommand Commands triggers by new reaction
type ReactionAddCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd)

//ReactionRemoveCommand Commands triggers by new reaction
type ReactionRemoveCommand func(b *Bot, s *discordgo.Session, m *discordgo.MessageReactionRemove)
