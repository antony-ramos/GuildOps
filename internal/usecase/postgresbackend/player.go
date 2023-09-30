package postgresbackend

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

// SearchPlayer is a function which call backend to Search a Player Object.
// TODO: refactor this function to avoid cyclomatic complexity
//
//nolint:gocyclo
func (pg *PG) SearchPlayer(ctx context.Context, playerID int, name, discordName string) ([]entity.Player, error) {
	ctx, span := otel.Tracer("postgresbackend").Start(ctx, "SearchPlayer")
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

		var sql string
		var err error
		var rows pgx.Rows
		switch {
		case playerID != -1:
			sql, _, err = pg.Builder.Select("id", "name", "discord_id").From("players").Where("player_id = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
			}
			rows, err = pg.Pool.Query(ctx, sql, playerID)
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
			}
		case name != "":
			sql, _, err = pg.Builder.Select("id", "name", "discord_id").From("players").Where("name = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
			}
			rows, err = pg.Pool.Query(ctx, sql, name)
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
			}
		case discordName != "":
			sql, _, err = pg.Builder.Select("id", "name", "discord_id").From("players").Where("discord_id = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
			}
			rows, err = pg.Pool.Query(ctx, sql, discordName)
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
			}
		}
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				player, err := func() (entity.Player, error) {
					var player entity.Player
					err := rows.Scan(&player.ID, &player.Name, &player.DiscordName)
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
					}

					// populate player.Loot list
					sql, _, err = pg.Builder.
						Select("loots.id", "loots.name", "loots.raid_id", "raids.name", "raids.difficulty", "raids.date").
						From("loots").
						Join("raids ON raids.id = loots.raid_id").
						Where("loots.player_id = $1").ToSql()
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Builder.Select: %w", err)
					}
					playerRows, err := pg.Pool.Query(ctx, sql, strconv.FormatInt(int64(player.ID), 10))
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Pool.Query: %w", err)
					}
					defer playerRows.Close()
					for playerRows.Next() {
						loot := entity.Loot{}
						raid := entity.Raid{}

						err := playerRows.Scan(&loot.ID, &loot.Name, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date)
						if err != nil {
							return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
						}
						loot.Raid = &raid
						loot.Player = &player
						player.Loots = append(player.Loots, loot)
					}

					// populate player.MissedRaids list
					sql, _, err = pg.Builder.
						Select("absences.id", "absences.player_id", "absences.raid_id",
							"raids.name", "raids.difficulty", "raids.date").
						From("absences").
						Join("raids ON raids.id = absences.raid_id").
						Where("absences.player_id = $1").ToSql()
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Builder.Select: %w", err)
					}
					playerRows, err = pg.Pool.Query(ctx, sql, strconv.FormatInt(int64(player.ID), 10))
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Pool.Query: %w", err)
					}
					defer playerRows.Close()
					for playerRows.Next() {
						raid := entity.Raid{}
						err := playerRows.Scan(nil, nil, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date)
						if err != nil {
							return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
						}

						player.MissedRaids = append(player.MissedRaids, raid)
					}

					// populate player.Strikes list
					sql, _, err = pg.Builder.
						Select("id", "season", "reason", "created_at").
						From("strikes").
						Where("player_id = $1").ToSql()
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Builder.Select: %w", err)
					}
					strikesRows, err := pg.Pool.Query(ctx, sql, strconv.FormatInt(int64(player.ID), 10))
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Pool.Query: %w", err)
					}
					defer strikesRows.Close()
					for strikesRows.Next() {
						strike := entity.Strike{}
						err := strikesRows.Scan(&strike.ID, &strike.Season, &strike.Reason, &strike.Date)
						if err != nil {
							return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
						}
						player.Strikes = append(player.Strikes, strike)
					}

					// populate player.Fails list
					sql, _, err = pg.Builder.
						Select("id", "player_id", "raid_id", "reason").
						From("fails").
						Where("player_id = $1").ToSql()
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Builder.Select: %w", err)
					}
					failsRows, err := pg.Pool.Query(ctx, sql, strconv.FormatInt(int64(player.ID), 10))
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - playerRows.Pool.Query: %w", err)
					}
					defer failsRows.Close()
					for failsRows.Next() {
						fail := entity.Fail{}
						err := failsRows.Scan(&fail.ID, &fail.Player.ID, &fail.Raid.ID, &fail.Reason)
						if err != nil {
							return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
						}
						player.Fails = append(player.Fails, fail)
					}

					return player, nil
				}()
				if err != nil {
					return nil, err
				}
				players = append(players, player)
			}
		}
		return players, nil
	}
}

// CreatePlayer is a function which call backend to Create a Player Object.
func (pg *PG) CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error) {
	ctx, span := otel.Tracer("postgresbackend").Start(ctx, "CreatePlayer")
	span.SetAttributes(attribute.String("playerName", player.Name))
	defer span.End()
	select {
	case <-ctx.Done():
		return entity.Player{}, fmt.Errorf("database - CreatePlayer - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Select("name").From("players").Where("name = $1").ToSql()
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, player.Name)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer: player already exists")
		}

		if player.DiscordName == "" {
			player.DiscordName = "tmp_" + strconv.FormatInt(time.Now().Unix(), 10)
		}

		sql, args, errInsert := pg.Builder.
			Insert("players").
			Columns("name", "discord_id").
			Values(player.Name, player.DiscordName).ToSql()
		if errInsert != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Builder.Insert: %w", errInsert)
		}
		_, err = pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Pool.Exec: %w", err)
		}
		player, err := pg.SearchPlayer(ctx, -1, player.Name, "")
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.SearchPlayer: %w", err)
		}
		return player[0], nil
	}
}

// ReadPlayer is a function which call backend to Read a Player Object.
func (pg *PG) ReadPlayer(ctx context.Context, playerID int) (entity.Player, error) {
	ctx, span := otel.Tracer("postgresbackend").Start(ctx, "ReadPlayer")
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
	ctx, span := otel.Tracer("postgresbackend").Start(ctx, "UpdatePlayer")
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
func (pg *PG) DeletePlayer(ctx context.Context, playerID int) error {
	ctx, span := otel.Tracer("postgresbackend").Start(ctx, "DeletePlayer")
	span.SetAttributes(
		attribute.Int("playerID", playerID))
	defer span.End()

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeletePlayer - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Delete("players").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Builder.Delete: %w", err)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, playerID)
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database - DeletePlayer - player not found")
		}
		return nil
	}
}
