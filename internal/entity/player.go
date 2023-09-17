package entity

import (
	"fmt"
	"strings"
)

//TODO Tests

type Player struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Strikes     []Strike `json:"strikes"`
	Loots       []Loot   `json:"loots"`
	MissedRaids []Raid   `json:"missed_raids"`
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
