package postgresbackend_test

import (
	"context"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPG_CreateLoot(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		loot := entity.Loot{
			ID:   1,
			Name: "lootname",
			Raid: &entity.Raid{
				ID:         1,
				Name:       "raidname",
				Difficulty: "difficulty",
				Date:       time.Now(),
			},
			Player: &entity.Player{
				ID:   1,
				Name: "playername",
			},
		}

		columns := []string{"name", "raid_id", "player_id"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT name, raid_id, player_id FROM loots WHERE name = $1 AND raid_id = $2 AND player_id = $3",
			loot.Name, loot.Raid.ID, loot.Player.ID).
			Return(pgxRows, nil)

		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO loots (name,raid_id,player_id) VALUES ($1,$2,$3)",
			loot.Name, loot.Raid.ID, loot.Player.ID).
			Return(nil, nil)

		loot, err := pgBackend.CreateLoot(context.Background(), loot)
		assert.NoError(t, err)
		assert.Equal(t, loot, loot)
	})
}

func TestPG_SearchLoot(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		loot := entity.Loot{
			ID:   1,
			Name: "lootname",
			Raid: &entity.Raid{
				ID:         1,
				Name:       "raidname",
				Difficulty: "difficulty",
				Date:       time.Now(),
			},
			Player: &entity.Player{
				ID:   1,
				Name: "playername",
			},
		}

		columns := []string{
			"loots.id", "loots.name", "loots.raid_id",
			"raids.name", "raids.difficulty", "raids.date",
			"loots.player_id", "players.name",
		}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(loot.ID, loot.Name, loot.Raid.ID,
			loot.Raid.Name, loot.Raid.Difficulty, loot.Raid.Date,
			loot.Player.ID, loot.Player.Name).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT loots.id, loots.name, loots.raid_id, "+
				"raids.name, raids.difficulty, raids.date, loots.player_id, p"+
				"layers.name FROM loots JOIN raids ON raids.id = loots.raid_id "+
				"JOIN players ON players.id = loots.player_id "+
				"WHERE raids.date = $1 AND raids.difficulty = $2",
			loot.Raid.Date, loot.Raid.Difficulty).
			Return(pgxRows, nil)

		loots, err := pgBackend.SearchLoot(context.Background(), loot.Name, loot.Raid.Date, loot.Raid.Difficulty)
		assert.NoError(t, err)
		assert.Equal(t, loots, []entity.Loot{loot})
	})
}

func TestPG_UpdateLoot(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		loot := entity.Loot{
			ID:   1,
			Name: "lootname",
			Raid: &entity.Raid{
				ID:         1,
				Name:       "raidname",
				Difficulty: "difficulty",
				Date:       time.Now(),
			},
			Player: &entity.Player{
				ID:   1,
				Name: "playername",
			},
		}

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE loots SET name = $1, raid_id = $2, player_id = $3 WHERE id = $4",
			gomock.Any(), gomock.Any()).
			Return(nil, nil)

		err := pgBackend.UpdateLoot(context.Background(), loot)
		assert.NoError(t, err)
	})
}
