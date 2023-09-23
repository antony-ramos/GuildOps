package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

type PlayerUseCase struct {
	backend Backend
}

func NewPlayerUseCase(bk Backend) *PlayerUseCase {
	return &PlayerUseCase{backend: bk}
}

func (puc PlayerUseCase) CreatePlayer(ctx context.Context, playerName string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - CreatePlayer - ctx.Done: %w", ctx.Err())
	default:
		player := entity.Player{
			Name: playerName,
		}
		err := player.Validate()
		if err != nil {
			return fmt.Errorf("database - CreatePlayer - r.Validate: %w", err)
		}
		_, err = puc.backend.CreatePlayer(ctx, player)
		if err != nil {
			return fmt.Errorf("database - CreatePlayer - r.CreatePlayer: %w", err)
		}
		return nil
	}
}

func (puc PlayerUseCase) DeletePlayer(ctx context.Context, playerName string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - DeletePlayer - ctx.Done: %w", ctx.Err())
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("player %s not found", playerName)
		}

		strikes, err := puc.backend.SearchStrike(ctx, player[0].ID, time.Time{}, "", "")
		for _, strike := range strikes {
			err = puc.backend.DeleteStrike(ctx, strike.ID)
			if err != nil {
				return fmt.Errorf("database - DeletePlayer - r.DeleteStrike: %w", err)
			}
		}

		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.SearchStrike: %w", err)
		}
		err = puc.backend.DeletePlayer(ctx, player[0].ID)
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.DeletePlayer: %w", err)
		}
		return nil
	}
}

func (puc PlayerUseCase) ReadPlayer(ctx context.Context, playerName string) (entity.Player, error) {
	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("PlayerUseCase - ReadPlayer - ctx.Done: %w", ctx.Err())
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchPlayer: %w", err)
		}

		if len(player) == 0 {
			return entity.Player{}, fmt.Errorf("player %s not found", playerName)
		}

		strikes, err := puc.backend.SearchStrike(ctx, player[0].ID, time.Time{}, "", "")
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchStrike: %w", err)
		}
		player[0].Strikes = strikes

		return player[0], nil
	}
}
