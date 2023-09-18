package entity

type Loot struct {
	ID     int
	Name   string
	Player *Player
	Raid   *Raid
}

func (l Loot) Validate() error {
	return nil
}
