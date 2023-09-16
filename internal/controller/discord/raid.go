package discordHandler

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func (d Discord) InitRaid() map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error{
		"coven-raid-create": d.CreateRaidHandler,
		"coven-raid-del":    d.DeleteRaidHandler,
	}
}

var RaidDescriptors = []discordgo.ApplicationCommand{
	{
		Name:        "coven-raid-create",
		Description: "Créer un raid",
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
		Name:        "coven-raid-list",
		Description: "Lister les raids",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "ex: Milowenn",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date",
				Description: "ex: Milowenn",
				Required:    false,
			},
		},
	},
	{
		Name:        "coven-raid-del",
		Description: "Supprimer un raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: 4444-4444-4444",
				Required:    true,
			},
		},
	},
}

func (d Discord) CreateRaidHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var returnErr error
	var msg string

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name := optionMap["name"].StringValue()
	date, err := parseDate(optionMap["date"].StringValue())
	if err != nil {
		msg = "Erreur lors de la création du raid: " + err.Error()
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return returnErr
	}
	difficulty := optionMap["difficulté"].StringValue()

	raid, err := d.CreateRaid(ctx, name, difficulty, date[0])
	if err != nil {
		msg = "Erreur lors de la création du raid: " + err.Error()
		returnErr = err
	} else {
		msg = "Raid " + strconv.Itoa(raid.ID) + " créé avec succès"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return returnErr
}

func (d Discord) DeleteRaidHandler(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()
	//
	//var returnErr error
	//var msg string
	//
	//options := i.ApplicationCommandData().Options
	//optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	//
	//for _, opt := range options {
	//	optionMap[opt.Name] = opt
	//}
	//
	//id, err := uuid.Parse(optionMap["id"].StringValue())
	//if err != nil {
	//	return err
	//}
	//
	//err = d.DeleteRaid(ctx, id)
	//if err != nil {
	//	msg = "Erreur lors de la suppression du joueur: " + err.Error()
	//	returnErr = err
	//} else {
	//	msg = "Joueur " + id.String() + " supprimé avec succès"
	//}
	//
	//_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//	Type: discordgo.InteractionResponseChannelMessageWithSource,
	//	Data: &discordgo.InteractionResponseData{
	//		Content: msg,
	//	},
	//})
	//return returnErr
	return nil
}
