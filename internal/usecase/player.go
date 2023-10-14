package usecase

import (
	"context"
	"fmt"
	"strings"
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
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Player/CreatePlayer")
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
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Player/DeletePlayer")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - DeletePlayer - ctx.Done: request took too much time to be proceed")
	default:
		playerName := strings.ToLower(playerName)
		err := puc.backend.DeletePlayer(ctx, entity.Player{Name: playerName})
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.DeletePlayer: %w", err)
		}
		return nil
	}
}

func (puc PlayerUseCase) ReadPlayer(ctx context.Context, playerName, playerLinkName string) (entity.Player, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Player/ReadPlayer")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.String("playerLinkName", playerLinkName),
	)
	logger.FromContext(ctx).Debug("read player use case")

	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("PlayerUseCase - ReadPlayer - ctx.Done: request took too much time to be proceed")
	default:
		playerName := strings.ToLower(playerName)
		playerLinkName := strings.ToLower(playerLinkName)

		player := entity.Player{
			ID: -1,
		}
		if playerName != "" {
			plrs, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
			if err != nil {
				return entity.Player{}, fmt.Errorf("get basic info for player: %w", err)
			}
			if len(plrs) == 0 {
				return entity.Player{}, fmt.Errorf("player %s not found", playerName)
			}
			player = plrs[0]
		} else if playerLinkName != "" {
			plrs, err := puc.backend.SearchPlayer(ctx, -1, "", playerLinkName)
			if err != nil {
				return entity.Player{}, fmt.Errorf("get basic info for player: %w", err)
			}
			if len(plrs) == 0 {
				return entity.Player{}, fmt.Errorf("didn't find a player linked to this discord user named %s", playerLinkName)
			}
			player = plrs[0]
		}

		if player.ID == -1 {
			return entity.Player{}, fmt.Errorf("player %s not found", playerName)
		}

		strikes, err := puc.backend.SearchStrike(ctx, player.ID, time.Time{}, "", "")
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchStrike: %w", err)
		}
		player.Strikes = strikes

		fails, err := puc.backend.SearchFail(ctx, "", player.ID, -1, "")
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchFail: %w", err)
		}
		player.Fails = fails

		loots, err := puc.backend.SearchLoot(ctx, "", time.Time{}, "", player.Name)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchLoot: %w", err)
		}
		player.Loots = loots

		missedRaids, err := puc.backend.SearchAbsence(ctx, "", player.ID, time.Time{})
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.SearchAbsence: %w", err)
		}
		for _, missedRaid := range missedRaids {
			player.MissedRaids = append(player.MissedRaids, *missedRaid.Raid)
		}

		for k, fail := range fails {
			r, err := puc.backend.ReadRaid(ctx, fail.Raid.ID)
			if err != nil {
				return entity.Player{}, errors.Wrap(err, "read player, for each fail, read raid")
			}
			fails[k].Raid = &r
		}

		return player, nil
	}
}

func (puc PlayerUseCase) LinkPlayer(ctx context.Context, playerName string, discordID string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Player/LinkPlayer")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.String("discordID", discordID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("PlayerUseCase - LinkPlayer - ctx.Done: request took too much time to be proceed")
	default:
		playerName := strings.ToLower(playerName)
		discordID := strings.ToLower(discordID)

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
