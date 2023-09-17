package entity

import (
	"fmt"
	"time"
)

//TODO Tests

type Strike struct {
	Date   time.Time `json:"date"`
	ID     int       `json:"id"`
	Season string    `json:"season"`
	Reason string    `json:"reason"`
	Player *Player   `json:"player"`
}

var seasons = []string{"s1", "s2", "s3"}

func (s Strike) Validate() error {
	if len(s.Reason) == 0 {
		return fmt.Errorf("reason must not be empty")
	}
	if len(s.Reason) > 100 {
		return fmt.Errorf("reason must not be longer than 255 characters")
	}

	isValidSeason := false
	for _, ss := range seasons {
		if s.Season == ss {
			isValidSeason = true
			break
		}
	}
	if !isValidSeason {
		return fmt.Errorf("season must be : %s", seasons)
	}

	if s.Date.After(time.Now()) {
		return nil
	}
	return nil
}
