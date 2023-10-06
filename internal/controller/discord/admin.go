package discordhandler

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
)

var AdminDescriptor = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-admin-absence-create",
		Description: "Create an absence for a raid or multiple raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from",
				Description: "ex: 11/05/23",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to",
				Description: "ex: 15/05/23",
				Required:    false,
			},
		},
	},
	{
		Name:        "guildops-admin-absence-delete",
		Description: "Delete an absence for a raid or multiple raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from",
				Description: "(ex: 11/05/23)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to",
				Description: "(ex: 15/05/23)",
				Required:    false,
			},
		},
	},
}

func (d Discord) InitAdmin() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-admin-absence-create": d.AdminHandler,
		"guildops-admin-absence-delete": d.AdminHandler,
	}
}

func (d Discord) AdminHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user := optionMap["name"].StringValue()
	from := optionMap["from"].StringValue()
	toDate := ""
	if len(optionMap) > 2 {
		toDate = optionMap["to"].StringValue()
	}

	return d.GenerateAbsenceHandlerMsg(
		ctx, user, from, toDate, interaction.ApplicationCommandData().Name == "guildops-admin-absence-create")
}
