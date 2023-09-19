package discordhandler

import (
	"fmt"
	"strings"
	"time"

	"github.com/antony-ramos/guildops/internal/usecase"
)

type Discord struct {
	*usecase.AbsenceUseCase
	*usecase.PlayerUseCase
	*usecase.StrikeUseCase
	*usecase.LootUseCase
	*usecase.RaidUseCase
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
