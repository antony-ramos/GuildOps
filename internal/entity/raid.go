package entity

import (
	"time"
)

//TODO Tests

// TODO Read Difficulties from backend

type Raid struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	Difficulty string    `json:"difficulty"`

	Absences []*Player `json:"absences"`
	Players  []*Player `json:"players"`
	Bench    []*Player `json:"bench"`

	Loots []*Loot `json:"loots"`
}

func (r Raid) Validate() error {

	return nil
}
