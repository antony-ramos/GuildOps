package discordhandler

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/antony-ramos/guildops/internal/entity"
)

type Discord struct {
	AbsenceUseCase
	PlayerUseCase
	StrikeUseCase
	LootUseCase
	RaidUseCase
	FailUseCase
}

type AbsenceUseCase interface {
	CreateAbsence(ctx context.Context, playerName string, date time.Time) error
	DeleteAbsence(ctx context.Context, playerName string, date time.Time) error
	ListAbsence(ctx context.Context, date time.Time) ([]entity.Absence, error)
}

type PlayerUseCase interface {
	CreatePlayer(ctx context.Context, playerName string) (int, error)
	DeletePlayer(ctx context.Context, playerName string) error
	ReadPlayer(ctx context.Context, playerName, playerLinkName string) (entity.Player, error)
	LinkPlayer(ctx context.Context, playerName string, discordID string) error
}

type RaidUseCase interface {
	CreateRaid(ctx context.Context, raidName, difficulty string, date time.Time) (entity.Raid, error)
	DeleteRaid(ctx context.Context, raidID int) error
	ReadRaid(ctx context.Context, date time.Time) (entity.Raid, error)
}

type StrikeUseCase interface {
	CreateStrike(ctx context.Context, strikeReason, playerName string) error
	DeleteStrike(ctx context.Context, id int) error
	ReadStrikes(ctx context.Context, playerName string) ([]entity.Strike, error)
}

type LootUseCase interface {
	CreateLoot(ctx context.Context, lootName string, raidDate time.Time, playerName string) error
	ListLootOnPLayer(ctx context.Context, playerName string) ([]entity.Loot, error)
	ListLootOnRaid(ctx context.Context, raidDate time.Time) ([]entity.Loot, error)
	SelectPlayerToAssign(
		ctx context.Context, playerNames []string, difficulty string,
	) (entity.Player, error)
	DeleteLoot(ctx context.Context, lootID int) error
}

type FailUseCase interface {
	CreateFail(ctx context.Context, failReason string, date time.Time, playerName string) error
	ListFailOnPLayer(ctx context.Context, playerName string) ([]entity.Fail, error)
	ListFailOnRaid(ctx context.Context, date time.Time) ([]entity.Fail, error)
	ListFailOnRaidAndPlayer(
		ctx context.Context, raidName string, playerName string,
	) ([]entity.Fail, error)
	DeleteFail(ctx context.Context, failID int) error
	UpdateFail(ctx context.Context, failID int, failReason string) error
	ReadFail(ctx context.Context, failID int) (entity.Fail, error)
}

// HumanReadableError returns the error message without the package name.
func HumanReadableError(err error) string {
	str := strings.Split(err.Error(), ": ")
	if len(str) == 1 {
		return str[0]
	} else if len(str) > 1 {
		return strings.Join(str[1:], ": ")
	}
	return err.Error()
}

// ParseDate returns a list of dates between the two dates passed as parameters.
// If toDate is empty, it will be set to fromDate.
func ParseDate(fromDate, toDate string) ([]time.Time, error) {
	layout := "02/01/06"
	startDate, err := time.Parse(layout, fromDate)
	if err != nil {
		return nil, errors.Wrap(err, "parse date")
	}

	if toDate == "" {
		toDate = fromDate
	}

	endDate, err := time.Parse(layout, toDate)
	if err != nil {
		return nil, errors.Wrap(err, "parse date")
	}

	// check if endDate is before startDate
	if endDate.Before(startDate) {
		return nil, errors.New("endDate is before startDate")
	}

	var dates []time.Time
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dates = append(dates, currentDate)
	}
	return dates, nil
}
