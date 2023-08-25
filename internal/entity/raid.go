package entity

import "time"

type Raid struct {
	date time.Time
}

func (r Raid) Validate() error {
	//TODO
	return nil
}
