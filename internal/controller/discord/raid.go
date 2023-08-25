package discordHandler

//
//import (
//	"fmt"
//	"github.com/bwmarrin/discordgo"
//)
//
//var RaidDescriptor = discordgo.ApplicationCommand{
//	Name:        "raid",
//	Description: "Lister les absences pour un raid",
//	Options: []*discordgo.ApplicationCommandOption{
//		{
//			Type:        discordgo.ApplicationCommandOptionString,
//			Name:        "date",
//			Description: "(ex: 11-05-23 | ou 11-05-23 au 13-05-23)",
//			Required:    true,
//		},
//	},
//}
//
//func RaidHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
//
//	options := i.ApplicationCommandData().Options
//	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
//	for _, opt := range options {
//		optionMap[opt.Name] = opt
//	}
//
//	// Créer le thread avec le message parent
//	_, err := s.ThreadStart(i.ChannelID, "toto", discordgo.ChannelTypeGuildPublicThread, 60)
//	if err != nil {
//		// Gérer l'erreur
//		return
//	}
//
//	discordMessage := "Unknown Error"
//
//	beginDate, _, err := utils.ExtractDates(optionMap["date"].StringValue())
//	if err != nil {
//		discordMessage = fmt.Sprintf("date format is invalid : %s", err.Error())
//	} else {
//		id, err := b.Pub(adapters.GetRaidAbs, beginDate)
//		err = b.Sub(id, func(message adapters.BusMessage) bool {
//			if tmp, ok := message.Data.([]interface{}); ok {
//				var abs []string
//				for _, v := range tmp {
//					if str, ok := v.(string); ok {
//						abs = append(abs, str)
//					}
//				}
//
//				if abs[0] == "NO_ABS" {
//					discordMessage = fmt.Sprintf("Aucune absence au raid du %s %d %s %d \n", utils.GetDay(beginDate), beginDate.Day(), utils.GetMonthFR(beginDate), beginDate.Year())
//				} else {
//					discordMessage = fmt.Sprintf("Absences au raid du %s %d %s %d :\n", utils.GetDay(beginDate), beginDate.Day(), utils.GetMonthFR(beginDate), beginDate.Year())
//					for _, pseudo := range abs {
//						discordMessage += fmt.Sprintf("- %s\n", pseudo)
//					}
//				}
//				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//					// Ignore type for now, they will be discussed in "responses"
//					Type: discordgo.InteractionResponseChannelMessageWithSource,
//					Data: &discordgo.InteractionResponseData{
//						Content: fmt.Sprintf(
//							discordMessage,
//						),
//					},
//				})
//				return false
//			}
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				// Ignore type for now, they will be discussed in "responses"
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: fmt.Sprintf(discordMessage),
//					Flags:   discordgo.MessageFlagsEphemeral,
//				},
//			})
//			return true
//		})
//		if err != nil {
//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				// Ignore type for now, they will be discussed in "responses"
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: fmt.Sprintf(discordMessage),
//					Flags:   discordgo.MessageFlagsEphemeral,
//				},
//			})
//			//TODO log
//			return
//		}
//	}
//}
