package usecase

import (
	"context"
	"time"

	"github.com/coven-discord-bot/internal/entity"
)

type Backend interface {
	Player
	Strike
	Raid
	Loot
	Absence
}

type Player interface {
	SearchPlayer(ctx context.Context, id int, name string) ([]entity.Player, error)
	CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error)
	ReadPlayer(ctx context.Context, playerID int) (entity.Player, error)
	UpdatePlayer(ctx context.Context, player entity.Player) error
	DeletePlayer(ctx context.Context, playerID int) error
}

type Strike interface {
	SearchStrike(ctx context.Context, playerID int, Date time.Time, Season, Reason string) ([]entity.Strike, error)
	CreateStrike(ctx context.Context, strike entity.Strike, player entity.Player) error
	ReadStrike(ctx context.Context, strikeID int) (entity.Strike, error)
	UpdateStrike(ctx context.Context, strike entity.Strike) error
	DeleteStrike(ctx context.Context, strikeID int) error
}

type Raid interface {
	SearchRaid(ctx context.Context, raidName string, date time.Time, difficulty string) ([]entity.Raid, error)
	CreateRaid(ctx context.Context, raid entity.Raid) (entity.Raid, error)
	ReadRaid(ctx context.Context, raidID int) (entity.Raid, error)
	UpdateRaid(ctx context.Context, raid entity.Raid) error
	DeleteRaid(ctx context.Context, raidID int) error
}

type Loot interface {
	SearchLoot(ctx context.Context, name string, date time.Time, difficulty string) ([]entity.Loot, error)
	CreateLoot(ctx context.Context, loot entity.Loot) (entity.Loot, error)
	ReadLoot(ctx context.Context, lootID int) (entity.Loot, error)
	UpdateLoot(ctx context.Context, loot entity.Loot) error
	DeleteLoot(ctx context.Context, lootID int) error
}

type Absence interface {
	SearchAbsence(ctx context.Context, playerName string, playerID int, date time.Time) ([]entity.Absence, error)
	CreateAbsence(ctx context.Context, absence entity.Absence) (entity.Absence, error)
	ReadAbsence(ctx context.Context, absenceID int) (entity.Absence, error)
	UpdateAbsence(ctx context.Context, absence entity.Absence) error
	DeleteAbsence(ctx context.Context, absenceID int) error
}
