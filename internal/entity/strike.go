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

func (s Strike) Validate() error {
	if len(s.Reason) == 0 {
		return fmt.Errorf("reason must not be empty")
	}
	if len(s.Reason) > 100 {
		return fmt.Errorf("reason must not be longer than 255 characters")
	}

	if s.Date.After(time.Now()) {
		return nil
	}
	return nil
}
