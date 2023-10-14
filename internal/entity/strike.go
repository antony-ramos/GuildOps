package entity

import (
	"fmt"
	"time"
)

type Strike struct {
	Date   time.Time
	ID     int
	Season string
	Reason string

	Player *Player
}

func SeasonCalculator(date time.Time) string {
	if date.After(time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)) &&
		date.Before(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return "DF/S2"
	} else {
		return "Unknown"
	}
}

func NewStrike(reason string) (Strike, error) {
	if len(reason) == 0 {
		return Strike{}, fmt.Errorf("reason must not be empty")
	}
	if len(reason) > 255 {
		return Strike{}, fmt.Errorf("reason must not be longer than 255 characters")
	}

	return Strike{
		Reason: reason,
		Date:   time.Now(),
		Season: SeasonCalculator(time.Now()),
	}, nil
}
