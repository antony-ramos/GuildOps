package entity

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

type Raid struct {
	ID         int
	Name       string
	Date       time.Time
	Difficulty string

	Absences []*Player
	Players  []*Player
	Bench    []*Player

	Loots []Loot
}

func NewRaid(name, difficulty string, date time.Time) (Raid, error) {
	if len(name) == 0 {
		return Raid{}, fmt.Errorf("name cannot be empty")
	}

	difficulty = strings.ToLower(difficulty)
	name = strings.ToLower(name)

	if difficulty != "normal" && difficulty != "heroic" && difficulty != "mythic" {
		return Raid{}, fmt.Errorf("difficulty must be normal, heroic, or mythic")
	}

	if len(name) < 1 || len(name) > 12 {
		return Raid{}, fmt.Errorf("name must be between 1 and 12 characters")
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return Raid{}, fmt.Errorf("name must only contain letters")
		}
	}

	return Raid{
		Name:       name,
		Difficulty: difficulty,
		Date:       date,
	}, nil
}
