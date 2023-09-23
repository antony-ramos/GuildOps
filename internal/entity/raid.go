package entity

import (
	"fmt"
	"strings"
	"time"
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

func (r Raid) Validate() error {
	if len(r.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	if r.Difficulty != "normal" && r.Difficulty != "heroic" && r.Difficulty != "mythic" {
		return fmt.Errorf("difficulty must be normal, heroic, or mythic")
	}

	// r.Name must be between 1 and 12 characters, only a-z
	if len(r.Name) < 1 || len(r.Name) > 12 {
		return fmt.Errorf("name must be between 1 and 12 characters")
	}
	// only a-z characters or spaces
	for _, c := range r.Name {
		if c != ' ' && (c < 'a' || c > 'z') {
			return fmt.Errorf("name must only contain a-z characters or spaces")
		}
	}
	if r.Name != strings.ToLower(r.Name) {
		return fmt.Errorf("name must be lowercase")
	}
	return nil
}
