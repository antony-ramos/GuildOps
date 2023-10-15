package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

type LootUseCase struct {
	backend Backend
}

func NewLootUseCase(bk Backend) *LootUseCase {
	return &LootUseCase{backend: bk}
}

func (puc LootUseCase) CreateLoot(ctx context.Context, lootName string, raidDate time.Time, playerName string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Loot/CreateLoot")
	defer span.End()
	span.SetAttributes(
		attribute.String("lootName", lootName),
		attribute.String("raidDate", raidDate.Format("02/01/2006")),
		attribute.String("playerName", playerName),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("LootUseCase - CreateLoot - ctx.Done: request took too much time to be proceed")
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", raidDate, "")
		if err != nil {
			return fmt.Errorf("CreateLoot - backend.ReadRaid: %w", err)
		}
		if len(raids) == 0 {
			return fmt.Errorf("raid not found")
		}
		raid := raids[0]
		if raid.ID == 0 {
			return fmt.Errorf("raid not found")
		}

		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("check if player exists in database: %w", err) // OK
		}

		if len(player) == 0 {
			return fmt.Errorf("no player found")
		}

		loot, err := entity.NewLoot(-1, lootName, &player[0], &raid)
		if err != nil {
			return fmt.Errorf("create a loot object: %w", err)
		}

		_, err = puc.backend.CreateLoot(ctx, loot)
		if err != nil {
			return fmt.Errorf("CreateLoot - backend.CreateLoot: %w", err)
		}

		if err != nil {
			return fmt.Errorf("CreateLoot - backend.CreateLoot: %w", err)
		}
		return nil
	}
}

func (puc LootUseCase) ListLootOnPLayer(ctx context.Context, playerName string) ([]entity.Loot, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Loot/ListLootOnPLayer")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
	)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("LootUseCase - ListLootOnPLayer - ctx.Done: request took too much time to be proceed")
	default:
		loots, err := puc.backend.SearchLoot(ctx, "", time.Time{}, "", playerName)
		if err != nil {
			return nil, fmt.Errorf("list loots on player: %w", err)
		}
		return loots, nil
	}
}

func (puc LootUseCase) ListLootOnRaid(ctx context.Context, date time.Time) ([]entity.Loot, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Loot/ListLootOnRaid")
	defer span.End()
	span.SetAttributes(
		attribute.String("date", date.Format("02/01/2006")),
	)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("LootUseCase - ListLootOnPLayer - ctx.Done: request took too much time to be proceed")
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return nil, fmt.Errorf("checking if raid exists on this date: %w", err)
		}

		if len(raids) == 0 {
			return nil, fmt.Errorf("raid not found")
		}

		loots, err := puc.backend.SearchLoot(ctx, "", raids[0].Date, raids[0].Difficulty, "")
		if err != nil {
			return nil, fmt.Errorf("ListLootOnPLayer - backend.SearchLoot: %w", err)
		}
		raids[0].Loots = loots

		return raids[0].Loots, nil
	}
}

func (puc LootUseCase) SelectPlayerToAssign(
	ctx context.Context, playerNames []string, difficulty string,
) (entity.Player, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Loot/SelectPlayerToAssign")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerNames", fmt.Sprintf("%v", playerNames)),
		attribute.String("difficulty", difficulty),
	)

	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("LootUseCase - SelectPlayerToAssign - " +
			"ctx.Done: request took too much time to be proceed")
	default:
		if len(playerNames) == 0 {
			return entity.Player{}, fmt.Errorf("player list empty")
		}

		difficulty = strings.ToLower(difficulty)
		if difficulty != "normal" && difficulty != "heroic" && difficulty != "mythic" {
			return entity.Player{}, fmt.Errorf("difficulty not valid. Must be normal, Heroic or mythic")
		}

		playerList := make([]entity.Player, 0)
		for _, playerName := range playerNames {
			player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
			if len(player) == 0 || err != nil {
				return entity.Player{}, fmt.Errorf("player %s not found", playerName)
			}

			loots, err := puc.backend.SearchLoot(ctx, "", time.Time{}, difficulty, playerName)
			if err != nil {
				return entity.Player{},
					fmt.Errorf("in a loot to check if each players given by parameters exists in database: %w", err)
			}
			player[0].Loots = loots
			playerList = append(playerList, player[0])
		}

		counter := make(map[string]int)
		for _, player := range playerList {
			counter[player.Name] = 0
			for _, loot := range player.Loots {
				if loot.Raid.Difficulty == difficulty {
					counter[player.Name]++
				}
			}
		}

		minimum := 1000
		for _, value := range counter {
			if value < minimum {
				minimum = value
			}
		}
		winningPlayers := make([]entity.Player, 0)
		for _, player := range playerList {
			if counter[player.Name] == minimum {
				winningPlayers = append(winningPlayers, player)
			}
		}
		if len(winningPlayers) > 0 {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(winningPlayers))))
			winner := n.Int64()
			return winningPlayers[winner], nil
		}

		return entity.Player{}, fmt.Errorf("no player found")
	}
}

func (puc LootUseCase) DeleteLoot(ctx context.Context, lootID int) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Loot/DeleteLoot")
	defer span.End()
	span.SetAttributes(
		attribute.Int("lootID", lootID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("LootUseCase - DeleteLoot - ctx.Done: request took too much time to be proceed")
	default:
		err := puc.backend.DeleteLoot(ctx, lootID)
		if err != nil {
			return fmt.Errorf("DeleteLoot - backend.DeleteLoot: %w", err)
		}
		return nil
	}
}
