package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/bwmarrin/discordgo"
)

var PlayerDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "guildops-player-create",
		Description: "Créer un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-delete",
		Description: "Supprimer un joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-get",
		Description: "Infos sur le joueur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-link",
		Description: "link your discord account to your player name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: milowenn",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-player-info",
		Description: "Show info about yourself",
	},
}

func (d Discord) InitPlayer() map[string]func(
	ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) error{
		"guildops-player-create": d.PlayerHandler,
		"guildops-player-delete": d.PlayerHandler,
		"guildops-player-get":    d.GetPlayerHandler,
		"guildops-player-link":   d.LinkPlayerHandler,
		"guildops-player-info":   d.GetPlayerHandler,
	}
}

func (d Discord) PlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	var returnErr error
	name := optionMap["name"].StringValue()
	if interaction.ApplicationCommandData().Name == "guildops-player-create" {
		id, err := d.CreatePlayer(ctx, name)
		if err != nil {
			msg = "Erreur lors de la création du joueur: " + HumanReadableError(err)
			returnErr = err
		} else {
			msg = "Joueur " + name + " créé avec succès : ID " + strconv.Itoa(id)
		}
	} else {
		err := d.DeletePlayer(ctx, name)
		if err != nil {
			msg = "Erreur lors de la suppression du joueur: " + HumanReadableError(err)
			returnErr = err
		} else {
			msg = "Joueur " + name + " supprimé avec succès"
		}
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

func (d Discord) GetPlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	var name string
	var player entity.Player
	var err error

	commandInfo := "guildops-player-info"

	if interaction.ApplicationCommandData().Name == commandInfo {
		player, err = d.ReadPlayer(ctx, "", interaction.Member.User.Username)
	} else {
		name = optionMap["name"].StringValue()
		player, err = d.ReadPlayer(ctx, name, "")
	}

	// Show on string all info about player
	if err != nil {
		msg = "Erreur lors de la récupération du joueur: " + HumanReadableError(err)
		if !d.Fake {
			data := discordgo.InteractionResponseData{
				Content: msg,
			}
			if interaction.ApplicationCommandData().Name == commandInfo {
				data = discordgo.InteractionResponseData{
					Content: msg,
					Flags:   discordgo.MessageFlagsEphemeral,
				}
			}

			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &data,
			})
		}
		return fmt.Errorf("database - GetPlayerHandler - r.ReadPlayer: %w", err)
	}
	msg += "Name : **" + player.Name + "**\n"
	msg += "ID : **" + strconv.Itoa(player.ID) + "**\n"
	if player.DiscordName != "" {
		msg += "Discord ID : **" + interaction.Member.User.ID + "**\n"
	}

	// For each difficulty, show the number of loots
	lootCounter := make(map[string]int)
	for _, loot := range player.Loots {
		lootCounter[loot.Raid.Difficulty]++
	}
	if len(lootCounter) > 0 {
		msg += "**Loots Count:** \n"
		for difficulty, count := range lootCounter {
			msg += "*  " + difficulty + " | " + strconv.Itoa(count) + " loots \n"
		}
	}

	if len(player.Strikes) > 0 {
		msg += "**Strikes (" + strconv.Itoa(len(player.Strikes)) + ") :** \n"
		for _, strike := range player.Strikes {
			msg += "*  " + strike.Date.Format("02/01/2006") +
				" | " + strike.Reason + " | " + strike.Season + " | " + strconv.Itoa(strike.ID) + "\n"
		}
	}
	if len(player.MissedRaids) > 0 {
		msg += "**Absences (" + strconv.Itoa(len(player.MissedRaids)) + ") :** \n"
		for _, raid := range player.MissedRaids {
			msg += "*  " + raid.Date.Format("02/01/06") +
				" | " + raid.Difficulty +
				" | " + raid.Name + "\n"
		}
	}

	if len(player.Loots) > 0 {
		msg += "**Loots (" + strconv.Itoa(len(player.Loots)) + ") :** \n"
		for _, loot := range player.Loots {
			msg += "*  " + loot.Raid.Difficulty +
				" | " + loot.Raid.Date.Format("02/01/06") +
				" | " + loot.Name + "\n"
		}
	}

	if len(player.Fails) > 0 {
		msg += "**Fails (" + strconv.Itoa(len(player.Fails)) + ") :** \n"
		for _, fail := range player.Fails {
			msg += "*  " + fail.Raid.Date.Format("02/01/2006") +
				" | " + fail.Reason + "\n"
		}
	}

	if !d.Fake {
		data := discordgo.InteractionResponseData{
			Content: msg,
		}
		if interaction.ApplicationCommandData().Name == commandInfo {
			data = discordgo.InteractionResponseData{
				Content: msg,
				Flags:   discordgo.MessageFlagsEphemeral,
			}
		}
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &data,
		})
	}
	return nil
}

func (d Discord) GenerateLinkPlayerMsg(ctx context.Context, discordName, playerName string) (string, error) {
	select {
	case <-ctx.Done():
		return ctxError,
			fmt.Errorf("discord - GenerateLinkPlayerMsg - ctx.Done: request took too much time to be proceed")
	default:
		var msg string

		player, err := d.ReadPlayer(ctx, playerName, "")
		if err != nil {
			msg = "Error while reading player: " + HumanReadableError(err)
			return msg, fmt.Errorf("discord - LinkPlayerHandler - r.ReadPlayer: %w", err)
		}

		err = d.LinkPlayer(ctx, player.Name, discordName)
		if err != nil {
			msg = "Error while linking player: " + HumanReadableError(err)
			return msg, fmt.Errorf("discord - LinkPlayerHandler - r.LinkPlayer: %w", err)
		}

		msg += "You are now linked to this player : \n"
		msg += "Name : **" + player.Name + "**\n"
		msg += "ID : **" + strconv.Itoa(player.ID) + "**\n"
		if player.DiscordName != "" {
			msg += "Discord ID : **" + discordName + "**\n"
		}
		return msg, nil
	}
}

func (d Discord) LinkPlayerHandler(
	ctx context.Context, session *discordgo.Session, interaction *discordgo.InteractionCreate,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	playerName := optionMap["name"].StringValue()
	discordName := interaction.Member.User.Username

	msg, err := d.GenerateLinkPlayerMsg(ctx, discordName, playerName)

	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return err
}
