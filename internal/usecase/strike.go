package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

type StrikeUseCase struct {
	backend Backend
}

// NewStrikeUseCase is a StrikeUseCase Object generator.
func NewStrikeUseCase(bk Backend) *StrikeUseCase {
	return &StrikeUseCase{backend: bk}
}

func SeasonCalculator(date time.Time) string {
	if date.After(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)) &&
		date.Before(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return "DF/S2"
	} else {
		return "Unknown"
	}
}

// CreateStrike is a function which call backend to Create a Strike Object.
func (puc StrikeUseCase) CreateStrike(ctx context.Context, strikeReason, playerName string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("StrikeUseCase - CreateStrike - ctx.Done: %w", ctx.Err())
	default:
		strike := entity.Strike{
			Reason: strikeReason,
			Date:   time.Now(),
			Season: SeasonCalculator(time.Now()),
		}
		err := strike.Validate()
		if err != nil {
			return fmt.Errorf("database - CreateStrike - r.Validate: %w", err)
		}

		player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return fmt.Errorf("database - CreateStrike - r.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return errors.New("player not found")
		}
		err = puc.backend.CreateStrike(ctx, strike, player[0])
		if err != nil {
			return fmt.Errorf("database - CreateStrike - r.CreateStrike: %w", err)
		}
		return nil
	}
}

// DeleteStrike is a function which call backend to Delete a Strike Object.
func (puc StrikeUseCase) DeleteStrike(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("StrikeUseCase - DeleteStrike - ctx.Done: %w", ctx.Err())
	default:
		err := puc.backend.DeleteStrike(ctx, id)
		if err != nil {
			return fmt.Errorf("database - DeleteStrike - r.DeleteStrike: %w", err)
		}
		return nil
	}
}

// ReadStrikes is a function which call backend to Read all strikes on a player.
func (puc StrikeUseCase) ReadStrikes(ctx context.Context, playerName string) ([]entity.Strike, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("StrikeUseCase - ReadStrikes - ctx.Done: %w", ctx.Err())
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return nil, fmt.Errorf("database - ReadStrikes - r.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return nil, errors.New("player not found")
		}
		strikes, err := puc.backend.SearchStrike(ctx, player[0].ID, time.Time{}, "", "")
		if err != nil {
			return nil, fmt.Errorf("database - ReadStrikes - r.SearchStrike: %w", err)
		}
		if len(strikes) == 0 {
			return nil, errors.New("no strikes found")
		}

		return strikes, nil
	}
}
