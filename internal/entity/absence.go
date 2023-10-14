package entity

import "fmt"

type Absence struct {
	ID     int
	Player *Player
	Raid   *Raid
}

func NewAbsence(index int, player *Player, raid *Raid) (Absence, error) {
	if player == nil {
		return Absence{}, fmt.Errorf("player cannot be nil")
	}
	if raid == nil {
		return Absence{}, fmt.Errorf("raid cannot be nil")
	}

	return Absence{
		ID:     index,
		Player: player,
		Raid:   raid,
	}, nil
}
