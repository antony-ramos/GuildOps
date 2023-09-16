package discordHandler

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var PlayerDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-player-create",
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
		Name:        "coven-player-delete",
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
		Name:        "coven-player-get",
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
}

func (d Discord) InitPlayer() map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-player-create": d.CreatePlayerHandler,
		"coven-player-delete": d.DeletePlayerHandler,
		"coven-player-get":    d.GetPlayerHandler,
	}
}

func (d Discord) CreatePlayerHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	name := optionMap["name"].StringValue()
	err := d.CreatePlayer(ctx, name, "")
	var returnErr error
	if err != nil {
		msg = "Erreur lors de la création du joueur: " + err.Error()
		returnErr = err
	} else {
		msg = "Joueur " + name + " créé avec succès"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) DeletePlayerHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	name := optionMap["name"].StringValue()
	err := d.DeletePlayer(ctx, name)
	var returnErr error
	if err != nil {
		msg = "Erreur lors de la suppression du joueur: " + err.Error()
		returnErr = err
	} else {
		msg = "Joueur " + name + " supprimé avec succès"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) GetPlayerHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var msg string
	name := optionMap["name"].StringValue()
	player, err := d.ReadPlayer(ctx, name)
	// Show on string all info about player
	if err != nil {
		msg = "Erreur lors de la récupération du joueur: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return err
	}
	msg += "Name : **" + player.Name + "**\n"
	msg += "ID : **" + strconv.Itoa(player.ID) + "**\n"

	// For each difficulty, show the number of loots
	lootCounter := make(map[string]int)
	for _, loot := range player.Loots {
		lootCounter[loot.Raid.Difficulty]++
	}
	if len(lootCounter) > 0 {
		msg += "**Loots Count:** \n"
		for difficulty, count := range lootCounter {
			msg += "  " + difficulty + " | " + strconv.Itoa(count) + " loots \n"
		}
	}

	if len(player.Strikes) > 0 {
		msg += "**Strikes (" + strconv.Itoa(len(player.Strikes)) + ") :** \n"
		for _, strike := range player.Strikes {
			msg += "  " + strike.Reason +
				" | " + strike.Date.Format("02/01/06") + " | " + strike.Season + " | " + strconv.Itoa(strike.ID) + "\n"
		}
	}
	if len(player.MissedRaids) > 0 {
		msg += "**Absences (" + strconv.Itoa(len(player.MissedRaids)) + ") :** \n"
		for _, raid := range player.MissedRaids {
			msg += "  " + raid.Name +
				" | " + raid.Difficulty +
				" | " + raid.Date.Format("02/01/06") + "\n"
		}
	}

	if len(player.Loots) > 0 {
		msg += "**Loots :** \n"
		for _, loot := range player.Loots {

			msg += "  " + loot.Raid.Difficulty +
				" | " + loot.Raid.Date.Format("02/01/06") +
				" | " + loot.Name + "\n"
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return nil
}
