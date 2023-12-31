package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Option -.
type Option func(discord *Discord)

func Token(token string) Option {
	return func(d *Discord) {
		d.token = token
	}
}

func GuildID(guildID int) Option {
	return func(d *Discord) {
		d.guildID = guildID
	}
}

func CommandHandlers(m map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error),
) Option {
	return func(d *Discord) {
		d.commandHandlers = m
	}
}

func Command(m []*discordgo.ApplicationCommand) Option {
	return func(d *Discord) {
		d.commands = m
	}
}

func DeleteCommands(b bool) Option {
	return func(d *Discord) {
		d.DeleteCommands = b
	}
}
