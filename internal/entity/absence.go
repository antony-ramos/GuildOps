package entity

type Absence struct {
	ID     int
	Player *Player
	Raid   *Raid
}

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (a Absence) Validate() error {
	return nil
}
