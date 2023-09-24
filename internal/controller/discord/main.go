package discordhandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/antony-ramos/guildops/internal/usecase"
)

type Discord struct {
	AbsenceUseCase
	*usecase.PlayerUseCase
	*usecase.StrikeUseCase
	*usecase.LootUseCase
	*usecase.RaidUseCase
}

type AbsenceUseCase interface {
	CreateAbsence(ctx context.Context, playerName string, date time.Time) error
	DeleteAbsence(ctx context.Context, playerName string, date time.Time) error
	ListAbsence(ctx context.Context, date time.Time) ([]entity.Absence, error)
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

func parseDate(dateStr string) ([]time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	dateParts := strings.Split(dateStr, " au ")

	var dates []time.Time

	for _, datePart := range dateParts {
		date, err := time.Parse("02/01/06", datePart)
		if err != nil {
			return nil, fmt.Errorf("discord - parseDate - time.Parse: %w", err)
		}
		dates = append(dates, date)
	}

	return dates, nil
}
