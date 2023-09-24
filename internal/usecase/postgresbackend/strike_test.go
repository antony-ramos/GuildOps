package postgresbackend_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/jackc/pgconn"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPG_searchStrikeOnParam(t *testing.T) {
	t.Parallel()
	t.Run("Search with playerID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		columns := []string{"id", "player_id", "season", "reason", "created_at"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, "test", "test", time.Now()).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, season, reason, created_at FROM strikes WHERE player_id = $1", 1).
			Return(pgxRows, nil)

		strike, err := pgBackend.SearchStrikeOnParam(context.Background(), "player_id", 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(strike))
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.SearchStrikeOnParam(ctx, "player_id", 1)
		assert.Error(t, err)
	})

	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, season, reason, created_at FROM strikes WHERE player_id = $1", 1).
			Return(nil, errors.New("error"))

		_, err := pgBackend.SearchStrikeOnParam(context.Background(), "player_id", 1)
		assert.Error(t, err)
	})

	t.Run("Scan failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		columns := []string{"id", "player_id", "season", "reason", "created_at", "toto"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, "test", "test", time.Now(), "toto").ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, season, reason, created_at FROM strikes WHERE player_id = $1", 1).
			Return(pgxRows, nil)

		_, err := pgBackend.SearchStrikeOnParam(context.Background(), "player_id", 1)
		assert.Error(t, err)
	})
}

//nolint:dupl
func TestPG_DeleteStrike(t *testing.T) {
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
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM strikes WHERE id = $1", 1).
			Return(nil, nil)

		err := pgBackend.DeleteStrike(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := pgBackend.DeleteStrike(ctx, 1)
		assert.Error(t, err)
	})

	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM strikes WHERE id = $1", 1).
			Return(nil, errors.New("error"))

		err := pgBackend.DeleteStrike(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("strike is not deleted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		commandTag := pgconn.CommandTag("DELETE 0")
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM strikes WHERE id = $1", 1).
			Return(commandTag, nil)

		err := pgBackend.DeleteStrike(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestPG_CreateStrike(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO strikes (player_id,season,reason) VALUES ($1,$2,$3)", 1, "season", "reason").
			Return(nil, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		strike := entity.Strike{
			ID:     0,
			Player: &entity.Player{},
			Season: "season",
			Reason: "reason",
		}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		err := pgBackend.CreateStrike(context.Background(), strike, player)
		assert.NoError(t, err)
	})

	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		strike := entity.Strike{
			ID:     0,
			Player: &entity.Player{},
			Season: "season",
			Reason: "reason",
		}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		err := pgBackend.CreateStrike(ctx, strike, player)
		assert.Error(t, err)
	})

	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO strikes (player_id,season,reason) VALUES ($1,$2,$3)", 1, "season", "reason").
			Return(nil, errors.New("error"))

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		strike := entity.Strike{
			ID:     0,
			Player: &entity.Player{},
			Season: "season",
			Reason: "reason",
		}

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		err := pgBackend.CreateStrike(context.Background(), strike, player)
		assert.Error(t, err)
	})
}

func TestPG_SearchStrike(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"id", "player_id", "season", "reason", "created_at"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, "test", "test", time.Now()).ToPgxRows()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, season, reason, created_at FROM strikes WHERE player_id = $1", player.ID).
			Return(pgxRows, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		strikes, err := pgBackend.SearchStrike(context.Background(), player.ID, time.Time{}, "", "")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(strikes))
	})

	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.SearchStrike(ctx, player.ID, time.Time{}, "", "")
		assert.Error(t, err)
	})
}

func TestPG_ReadStrike(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strike := entity.Strike{
			ID:     1,
			Reason: "valid reason",
			Date:   time.Now(),
			Season: "DF/S2",
		}

		columns := []string{"id", "player_id", "season", "reason", "created_at"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(strike.ID, 0, strike.Season, strike.Reason, strike.Date).ToPgxRows()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, season, reason, created_at FROM strikes WHERE id = $1", strike.ID).
			Return(pgxRows, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		strikeInput, err := pgBackend.ReadStrike(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, strike, strikeInput)
	})
}

func TestPG_UpdateStrike(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strike := entity.Strike{
			ID:     1,
			Reason: "valid reason",
			Date:   time.Now(),
			Season: "DF/S2",
		}

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE strikes SET season = $1, reason = $2 WHERE id = $1", strike.ID).
			Return(nil, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		err := pgBackend.UpdateStrike(context.Background(), strike)
		assert.NoError(t, err)
	})

	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		strike := entity.Strike{
			ID:     1,
			Reason: "valid reason",
			Date:   time.Now(),
			Season: "DF/S2",
		}

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := pgBackend.UpdateStrike(ctx, strike)
		assert.Error(t, err)
	})
}
