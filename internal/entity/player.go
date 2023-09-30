package entity

import (
	"fmt"
	"strings"
	"unicode"
)

type Player struct {
	ID          int
	Name        string
	DiscordName string

	Strikes     []Strike
	Loots       []Loot
	MissedRaids []Raid
	Fails       []Fail
}

func (p *Player) Validate() error {
	if len(p.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}

	// p.Name must be between 1 and 12 characters, only a-z
	if len(p.Name) < 1 || len(p.Name) > 12 {
		return fmt.Errorf("name must be between 1 and 12 characters")
	}
	// only a-z characters
	for _, char := range p.Name {
		if !unicode.IsLetter(char) {
			return fmt.Errorf("name must only contain letters")
		}
	}

	p.Name = strings.ToLower(p.Name)

	return nil
}
