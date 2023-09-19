package usecase

import (
	"context"
	"errors"
	"github.com/antony-ramos/guildops/internal/entity"
	"time"
)

type StrikeUseCase struct {
	backend Backend
}

// NewStrikeUseCase is a StrikeUseCase Object generator
func NewStrikeUseCase(bk Backend) *StrikeUseCase {
	return &StrikeUseCase{backend: bk}
}

func SeasonCalculator(date time.Time) string {
	if date.After(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)) && date.Before(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return "DF/S2"
	} else {
		return "Unknown"
	}
}

// CreateStrike is a function which call backend to Create a Strike Object
func (puc StrikeUseCase) CreateStrike(ctx context.Context, strikeReason, playerName string) error {
	strike := entity.Strike{
		Reason: strikeReason,
		Date:   time.Now(),
		Season: SeasonCalculator(time.Now()),
	}
	err := strike.Validate()
	if err != nil {
		return err
	}

	p, err := puc.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return err
	}
	if len(p) == 0 {
		return errors.New("player not found")
	}
	err = puc.backend.CreateStrike(ctx, strike, p[0])
	if err != nil {
		return err
	}
	return nil
}

// DeleteStrike is a function which call backend to Delete a Strike Object
func (puc StrikeUseCase) DeleteStrike(ctx context.Context, ID int) error {
	err := puc.backend.DeleteStrike(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

// ReadStrikes is a function which call backend to Read all strikes on a player
func (puc StrikeUseCase) ReadStrikes(ctx context.Context, playerName string) ([]entity.Strike, error) {
	p, err := puc.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return nil, err

	}
	if len(p) == 0 {
		return nil, errors.New("player not found")
	}
	strikes, err := puc.backend.SearchStrike(ctx, p[0].ID, time.Time{}, "", "")
	if err != nil {
		return nil, err

	}
	if len(strikes) == 0 {
		return nil, errors.New("no strikes found")
	}

	return strikes, nil
}
