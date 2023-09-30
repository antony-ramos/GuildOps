package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

type PlayerUseCase struct {
	backend Backend
}

func NewPlayerUseCase(bk Backend) *PlayerUseCase {
	return &PlayerUseCase{backend: bk}
}

func (puc PlayerUseCase) CreatePlayer(ctx context.Context, playerName string) (int, error) {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "PlayerUseCase/CreatePlayer")
	span.SetAttributes(attribute.String("playerName", playerName))
	defer span.End()
	select {
	case <-ctx.Done():
		return -1, fmt.Errorf("PlayerUseCase - CreatePlayer - ctx.Done: request took too much time to be proceed")
	default:
		player := entity.Player{
			Name: playerName,
		}

		_, spanValidate := otel.Tracer("Entity").Start(ctx, "Player/Validate")
		span.SetAttributes(attribute.String("playerName", playerName))
		err := player.Validate()
		spanValidate.End()

		if err != nil {
			return -1, fmt.Errorf("database - CreatePlayer - r.Validate: %w", err)
		}
		player, err = puc.backend.CreatePlayer(ctx, player)
		if err != nil {
			return -1, fmt.Errorf("database - CreatePlayer - r.CreatePlayer: %w", err)
		}
		return player.ID, nil
	}
}

func (puc PlayerUseCase) DeletePlayer(ctx context.Context, playerName string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - DeletePlayer - ctx.Done: request took too much time to be proceed")
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("player %s not found", playerName)
		}

		strikes, err := puc.backend.SearchStrike(ctx, player[0].ID, time.Time{}, "", "")
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.SearchStrike: %w", err)
		}
		for _, strike := range strikes {
			err = puc.backend.DeleteStrike(ctx, strike.ID)
			if err != nil {
				return fmt.Errorf("database - DeletePlayer - r.DeleteStrike: %w", err)
			}
		}

		fails, err := puc.backend.SearchFail(ctx, "", player[0].ID, -1, "")
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.SearchFail: %w", err)
		}
		for _, fail := range fails {
			err = puc.backend.DeleteFail(ctx, fail.ID)
			if err != nil {
				return fmt.Errorf("database - DeletePlayer - r.DeleteFail: %w", err)
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
	ctx, span := otel.Tracer("UseCase").Start(ctx, "PlayerUseCase/ReadPlayer")
	span.SetAttributes(attribute.String("playerName", playerName))
	defer span.End()
	logger.FromContext(ctx).Debug("read player use case")

	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("PlayerUseCase - ReadPlayer - ctx.Done: request took too much time to be proceed")
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
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

		fails, err := puc.backend.SearchFail(ctx, "", player[0].ID, -1, "")
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchFail: %w", err)
		}
		player[0].Fails = fails

		for k, fail := range fails {
			r, err := puc.backend.ReadRaid(ctx, fail.Raid.ID)
			if err != nil {
				return entity.Player{}, errors.Wrap(err, "read player, for each fail, read raid")
			}
			fails[k].Raid = &r
		}

		return player[0], nil
	}
}

func (puc PlayerUseCase) LinkPlayer(ctx context.Context, playerName string, discordID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - LinkPlayer - ctx.Done: request took too much time to be proceed")
	default:
		alreadyLinked, err := puc.backend.SearchPlayer(ctx, -1, "", discordID)
		if err != nil {
			return fmt.Errorf("database - LinkPlayer - r.SearchPlayer: %w", err)
		}
		if len(alreadyLinked) > 0 {
			return fmt.Errorf("discord account already linked to player name %s. Contact Staff for modification",
				alreadyLinked[0].Name)
		}

		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("database - LinkPlayer - r.SearchPlayer: %w", err)
		}
		if len(player) == 0 {
			return fmt.Errorf("player %s not found", playerName)
		}
		player[0].DiscordName = discordID
		err = puc.backend.UpdatePlayer(ctx, player[0])
		if err != nil {
			return fmt.Errorf("database - LinkPlayer - r.UpdatePlayer: %w", err)
		}
		return nil
	}
}
