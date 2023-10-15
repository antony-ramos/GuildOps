package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

type StrikeUseCase struct {
	backend Backend
}

// NewStrikeUseCase is a StrikeUseCase Object generator.
func NewStrikeUseCase(bk Backend) *StrikeUseCase {
	return &StrikeUseCase{backend: bk}
}

// CreateStrike is a function which call backend to Create a Strike Object.
func (puc StrikeUseCase) CreateStrike(ctx context.Context, strikeReason, playerName string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Strike/CreateStrike")
	defer span.End()
	span.SetAttributes(
		attribute.String("strikeReason", strikeReason),
		attribute.String("playerName", playerName),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("StrikeUseCase - CreateStrike - ctx.Done: request took too much time to be proceed")
	default:
		playerName := strings.ToLower(playerName)
		strike, err := entity.NewStrike(strikeReason)
		if err != nil {
			return fmt.Errorf("create an object strike for backend: %w", err)
		}

		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("get basic info for player: %w", err)
		}
		if len(player) == 0 {
			return errors.New("player not found")
		}
		err = puc.backend.CreateStrike(ctx, strike, player[0].ID)
		if err != nil {
			return fmt.Errorf("database - CreateStrike - r.CreateStrike: %w", err)
		}
		return nil
	}
}

// DeleteStrike is a function which call backend to Delete a Strike Object.
func (puc StrikeUseCase) DeleteStrike(ctx context.Context, strikeID int) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Strike/DeleteStrike")
	defer span.End()
	span.SetAttributes(
		attribute.Int("strikeID", strikeID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("StrikeUseCase - DeleteStrike - ctx.Done: request took too much time to be proceed")
	default:
		err := puc.backend.DeleteStrike(ctx, strikeID)
		if err != nil {
			return fmt.Errorf("database DeleteStrike: r.DeleteStrike: %w", err)
		}
		return nil
	}
}

// ReadStrikes is a function which call backend to Read all strikes on a player.
func (puc StrikeUseCase) ReadStrikes(ctx context.Context, playerName string) ([]entity.Strike, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Strike/ReadStrikes")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("StrikeUseCase - ReadStrikes - ctx.Done: request took too much time to be proceed")
	default:
		playerName := strings.ToLower(playerName)
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, fmt.Errorf("get basic info about player: %w", err)
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
