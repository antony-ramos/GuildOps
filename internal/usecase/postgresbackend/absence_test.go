package postgresbackend_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPG_SearchAbsence(t *testing.T) {
	t.Parallel()

	t.Run("Searching on playerName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		columns := []string{"id", "player_id", "raid_id", "name", "difficulty", "date", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, 0, "test", "test", time.Now(), "test").ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT absences.id, absences.player_id, absences.raid_id, "+
				"raids.name, raids.difficulty, raids.date, players.name "+
				"FROM absences "+
				"JOIN raids ON raids.id = absences.raid_id "+
				"JOIN players ON players.id = absences.player_id "+
				"WHERE players.name = $1", "test").
			Return(pgxRows, nil)

		absence, err := pgBackend.SearchAbsence(context.Background(), "test", -1, time.Now())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(absence))
	})
}

func TestPG_CreateAbsence(t *testing.T) {
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

		abs := entity.Absence{
			ID: 0,
			Player: &entity.Player{
				ID:   0,
				Name: "test",
			},
			Raid: &entity.Raid{
				ID:   0,
				Name: "test",
				Date: time.Now(),
			},
		}

		columns := []string{"id", "player_id", "raid_id"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id FROM absences WHERE player_id = $1 AND raid_id = $2", abs.Player.ID, abs.Raid.ID).
			Return(pgxRows, nil)

		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO absences (player_id,raid_id) VALUES ($1,$2)", abs.Player.ID, abs.Raid.ID).
			Return(nil, nil)

		_, err := pgBackend.CreateAbsence(context.Background(), abs)
		assert.NoError(t, err)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		abs := entity.Absence{
			ID: 0,
			Player: &entity.Player{
				ID:   0,
				Name: "test",
			},
			Raid: &entity.Raid{
				ID:   0,
				Name: "test",
				Date: time.Now(),
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.CreateAbsence(ctx, abs)
		assert.Error(t, err)
	})

	t.Run("Error inserting", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		abs := entity.Absence{
			ID: 0,
			Player: &entity.Player{
				ID:   0,
				Name: "test",
			},
			Raid: &entity.Raid{
				ID:   0,
				Name: "test",
				Date: time.Now(),
			},
		}

		columns := []string{"id", "player_id", "raid_id"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id FROM absences WHERE player_id = $1 AND raid_id = $2", abs.Player.ID, abs.Raid.ID).
			Return(pgxRows, nil)

		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO absences (player_id,raid_id) VALUES ($1,$2)", abs.Player.ID, abs.Raid.ID).
			Return(nil, errors.New("error"))

		_, err := pgBackend.CreateAbsence(context.Background(), abs)
		assert.Error(t, err)
	})
}
