package entity

import "fmt"

type Fail struct {
	ID     int
	Reason string
	Player *Player
	Raid   *Raid
}

func NewFail(index int, reason string, player *Player, raid *Raid) (Fail, error) {
	if len(reason) == 0 {
		return Fail{}, fmt.Errorf("reason cannot be empty")
	}
	if player == nil {
		return Fail{}, fmt.Errorf("player cannot be nil")
	}
	if raid == nil {
		return Fail{}, fmt.Errorf("raid cannot be nil")
	}

	return Fail{
		ID:     index,
		Reason: reason,
		Player: player,
		Raid:   raid,
	}, nil
}
