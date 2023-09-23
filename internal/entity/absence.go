package entity

import "fmt"

type Absence struct {
	ID     int
	Player *Player
	Raid   *Raid
}

func (a Absence) Validate() error {
	// check if absence is valid
	if a.Player == nil {
		return fmt.Errorf("player cannot be nil")
	}
	if a.Raid == nil {
		return fmt.Errorf("raid cannot be nil")
	}
	err := a.Player.Validate()
	if err != nil {
		return fmt.Errorf("absence Validate error : %w", err)
	}
	err = a.Raid.Validate()
	if err != nil {
		return fmt.Errorf("absence Validate error : %w", err)
	}
	return nil
}
