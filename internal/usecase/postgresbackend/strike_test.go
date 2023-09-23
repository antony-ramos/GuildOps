package postgresbackend_test

import (
	"context"
	"errors"
	"testing"
	"time"

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
