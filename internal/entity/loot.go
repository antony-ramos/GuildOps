package entity

type Loot struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Raid     *Raid  `json:"raid"`
	PlayerID int    `json:"player_id"`
}

func (l Loot) Validate() error {
	return nil
}
