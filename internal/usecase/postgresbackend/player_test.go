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
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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
			Name: "playername",
		}

		rows := pgxpoolmock.NewRows([]string{"id"}).AddRow(1).ToPgxRows()
		rows.Next()
		mockPool.EXPECT().QueryRow(gomock.Any(),
			"INSERT INTO players (name,discord_id) VALUES ($1,$2) RETURNING \"id\"", player.Name).
			Return(rows)

		p, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.NoError(t, err)
		assert.NotEqual(t, entity.Player{}, p)
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

	t.Run("query row is nil", func(t *testing.T) {
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

		mockPool.EXPECT().QueryRow(gomock.Any(),
			"INSERT INTO players (name,discord_id) VALUES ($1,$2) RETURNING \"id\"", player.Name).
			Return(pgx.Row(nil))

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

		rows := pgxpoolmock.NewRows([]string{"id"}).RowError(0, &pgconn.PgError{Code: "23505"}).AddRow(1).ToPgxRows()
		rows.Next()
		mockPool.EXPECT().QueryRow(gomock.Any(),
			"INSERT INTO players (name,discord_id) VALUES ($1,$2) RETURNING \"id\"", player.Name).
			Return(rows)

		_, err := pgBackend.CreatePlayer(context.Background(), player)
		assert.Equal(t, errors.New("player already exists"), err)
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

//nolint:dupl
func TestPG_SearchPlayer(t *testing.T) {
	t.Parallel()
	t.Run("Success with playerID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := 1
		name := ""
		discordName := ""

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

		columns := []string{"id", "name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name, player.DiscordName).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "id", playerID).
			Return(pgxRows, nil)

		p, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.NoError(t, err)
		assert.Equal(t, player, p[0])
	})

	t.Run("Failed with playerID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := 1
		name := ""
		discordName := ""

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "id", playerID).
			Return(nil, errors.New("error"))

		_, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.Error(t, err)
	})

	t.Run("Success with name", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := "playername"
		discordName := ""

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

		columns := []string{"id", "name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name, player.DiscordName).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "name", name).
			Return(pgxRows, nil)

		p, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.NoError(t, err)
		assert.Equal(t, player, p[0])
	})

	t.Run("Failed with name", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := "miss"
		discordName := ""

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "name", name).
			Return(nil, errors.New("error"))

		_, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.Error(t, err)
	})

	t.Run("Success with discordName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := ""
		discordName := "discordName"

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

		columns := []string{"id", "name", "discord_id"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(player.ID, player.Name, player.DiscordName).ToPgxRows()
		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "discord_id", discordName).
			Return(pgxRows, nil)

		p, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)

		assert.NoError(t, err)
		assert.Equal(t, player, p[0])
	})

	t.Run("Failed with discordID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := ""
		discordName := "miss"

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		mockPool.EXPECT().Query(gomock.Any(),
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "discord_id", discordName).
			Return(nil, errors.New("error"))

		_, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.Error(t, err)
	})

	t.Run("Failed with discordName, Rows are incorrect", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := ""
		discordName := "discordName"

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
			"SELECT id, name, discord_id FROM players WHERE $1 = $2", "discord_id", discordName).
			Return(pgxRows, nil)

		_, err := pgBackend.SearchPlayer(context.Background(), playerID, name, discordName)
		assert.Error(t, err)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		playerID := -1
		name := ""
		discordName := ""

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := pgBackend.SearchPlayer(ctx, playerID, name, discordName)
		assert.Error(t, err)
	})
}

func TestPG_UpdatePlayer(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
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
			ID:          1,
			Name:        "playername",
			DiscordName: "toto",
		}

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE players SET name = $1, discord_id = $2 WHERE id = $3", player.Name, player.DiscordName, player.ID).
			Return(nil, nil)

		err := pgBackend.UpdatePlayer(context.Background(), player)
		assert.NoError(t, err)
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

		err := pgBackend.UpdatePlayer(ctx, entity.Player{})
		assert.Error(t, err)
	})

	t.Run("Update failed", func(t *testing.T) {
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
			ID:          1,
			Name:        "playername",
			DiscordName: "toto",
		}

		mockPool.EXPECT().Exec(gomock.Any(),
			"UPDATE players SET name = $1, discord_id = $2 WHERE id = $3", player.Name, player.DiscordName, player.ID).
			Return(nil, errors.New("error"))

		err := pgBackend.UpdatePlayer(context.Background(), player)
		assert.Error(t, err)
	})
}

func TestPG_DeletePlayer(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
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
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM players WHERE id = $1", 1).
			Return(nil, nil)

		err := pgBackend.DeletePlayer(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("Context cancelled", func(t *testing.T) {
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

		err := pgBackend.DeletePlayer(ctx, 1)
		assert.Error(t, err)
	})

	t.Run("Query failed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM players WHERE id = $1", 1).
			Return(nil, errors.New("error"))

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		err := pgBackend.DeletePlayer(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("Player is not deleted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		commandTag := pgconn.CommandTag("DELETE 0")
		mockPool.EXPECT().Exec(gomock.Any(),
			"DELETE FROM players WHERE id = $1", 1).
			Return(commandTag, nil)

		pgBackend := postgresbackend.PG{
			Postgres: &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mockPool,
			},
		}

		err := pgBackend.DeletePlayer(context.Background(), 1)
		assert.Error(t, err)
	})
}
