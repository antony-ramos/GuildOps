package postgresbackend_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/antony-ramos/guildops/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"
	"github.com/antony-ramos/guildops/pkg/postgres"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

//nolint:dupl
func TestPG_DeleteRaid(t *testing.T) {
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
			"DELETE FROM raids WHERE id = $1", 1).
			Return(nil, nil)

		err := pgBackend.DeleteRaid(context.Background(), 1)
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

		err := pgBackend.DeleteRaid(ctx, 1)
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
			"DELETE FROM raids WHERE id = $1", 1).
			Return(nil, errors.New("error"))

		err := pgBackend.DeleteRaid(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("raid is not deleted", func(t *testing.T) {
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
			"DELETE FROM raids WHERE id = $1", 1).
			Return(commandTag, nil)

		err := pgBackend.DeleteRaid(context.Background(), 1)
		assert.Error(t, err)
	})
}


func TestPG_UpdateRaid(t *testing.T) {
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

		raid := entity.Raid{
			ID:   1,
			Name: "test",
			Date: time.Now(),
		}

		columns := []string{"id", "name", "date", "difficulty"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(raid.ID, raid.Name, raid.Date, raid.Difficulty).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, date, difficulty FROM raids WHERE id = $1", raid.ID).
			Return(pgxRows, nil)

		raid = entity.Raid{
			ID:   1,
			Name: "test two",
			Date: time.Now(),
		}

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE raids SET name = $1, date = $2, difficulty = $3 WHERE id = $1", raid.ID).
			Return(nil, nil)

		err := pgBackend.UpdateRaid(context.Background(), raid)
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

		raid := entity.Raid{
			ID:   1,
			Name: "test",
			Date: time.Now(),
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := pgBackend.UpdateRaid(ctx, raid)
		assert.Error(t, err)
	})
}
