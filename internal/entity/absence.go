package entity

import (
	"time"
)

type Absence struct {
	ID     int     `json:"id"`
	Player *Player `json:"player"`
	Raid   *Raid   `json:"raid"`
}

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrorNameCannotBeEmpty       = Error{Message: "name of player cannot be empty"}
	ErrorDateCannotBeBeforeToday = Error{Message: "date of absence cannot be before today"}
	ErrorDateIsNotMonWesThu      = Error{Message: "date must be monday, wednesday or thursday"}
)

func validateDate(dateToCheck time.Time) error {
	today := time.Now()
	if dateToCheck.Before(today) {
		return ErrorDateCannotBeBeforeToday
	}

	dayOfWeek := dateToCheck.Weekday()
	if dayOfWeek != time.Monday && dayOfWeek != time.Wednesday && dayOfWeek != time.Thursday {
		return ErrorDateIsNotMonWesThu
	}

	return nil
}

func (a Absence) Validate() error {
	return nil
}
