package entity

import "fmt"

type Fail struct {
	ID     int
	Reason string
	Player *Player
	Raid   *Raid
}

func (f Fail) Validate() error {
	if f.Player == nil {
		return fmt.Errorf("player cannot be nil")
	}
	if f.Raid == nil {
		return fmt.Errorf("raid cannot be nil")
	}
	err := f.Player.Validate()
	if err != nil {
		return fmt.Errorf("validate absence player : %w", err)
	}
	err = f.Raid.Validate()
	if err != nil {
		return fmt.Errorf("validate absence raid : %w", err)
	}
	if len(f.Reason) == 0 {
		return fmt.Errorf("reason cannot be empty")
	}
	return nil
}
