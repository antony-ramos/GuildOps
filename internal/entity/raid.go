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

	Loots []*Loot
}

func (r *Raid) Validate() error {
	if len(r.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	r.Difficulty = strings.ToLower(r.Difficulty)

	if r.Difficulty != "normal" && r.Difficulty != "heroic" && r.Difficulty != "mythic" {
		return fmt.Errorf("difficulty must be normal, heroic, or mythic")
	}

	// r.Name must be between 1 and 12 characters, only a-z
	if len(r.Name) < 1 || len(r.Name) > 12 {
		return fmt.Errorf("name must be between 1 and 12 characters")
	}

	for _, char := range r.Name {
		if !unicode.IsLetter(char) {
			return fmt.Errorf("name must only contain letters")
		}
	}

	r.Name = strings.ToLower(r.Name)

	if r.Name != strings.ToLower(r.Name) {
		return fmt.Errorf("name must be lowercase")
	}
	return nil
}
