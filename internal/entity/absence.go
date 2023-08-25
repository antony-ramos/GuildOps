package entity

import "time"

type Absence struct {
	Name string
	Date time.Time
}

func (a Absence) Validate() error {
	//TODO
	return nil
}
