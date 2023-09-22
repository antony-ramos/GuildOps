package entity

import "fmt"

type Absence struct {
	ID     int
	Player *Player
	Raid   *Raid
}

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (a Absence) Validate() error {
	// check if absence is valid
	if a.Player == nil {
		return Error{Message: "player cannot be nil"}
	}
	if a.Raid == nil {
		return Error{Message: "raid cannot be nil"}
	}
	err := a.Player.Validate()
	if err != nil {
		return Error{Message: fmt.Sprintf("Absence Invalid : %s", err)}
	}
	err = a.Raid.Validate()
	if err != nil {
		return Error{Message: fmt.Sprintf("Absence Invalid : %s", err)}
	}
	return nil
}
