package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
)

func (d Discord) InitRaid() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"guildops-raid-create": d.CreateRaidHandler,
		"guildops-raid-del":    d.DeleteRaidHandler,
		"guildops-raid-list":   d.ListRaidHandler,
	}
}

var RaidDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-raid-create",
		Description: "Create a raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Raid Milo",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: 11/05/23",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "difficulté",
				Description: "ex: HM",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-raid-del",
		Description: "Remove a raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "ex: 4546646",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-raid-list",
		Description: "List all raids on a date range",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from",
				Description: "ex: 02/10/23",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to",
				Description: "ex: 02/10/23",
				Required:    false,
			},
		},
	},
}

func (d Discord) CreateRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	returnErr := error(nil)
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name := optionMap["name"].StringValue()
	date, err := ParseDate(optionMap["date"].StringValue(), "")
	if err != nil {
		msg = "Erreur lors de la création du raid: " + HumanReadableError(err)
		if !d.Fake {
			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		}
		return fmt.Errorf("discord - CreateRaidHandler - parseDate: %w", err)
	}
	difficulty := optionMap["difficulté"].StringValue()

	raid, err := d.CreateRaid(ctx, name, difficulty, date[0])
	if err != nil {
		msg = "Erreur lors de la création du raid: " + HumanReadableError(err)
		returnErr = err
	} else {
		msg = "Raid " + strconv.Itoa(raid.ID) + " créé avec succès"
	}
	if !d.Fake {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
	return returnErr
}

func (d Discord) DeleteRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var returnErr error
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	raidID := optionMap["raidID"].IntValue()

	err := d.DeleteRaid(ctx, int(raidID))
	if err != nil {
		msg = "Erreur lors de la suppression du joueur: " + HumanReadableError(err)
		returnErr = err
	} else {
		msg = "Joueur " + strconv.Itoa(int(raidID)) + " supprimé avec succès"
	}

	if !d.Fake {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
	return returnErr
}

func (d Discord) ListRaidHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var returnErr error
	var msg string

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	from := optionMap["from"].StringValue()
	toDate := ""
	if len(optionMap) > 1 {
		toDate = optionMap["to"].StringValue()
	}

	dates, err := ParseDate(from, toDate)
	if err != nil {
		msg = "error while list raids: " + HumanReadableError(err)
		if !d.Fake {
			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		}
		return fmt.Errorf("discord - CreateRaidHandler - parseDate: %w", err)
	}
	if len(dates) == 0 {
		msg = "no date found"
		if !d.Fake {
			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		}
		return nil
	}

	var raids []entity.Raid
	for _, date := range dates {
		raid, err := d.ReadRaid(ctx, date)
		if err == nil {
			raids = append(raids, raid)
		}
	}

	msg = "Raid List:\n"
	for _, raid := range raids {
		msg += "* " + raid.Name + " " + raid.Date.Format("Mon 02/01/06") + " " + raid.Difficulty + "\n"
	}

	if len(raids) == 0 {
		msg = "no raid found"
	}

	if !d.Fake {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
	return returnErr
}
