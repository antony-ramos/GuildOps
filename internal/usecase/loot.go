package usecase

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"math/rand"
)

type LootUseCase struct {
	backend Backend
}

func NewLootUseCase(bk Backend) *LootUseCase {
	return &LootUseCase{backend: bk}
}

func (puc LootUseCase) CreateLoot(ctx context.Context, lootName string, raidID int, playerName string) error {
	raid, err := puc.backend.ReadRaid(ctx, raidID)
	if err != nil {
		return err
	}

	player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return err
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
		return err
	}

	loot, err = puc.backend.CreateLoot(ctx, loot)
	if err != nil {
		return err
	}

	if err != nil {
		return err

	}
	return nil
}

func (puc LootUseCase) ListLootOnPLayer(ctx context.Context, playerName string) ([]entity.Loot, error) {
	player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
	if err != nil {
		return nil, err
	}

	return player[0].Loots, nil
}

func (puc LootUseCase) SelectPlayerToAssign(ctx context.Context, playerNames []string, difficulty string) (entity.Player, error) {
	playerList := make([]entity.Player, 0)
	for _, playerName := range playerNames {
		player, err := puc.backend.SearchPlayer(ctx, -1, playerName)
		if err != nil {
			return entity.Player{}, err
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
	
	min := 1000
	for _, value := range counter {
		if value < min {
			min = value
		}
	}
	winningPlayers := make([]entity.Player, 0)
	for _, player := range playerList {
		if counter[player.Name] == min {
			winningPlayers = append(winningPlayers, player)
		}
	}
	if len(winningPlayers) > 0 {
		return winningPlayers[rand.Intn(len(winningPlayers))], nil
	}

	return entity.Player{}, fmt.Errorf("no player found")
}

func (puc LootUseCase) DeleteLoot(ctx context.Context, lootID int) error {
	err := puc.backend.DeleteLoot(ctx, lootID)
	if err != nil {
		return err

	}
	return nil
}
