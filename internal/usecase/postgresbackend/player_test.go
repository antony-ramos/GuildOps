package postgresbackend_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPG_CreatePlayer(t *testing.T) {
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

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT name FROM players WHERE name = $1", player.Name).
			Return(pgxRows, nil)

		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO players (name,discord_id) VALUES ($1,$2)", player.Name).
			Return(nil, nil)

		columns = []string{"player_id", "name", "discord_id"}
		pgxRows = pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name, player.DiscordName).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE name = $1", player.Name).
			Return(pgxRows, nil)

		columns = []string{"id", "name", "raid_id", "raid_name", "raid_difficulty", "raid_date"}
		pgxRows = pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT loots.id, loots.name, loots.raid_id, raids.name, raids.difficulty, raids.date "+
				"FROM loots JOIN raids ON raids.id = loots.raid_id "+
				"WHERE loots.player_id = $1", "1").
			Return(pgxRows, nil)

		columns = []string{
			"absences.id", "absences.player_id", "absences.raid_id",
			"raids.name", "raids.difficulty", "raids.date",
		}
		pgxRows = pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT absences.id, absences.player_id, absences.raid_id, raids.name, raids.difficulty, raids.date "+
				"FROM absences JOIN raids ON raids.id = absences.raid_id "+
				"WHERE absences.player_id = $1", "1").
			Return(pgxRows, nil)

		p, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.NoError(t, err)
		assert.Equal(t, player, p)
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

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.CreatePlayer(ctx, player)
		assert.Error(t, err)
	})

	t.Run("Search if player exists", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT name FROM players WHERE name = $1", player.Name).
			Return(nil, errors.New("error"))

		p, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.Error(t, err)
		assert.Equal(t, entity.Player{}, p)
	})

	t.Run("Cannot insert player", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT name FROM players WHERE name = $1", player.Name).
			Return(pgxRows, nil)

		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO players (name,discord_id) VALUES ($1,$2)", player.Name).
			Return(nil, errors.New("error"))

		p, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.Error(t, err)
		assert.Equal(t, entity.Player{}, p)
	})

	t.Run("Player already exists", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		// create empty rows
		columns := []string{"name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.Name, player.DiscordName).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT name FROM players WHERE name = $1", player.Name).
			Return(pgxRows, nil)

		_, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.Error(t, err)
	})
}

func TestPG_ReadPlayer(t *testing.T) {
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

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"id", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name FROM players WHERE id = $1", strconv.FormatInt(int64(player.ID), 10)).
			Return(pgxRows, nil)

		p, err := pgBackend.ReadPlayer(context.Background(), player.ID)
		assert.NoError(t, err)
		assert.Equal(t, player, p)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.ReadPlayer(ctx, 1)
		assert.Error(t, err)
	})

	// Error rows scan
	t.Run("Error Select", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"id", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name FROM players WHERE id = $1", strconv.FormatInt(int64(player.ID), 10)).
			Return(pgxRows, errors.New("error"))

		_, err := pgBackend.ReadPlayer(context.Background(), player.ID)
		assert.Error(t, err)
	})
}
