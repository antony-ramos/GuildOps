package entity

import (
	"fmt"
	"strings"
)

type Player struct {
	ID          int
	Name        string
	DiscordName string

	Strikes     []Strike
	Loots       []Loot
	MissedRaids []Raid
}

func (p Player) Validate() error {
	if len(p.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	// p.Name must be between 1 and 12 characters, only a-z
	if len(p.Name) < 1 || len(p.Name) > 12 {
		return fmt.Errorf("name must be between 1 and 12 characters")
	}
	// only a-z characters
	for _, c := range p.Name {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("name must only contain a-z characters")
		}
	}

	if p.Name != strings.ToLower(p.Name) {
		return fmt.Errorf("name must be lowercase")
	}
	return nil
}
