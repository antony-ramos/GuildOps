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

func NewPlayer(index int, name, discordName string) (Player, error) {
	if len(name) == 0 {
		return Player{}, fmt.Errorf("name cannot be empty")
	}

	name = strings.ToLower(name)

	if len(name) < 1 || len(name) > 12 {
		return Player{}, fmt.Errorf("name must be between 1 and 12 characters")
	}

	for _, char := range name {
		if !unicode.IsLetter(char) {
			return Player{}, fmt.Errorf("name must only contain letters")
		}
	}

	return Player{
		ID:          index,
		Name:        name,
		DiscordName: discordName,
	}, nil
}
