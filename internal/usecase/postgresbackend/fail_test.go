package postgresbackend_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestPG_searchFailOnParam(t *testing.T) {
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

		columns := []string{"id", "player_id", "raid_id", "reason"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, 0, "test").ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id, reason FROM fails WHERE player_id = $1", 1).
			Return(pgxRows, nil)

		fail, err := pgBackend.SearchFailOnParam(context.Background(), "player_id", 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(fail))
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

		_, err := pgBackend.SearchFailOnParam(ctx, "player_id", 1)
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
			"SELECT id, player_id, raid_id, reason FROM fails WHERE player_id = $1", 1).
			Return(nil, errors.New("query failed"))

		_, err := pgBackend.SearchFailOnParam(context.Background(), "player_id", 1)
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

		columns := []string{"id", "player_id", "raid_id", "reason", "toto"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, "test", "test", time.Now()).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id, reason FROM fails WHERE player_id = $1", 1).
			Return(pgxRows, nil)

		_, err := pgBackend.SearchFailOnParam(context.Background(), "player_id", 1)
		assert.Error(t, err)
	})
}

//nolint:dupl
func TestPG_DeleteFail(t *testing.T) {
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
			"DELETE FROM fails WHERE id = $1", 1).
			Return(nil, nil)

		err := pgBackend.DeleteFail(context.Background(), 1)
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

		err := pgBackend.DeleteFail(ctx, 1)
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
			"DELETE FROM fails WHERE id = $1", 1).
			Return(nil, errors.New("error"))

		err := pgBackend.DeleteFail(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("fail is not deleted", func(t *testing.T) {
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
			"DELETE FROM fails WHERE id = $1", 1).
			Return(commandTag, nil)

		err := pgBackend.DeleteFail(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestPG_CreateFail(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO fails (player_id,raid_id,reason) VALUES ($1,$2,$3)", 0, 0, "reason").
			Return(nil, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		fail := entity.Fail{
			ID:     0,
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
			Reason: "reason",
		}

		_, err := pgBackend.CreateFail(context.Background(), fail)
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

		fail := entity.Fail{
			ID:     0,
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
			Reason: "reason",
		}

		_, err := pgBackend.CreateFail(ctx, fail)
		assert.Error(t, err)
	})

	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Exec(gomock.Any(),
			"INSERT INTO fails (player_id,raid_id,reason) VALUES ($1,$2,$3)", 0, 0, "reason").
			Return(nil, errors.New("error"))

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}

		fail := entity.Fail{
			ID:     0,
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
			Reason: "reason",
		}

		_, err := pgBackend.CreateFail(context.Background(), fail)
		assert.Error(t, err)
	})
}

func TestPG_SearchFail(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}

		columns := []string{"id", "player_id", "raid_id", "reason"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(0, 0, 0, "test").ToPgxRows()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id, reason FROM fails WHERE player_id = $1", player.ID).
			Return(pgxRows, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		fails, err := pgBackend.SearchFail(context.Background(), "", player.ID, -1, "")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(fails))
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

		_, err := pgBackend.SearchFail(ctx, "", player.ID, -1, "")
		assert.Error(t, err)
	})
}

func TestPG_ReadFail(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fail := entity.Fail{
			ID:     1,
			Reason: "valid reason",
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
		}

		columns := []string{"id", "player_id", "raid_id", "reason"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(fail.ID, 0, 0, fail.Reason).ToPgxRows()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, player_id, raid_id, reason FROM fails WHERE id = $1", fail.ID).
			Return(pgxRows, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		failInput, err := pgBackend.ReadFail(context.Background(), 1)
		assert.NoError(t, err)
		failInput.Player = fail.Player
		failInput.Raid = fail.Raid
		assert.Equal(t, fail, failInput)
	})
}

func TestPG_UpdateFail(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fail := entity.Fail{
			ID:     1,
			Reason: "valid reason",
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
		}

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE fails SET reason = $1 WHERE id = $2", fail.Reason, fail.ID).
			Return(nil, nil)

		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		err := pgBackend.UpdateFail(context.Background(), fail)
		assert.NoError(t, err)
	})

	t.Run("Context cancelled", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fail := entity.Fail{
			ID:     1,
			Reason: "valid reason",
			Player: &entity.Player{},
			Raid:   &entity.Raid{},
		}

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		pgBackend := postgresbackend.PG{Postgres: &postgres.Postgres{
			Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			Pool:    mockPool,
		}}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := pgBackend.UpdateFail(ctx, fail)
		assert.Error(t, err)
	})
}
