package entity

import "fmt"

type Loot struct {
	ID     int
	Name   string
	Player *Player
	Raid   *Raid
}

func (l Loot) Validate() error {
	if l.Name == "" {
		return fmt.Errorf("loot name is empty")
	}
	if len(l.Name) > 20 {
		return fmt.Errorf("loot name is too long")
	}
	if l.Player == nil {
		return fmt.Errorf("loot player is empty")
	}
	if l.Raid == nil {
		return fmt.Errorf("loot raid is empty")
	}
	return nil
}
