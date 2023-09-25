package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

// AbsenceUseCase is the use case for absences.
type AbsenceUseCase struct {
	backend Backend
}

// NewAbsenceUseCase returns a new AbsenceUseCase.
func NewAbsenceUseCase(bk Backend) *AbsenceUseCase {
	return &AbsenceUseCase{backend: bk}
}

// CreateAbsence creates an absence for a given player and date.
func (a AbsenceUseCase) CreateAbsence(ctx context.Context, playerName string, date time.Time) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("AbsenceUseCase - CreateAbsence:  ctx.Done: %w", ctx.Err())
	default:
		// Get player ID
		player, err := a.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("CreateAbsence:  backend.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("no player found")
		}

		// Get raid ID
		raids, err := a.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return fmt.Errorf("CreateAbsence:  backend.SearchRaid: %w", err)
		}
		if len(raids) == 0 {
			return fmt.Errorf("CreateAbsence:  no raid found on %s", date.Format("02-01-2006"))
		}

		// For each raid ID, create an absence
		for _, raid := range raids {
			raid := raid
			absence := entity.Absence{
				Player: &player[0],
				Raid:   &raid,
			}
			err := absence.Validate()
			if err != nil {
				return fmt.Errorf("CreateAbsence:  absence.Validate: %w", err)
			}
			_, err = a.backend.CreateAbsence(ctx, absence)
			if err != nil {
				return fmt.Errorf("CreateAbsence:  backend.CreateAbsence: %w", err)
			}
		}
		return nil
	}
}

// DeleteAbsence deletes an absence for a given player and date.
func (a AbsenceUseCase) DeleteAbsence(ctx context.Context, playerName string, date time.Time) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("AbsenceUseCase - DeleteAbsence - ctx.Done: %w", ctx.Err())
	default:
		player, err := a.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("DeleteAbsence - backend.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("no player found")
		}

		// Get absence ID
		absences, err := a.backend.SearchAbsence(ctx, "", player[0].ID, date)
		if err != nil {
			return fmt.Errorf("DeleteAbsence - backend.SearchAbsence: %w", err)
		}
		if len(absences) == 0 {
			return fmt.Errorf("no absence found")
		}

		// Delete absences
		for _, absence := range absences {
			err := a.backend.DeleteAbsence(ctx, absence.ID)
			if err != nil {
				return fmt.Errorf("DeleteAbsence - backend.DeleteAbsence: %w", err)
			}
		}
		return nil
	}
}

// ListAbsence returns a list of absences for a given date.
func (a AbsenceUseCase) ListAbsence(ctx context.Context, date time.Time) ([]entity.Absence, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("AbsenceUseCase - ListAbsence - ctx.Done: %w", ctx.Err())
	default:
		absences, err := a.backend.SearchAbsence(ctx, "", -1, date)
		if err != nil {
			return nil, fmt.Errorf("ListAbsence - backend.SearchAbsence: %w", err)
		}
		return absences, nil
	}
}
