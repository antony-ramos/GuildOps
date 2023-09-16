package entity

import (
	"fmt"
	"strings"
)

//TODO Tests

type Player struct {
	ID   int
	Name string

	Strikes     []Strike
	Loots       []Loot
	MissedRaids []Raid
}

func (p Player) Validate() error {
	if len(p.Name) == 0 {
		return ErrorNameCannotBeEmpty
	}

	if p.Name != strings.ToLower(p.Name) {
		return fmt.Errorf("name must be lowercase")
	}
	return nil
}
