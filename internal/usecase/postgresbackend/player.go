package postgresbackend

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgconn"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

// SearchPlayer is a function which call backend to Search a Player Object.
// It can search by playerID, name or discordName.
// players returned doesn't contain strikes, fails, missed raids and loots.
func (pg *PG) SearchPlayer(ctx context.Context, playerID int, name, discordName string) ([]entity.Player, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Player/SearchPlayer")
	span.SetAttributes(
		attribute.String("playerName", name),
		attribute.String("discordName", discordName),
		attribute.Int("playerID", playerID))
	defer span.End()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchPlayer - ctx.Done: request took too much time to be proceed")
	default:
		var players []entity.Player
		var param any
		var field string

		switch {
		case playerID != -1:
			param = playerID
			field = "id"
		case name != "":
			param = name
			field = "name"
		case discordName != "":
			param = discordName
			field = "discord_id"
		}

		sql, _, err := pg.Builder.Select("id", "name", "discord_id").From("players").Where("$1 = $2").ToSql()
		if err != nil {
			return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, field, param)
		if err != nil {
			return nil, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
		}

		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var player entity.Player
				err := rows.Scan(&player.ID, &player.Name, &player.DiscordName)
				if err != nil {
					return nil, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
				}
				players = append(players, player)
			}
		}
		return players, nil
	}
}

// CreatePlayer is a function which call backend to Create a Player Object.
func (pg *PG) CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Player/CreatePlayer")
	span.SetAttributes(attribute.String("playerName", player.Name))
	defer span.End()
	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("database - CreatePlayer - ctx.Done: request took too much time to be proceed")
	default:
		if player.DiscordName == "" {
			player.DiscordName = "tmp_" + strconv.FormatInt(time.Now().Unix(), 10)
		}

		req, args, errInsert := pg.Builder.
			Insert("players").
			Columns("name", "discord_id").
			Values(player.Name, player.DiscordName).
			Suffix("RETURNING \"id\"").
			ToSql()
		if errInsert != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Builder.Insert: %w", errInsert)
		}
		row := pg.Pool.QueryRow(ctx, req, args...)
		if row == nil {
			return entity.Player{}, fmt.Errorf("call insert player, returned row is empty")
		}
		err := row.Scan(&player.ID)
		if err != nil {
			var pgErr *pgconn.PgError
			ok := errors.As(err, &pgErr)
			if ok && pgErr.Code == "23505" {
				return entity.Player{}, fmt.Errorf("player already exists")
			}
			return entity.Player{}, fmt.Errorf("scan row from query row on insert player: %w", err)
		}
		return player, nil
	}
}

// ReadPlayer is a function which call backend to Read a Player Object.
func (pg *PG) ReadPlayer(ctx context.Context, playerID int) (entity.Player, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Player/ReadPlayer")
	span.SetAttributes(
		attribute.Int("playerID", playerID))
	defer span.End()

	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("database - ReadPlayer - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Select("id", "name").From("players").Where("id = $1").ToSql()
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.Builder.Select: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, strconv.FormatInt(int64(playerID), 10))
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		var player entity.Player
		if rows.Next() {
			err := rows.Scan(&player.ID, &player.Name)
			if err != nil {
				return entity.Player{}, fmt.Errorf("database - ReadPlayer - rows.Scan: %w", err)
			}
			return player, nil
		}
		return entity.Player{}, fmt.Errorf("database - ReadPlayer - player not found")
	}
}

// UpdatePlayer is a function which call backend to Update a Player Object.
func (pg *PG) UpdatePlayer(ctx context.Context, player entity.Player) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Player/UpdatePlayer")
	span.SetAttributes(
		attribute.String("playerName", player.Name),
		attribute.String("discordName", player.DiscordName),
		attribute.Int("playerID", player.ID))
	defer span.End()

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdatePlayer - ctx.Done: request took too much time to be proceed")
	default:
		if player.DiscordName == "" {
			player.DiscordName = "tmp_" + strconv.FormatInt(time.Now().Unix(), 10)
		}
		sql, args, err := pg.Builder.Update("players").
			Set("name", player.Name).
			Set("discord_id", player.DiscordName).
			Where("id = ?", player.ID).ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdatePlayer - r.Builder.Update: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("database - UpdatePlayer - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

// DeletePlayer is a function which call backend to Delete a Player Object.
func (pg *PG) DeletePlayer(ctx context.Context, player entity.Player) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Player/DeletePlayer")
	span.SetAttributes(
		attribute.String("playerName", player.Name),
		attribute.String("discordName", player.DiscordName),
		attribute.Int("playerID", player.ID))
	defer span.End()

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeletePlayer - ctx.Done: request took too much time to be proceed")
	default:
		sqlQuery := pg.Builder.Delete("players")
		args := make([]any, 0)
		if player.ID != 0 {
			sqlQuery = sqlQuery.Where("id = ?", player.ID)
			args = append(args, player.ID)
		}
		if player.Name != "" {
			sqlQuery = sqlQuery.Where("name = ?", player.Name)
			args = append(args, player.Name)
		}
		if player.DiscordName != "" {
			sqlQuery = sqlQuery.Where("discord_id = ?", player.DiscordName)
			args = append(args, player.DiscordName)
		}
		sql, _, err := sqlQuery.ToSql()
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Builder.Delete: %w", err)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database - DeletePlayer - player not found")
		}
		return nil
	}
}
