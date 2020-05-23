package bot

import "github.com/bwmarrin/discordgo"

//GuildJoined Handle when join a guild
func (b *Bot) GuildJoined(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	s.ChannelMessageSend(event.Guild.SystemChannelID, "QuoteBot ready")
}
