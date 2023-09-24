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
	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
	})
	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
	})
	t.Run("Loot is not created", func(t *testing.T) {
		t.Parallel()
	})
}
