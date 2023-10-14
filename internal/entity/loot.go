package entity

import (
	"fmt"
	"strings"
)

type Loot struct {
	ID     int
	Name   string
	Player *Player
	Raid   *Raid
}

func NewLoot(index int, name string, player *Player, raid *Raid) (Loot, error) {
	if len(name) == 0 {
		return Loot{}, fmt.Errorf("name cannot be empty")
	}
	name = strings.ToLower(name)
	if len(name) < 1 || len(name) > 30 {
		return Loot{}, fmt.Errorf("name must be between 1 and 30 characters")
	}

	if player == nil {
		return Loot{}, fmt.Errorf("player cannot be nil")
	}

	if raid == nil {
		return Loot{}, fmt.Errorf("raid cannot be nil")
	}

	return Loot{
		ID:     index,
		Name:   name,
		Player: player,
		Raid:   raid,
	}, nil
}
