package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

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
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Absence/CreateAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.String("date", date.Format("02/01/06")),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("AbsenceUseCase - CreateAbsence:  ctx.Done: request took too much time to be proceed")
	default:
		if date.Before(
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())) {
			return errors.New("can't create a absence in the past")
		}

		player, err := entity.NewPlayer(-1, playerName, "")
		if err != nil {
			return fmt.Errorf("create player entity: %w", err)
		}

		// Get player ID
		players, err := a.backend.SearchPlayer(ctx, -1, player.Name, "")
		if err != nil {
			return fmt.Errorf("check if player exists in database: %w", err)
		}
		if len(players) == 0 {
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
			absence, err := entity.NewAbsence(-1, &players[0], &raid)
			if err != nil {
				return fmt.Errorf("create absence object: %w", err)
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
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Absence/DeleteAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.String("date", date.Format("02/01/06")),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("AbsenceUseCase - DeleteAbsence - ctx.Done: request took too much time to be proceed")
	default:
		player, err := a.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("check if player exists in database: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("no player found")
		}

		raid, err := a.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return fmt.Errorf("DeleteAbsence - backend.SearchRaid: %w", err)
		}
		if len(raid) == 0 {
			return fmt.Errorf("no raid found")
		}

		// Get absence ID
		absences, err := a.backend.SearchAbsence(ctx, "", player[0].ID, date)
		if err != nil {
			return fmt.Errorf("DeleteAbsence - backend.SearchAbsence: %w", err)
		}
		if len(absences) == 0 {
			return fmt.Errorf("absence not found")
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
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Absence/ListAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.String("date", date.Format("02/01/06")),
	)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("AbsenceUseCase - ListAbsence - ctx.Done: request took too much time to be proceed")
	default:
		absences, err := a.backend.SearchAbsence(ctx, "", -1, date)
		if err != nil {
			return nil, fmt.Errorf("ListAbsence - backend.SearchAbsence: %w", err)
		}
		return absences, nil
	}
}
