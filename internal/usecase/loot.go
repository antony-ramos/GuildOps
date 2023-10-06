package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

type LootUseCase struct {
	backend Backend
}

func NewLootUseCase(bk Backend) *LootUseCase {
	return &LootUseCase{backend: bk}
}

func (puc LootUseCase) CreateLoot(ctx context.Context, lootName string, raidDate time.Time, playerName string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("LootUseCase - CreateLoot - ctx.Done: request took too much time to be proceed")
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", raidDate, "")
		if err != nil {
			return fmt.Errorf("CreateLoot - backend.ReadRaid: %w", err)
		}
		raid := raids[0]
		if raid.ID == 0 {
			return fmt.Errorf("raid not found")
		}

		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return fmt.Errorf("CreateLoot - backend.SearchPlayer: %w", err)
		}

		if len(player) == 0 {
			return fmt.Errorf("no player found")
		}

		loot := entity.Loot{
			Name:   lootName,
			Player: &player[0],
			Raid:   &raid,
		}
		err = loot.Validate()
		if err != nil {
			return fmt.Errorf("CreateLoot - loot.Validate: %w", err)
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
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("LootUseCase - ListLootOnPLayer - ctx.Done: request took too much time to be proceed")
	default:
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, fmt.Errorf("ListLootOnPLayer - backend.SearchPlayer: %w", err)
		}

		return player[0].Loots, nil
	}
}

func (puc LootUseCase) ListLootOnRaid(ctx context.Context, date time.Time) ([]entity.Loot, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("LootUseCase - ListLootOnPLayer - ctx.Done: request took too much time to be proceed")
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return nil, fmt.Errorf("ListLootOnPLayer - backend.SearchPlayer: %w", err)
		}

		loots, err := puc.backend.SearchLoot(ctx, "", raids[0].Date, raids[0].Difficulty)
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
	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("LootUseCase - SelectPlayerToAssign - " +
			"ctx.Done: request took too much time to be proceed")
	default:
		if len(playerNames) == 0 {
			return entity.Player{}, fmt.Errorf("player list empty")
		}

		playerList := make([]entity.Player, 0)
		for _, playerName := range playerNames {
			player, err := puc.backend.SearchPlayer(ctx, -1, playerName, "")
			if err != nil {
				return entity.Player{}, fmt.Errorf("SelectPlayerToAssign - backend.SearchPlayer: %w", err)
			}
			if len(player) == 0 {
				return entity.Player{}, fmt.Errorf("no player found")
			}
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
