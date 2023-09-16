package entity

import (
	"time"
)

//TODO Tests

// TODO Read Difficulties from backend

type Raid struct {
	ID         int
	Name       string
	Date       time.Time
	Difficulty string

	Absences []*Player
	Players  []*Player
	Bench    []*Player

	Loots []*Loot
}

func (r Raid) Validate() error {

	return nil
}
