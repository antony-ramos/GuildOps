package usecase

import (
	"context"
	"fmt"
	"time"

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
		return ctx.Err()
	default:
		player := entity.Player{
			Name: playerName,
		}
		err := player.Validate()
		if err != nil {
			return err
		}
		_, err = puc.backend.CreatePlayer(ctx, player)
		if err != nil {
			return err
		}
		return nil
	}
}

func (puc PlayerUseCase) DeletePlayer(ctx context.Context, playerName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		p, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return fmt.Errorf("player %s not found", playerName)
		}

		strikes, err := puc.backend.SearchStrike(ctx, p[0].ID, time.Time{}, "", "")
		for _, strike := range strikes {
			err = puc.backend.DeleteStrike(ctx, strike.ID)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
		err = puc.backend.DeletePlayer(ctx, p[0].ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func (puc PlayerUseCase) ReadPlayer(ctx context.Context, playerName string) (entity.Player, error) {
	select {
	case <-ctx.Done():
		return entity.Player{}, ctx.Err()
	default:
		p, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return entity.Player{}, err
		}

		if len(p) == 0 {
			return entity.Player{}, fmt.Errorf("player %s not found", playerName)
		}

		player := p[0]
		strikes, err := puc.backend.SearchStrike(ctx, player.ID, time.Time{}, "", "")
		if err != nil {
			return entity.Player{}, err
		}
		player.Strikes = strikes

		return player, nil
	}
}
