package usecase

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"time"
)

type AbsenceUseCase struct {
	backend Backend
}

func NewAbsenceUseCase(bk Backend) *AbsenceUseCase {
	return &AbsenceUseCase{backend: bk}
}

func (a AbsenceUseCase) CreateAbsence(ctx context.Context, playerName string, date time.Time) error {
	// Get player ID
	player, err := a.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return err
	}
	if len(player) == 0 {
		return fmt.Errorf("no player found")
	}
	// Get raid ID
	raids, err := a.backend.SearchRaid(ctx, "", date, "")
	if err != nil {
		return err
	}
	if len(raids) == 0 {
		return fmt.Errorf("no raid found on this date %s", date)
	}
	// For each raid ID, create an absence
	for _, raid := range raids {
		absence := entity.Absence{
			Player: &player[0],
			Raid:   &raid,
		}
		err := absence.Validate()
		if err != nil {
			return err
		}
		absence, err = a.backend.CreateAbsence(ctx, absence)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a AbsenceUseCase) DeleteAbsence(ctx context.Context, playerName string, date time.Time) error {
	// Get player ID
	player, err := a.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return err
	}
	if len(player) == 0 {
		return fmt.Errorf("no player found")
	}

	absences, err := a.backend.SearchAbsence(ctx, "", player[0].ID, date)
	if err != nil {
		return err
	}
	if len(absences) == 0 {
		return fmt.Errorf("no absence found")
	}
	for _, absence := range absences {
		err := a.backend.DeleteAbsence(ctx, absence.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// List Absences on a specific date
func (a AbsenceUseCase) ListAbsence(ctx context.Context, date time.Time) ([]entity.Absence, error) {
	absences, err := a.backend.SearchAbsence(ctx, "", -1, date)
	if err != nil {
		return nil, err
	}
	return absences, nil
}
