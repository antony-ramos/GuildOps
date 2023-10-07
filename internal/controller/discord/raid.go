package discordhandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/alitto/pond"
	"github.com/bwmarrin/discordgo"
)

func (d Discord) InitRaid() map[string]func(
	ctx context.Context, interaction *discordgo.InteractionCreate) (string, error) {
	return map[string]func(ctx context.Context, interaction *discordgo.InteractionCreate) (string, error){
		"guildops-raid-create":          d.CreateRaidHandler,
		"guildops-raid-delete":          d.DeleteRaidHandler,
		"guildops-raid-list":            d.ListRaidHandler,
		"guildops-raid-create-multiple": d.GenerateRaidsOnRangeHandler,
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
				Name:        "difficulty",
				Description: "Must be one of: Normal, Heroic, Mythic",
				Required:    true,
			},
		},
	},
	{
		Name:        "guildops-raid-delete",
		Description: "Remove a raid",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "ex: 902837021961355265",
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
	{
		Name:        "guildops-raid-create-multiple",
		Description: "Create multiple raids on a date range",
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
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "difficulty",
				Description: "Must be one of: Normal, Heroic, Mythic",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "weekdays",
				Description: "Must be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday",
				Required:    true,
			},
		},
	},
}

// CreateRaidHandler call an usecase to create a raid
// and return a message to the user.
// It requires a raid name, a date and a difficulty field to be passed in the interaction.
func (d Discord) CreateRaidHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Strike/StrikeOnPlayerHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	name := optionMap["name"].StringValue()
	date, err := ParseDate(optionMap["date"].StringValue(), "")
	if err != nil {
		msg := "Error while creating raid: " + HumanReadableError(err)
		return msg, fmt.Errorf("create raid parse date: %w", err)
	}
	difficulty := optionMap["difficulty"].StringValue()
	span.SetAttributes(
		attribute.String("name", name),
		attribute.String("date", date[0].String()),
		attribute.String("difficulty", difficulty),
	)

	raid, err := d.CreateRaid(ctx, name, difficulty, date[0])
	if err != nil {
		msg := "Error while creating raid: " + HumanReadableError(err)
		return msg, fmt.Errorf("call create raid usecase: %w", err)
	}
	return "Raid successfully created with ID " + strconv.Itoa(raid.ID), nil
}

// DeleteRaidHandler call an usecase to delete a raid
// and return a message to the user.
// It requires a raid ID field to be passed in the interaction.
func (d Discord) DeleteRaidHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Raid/DeleteRaidHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	raidID := optionMap["id"].StringValue()
	span.SetAttributes(
		attribute.String("raid_id", raidID),
	)
	raid, err := strconv.Atoi(raidID)
	if err != nil {
		msg := "Error while deleting raid: " + HumanReadableError(err)
		return msg, fmt.Errorf("delete raid convert user output id to int: %w", err)
	}

	err = d.DeleteRaid(ctx, raid)
	if err != nil {
		msg := "Error while deleting raid: " + HumanReadableError(err)
		return msg, fmt.Errorf("call delete raid usecase : %w", err)
	}

	return "Raid with ID " + raidID + " successfully deleted", nil
}

// ListRaidHandler call an usecase to get raids on a date range
// and return a message to the user.
// It requires a 'from' date field to be passed in the interaction.
// Optional a 'to' date field can be passed.
func (d Discord) ListRaidHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Raid/ListRaidHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

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
	span.SetAttributes(
		attribute.String("from", from),
		attribute.String("to", toDate),
	)

	dates, err := ParseDate(from, toDate)
	if err != nil {
		msg := "error while list raids: " + HumanReadableError(err)
		return msg, fmt.Errorf("list raids parse date: %w", err)
	}

	raidsLock := &sync.Mutex{}
	var raids []entity.Raid
	pool := pond.New(len(dates), 5, pond.Context(ctx))

	for _, date := range dates {
		date := date
		pool.Submit(func() {
			raid, err := d.ReadRaid(ctx, date)
			if err == nil {
				raidsLock.Lock()
				raids = append(raids, raid)
				raidsLock.Unlock()
			}
		})
	}
	pool.StopAndWaitFor(2 * time.Second)
	select {
	case <-ctx.Done():
		msg := "error while list raids: " + HumanReadableError(ctx.Err())
		return msg, fmt.Errorf("list raids wait goroutines: %w", ctx.Err())
	default:
		if len(raids) == 0 {
			msg := "no raid found"
			return msg, nil
		}

		msg := "Raid List:\n"
		for _, raid := range raids {
			msg += "* " + raid.Name + " " +
				raid.Date.Format("Mon 02/01/06") + " " +
				raid.Difficulty + " " +
				strconv.Itoa(raid.ID) + "\n"
		}
		return msg, nil
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (d Discord) GenerateRaidsOnRangeHandler(
	ctx context.Context, interaction *discordgo.InteractionCreate,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ctx, span := otel.Tracer("Discord").Start(ctx, "Raid/GenerateRaidsOnRangeHandler")
	defer span.End()
	span.SetAttributes(
		attribute.String("request_from", interaction.Member.User.Username),
	)

	select {
	case <-ctx.Done():
		msg := "error while creating multiple raids: " + HumanReadableError(ctx.Err())
		return msg, fmt.Errorf("create multiple raids wait goroutines: %w", ctx.Err())
	default:

		options := interaction.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		from := optionMap["from"].StringValue()
		toDate := ""
		if len(optionMap) > 3 {
			toDate = optionMap["to"].StringValue()
		}

		dates, err := ParseDate(from, toDate)
		if err != nil {
			msg := "error while creating multiple raids: " + HumanReadableError(err)
			return msg, fmt.Errorf("create multiple raids parse date: %w", err)
		}

		difficulty := optionMap["difficulty"].StringValue()
		difficulty = strings.ToLower(difficulty)
		if difficulty != "normal" && difficulty != "heroic" && difficulty != "mythic" {
			return "difficulty must be one of: Normal, Heroic, Mythic",
				fmt.Errorf("create multiple raids parse difficulty: %w", err)
		}

		onWeekDays := optionMap["weekdays"].StringValue()
		onWeekDays = strings.ReplaceAll(onWeekDays, " ", "")
		onWeekDays = strings.ToLower(onWeekDays)
		weekDays := strings.Split(onWeekDays, ",")
		for _, weekDay := range weekDays {
			weekDay = strings.ToLower(weekDay)
			if weekDay != "monday" && weekDay != "tuesday" &&
				weekDay != "wednesday" && weekDay != "thursday" &&
				weekDay != "friday" && weekDay != "saturday" &&
				weekDay != "sunday" {
				return "week days must be one of: Monday, Tuesday, Wednesday, " +
					"Thursday, Friday, Saturday, Sunday", fmt.Errorf("create multiple raids parse week days: %w", err)
			}
		}

		raidsDays := make([]time.Time, 0)
		for index, date := range dates {
			weekDay := date.Weekday().String()
			weekDay = strings.ToLower(weekDay)
			if contains(weekDays, weekDay) {
				raidsDays = append(raidsDays, dates[index])
			}
		}

		raidsLock := &sync.Mutex{}
		var raids []entity.Raid
		pool := pond.New(len(dates), 5, pond.Context(ctx))

		for _, date := range raidsDays {
			date := date
			pool.Submit(func() {
				raid, err := d.CreateRaid(ctx, "Raid", difficulty, date)
				if err == nil {
					raidsLock.Lock()
					raids = append(raids, raid)
					raidsLock.Unlock()
				}
			})
		}
		pool.StopAndWait()

		if len(raids) == 0 || raids == nil {
			msg := "no raid created"
			return msg, nil
		}

		msg := "Raid List:\n"
		for _, raid := range raids {
			msg += "* " + raid.Name + " " +
				raid.Date.Format("Mon 02/01/06") + " " +
				raid.Difficulty + " " +
				strconv.Itoa(raid.ID) + "\n"
		}
		return msg, nil
	}
}
