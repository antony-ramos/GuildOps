package postgresbackend

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

func (pg *PG) SearchLoot(
	ctx context.Context, name string, date time.Time, difficulty, playerName string,
) ([]entity.Loot, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Loot/SearchLoot")
	defer span.End()
	span.SetAttributes(
		attribute.String("date", date.String()),
		attribute.String("difficulty", difficulty),
		attribute.String("playerName", playerName),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchLoot - ctx.Done: request took too much time to be proceed")
	default:
		var loots []entity.Loot
		var request string
		selectSQL := pg.Builder.
			Select("loots.id", "loots.name", "loots.raid_id",
				"raids.name", "raids.difficulty", "raids.date",
				"loots.player_id", "players.name").
			From("loots").
			Join("raids ON raids.id = loots.raid_id").Join("players ON players.id = loots.player_id")

		count := 0
		var args []any
		if name != "" {
			count++
			selectSQL = selectSQL.Where("loots.name = $" + strconv.Itoa(count))
			args = append(args, name)
		}
		if date != (time.Time{}) {
			count++
			selectSQL = selectSQL.Where("raids.date = $" + strconv.Itoa(count))
			args = append(args, date)
		}
		if difficulty != "" {
			count++
			selectSQL = selectSQL.Where("raids.difficulty = $" + strconv.Itoa(count))
			args = append(args, difficulty)
		}
		if playerName != "" {
			count++
			selectSQL = selectSQL.Where("players.name = $" + strconv.Itoa(count))
			args = append(args, playerName)
		}

		request, _, err := selectSQL.ToSql()
		if err != nil {
			return nil, fmt.Errorf("create query to search loot: %w", err)
		}

		rows, err := pg.Pool.Query(ctx, request, args...)
		if err != nil {
			return nil, fmt.Errorf("send query to search loot: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var loot entity.Loot
			var raid entity.Raid
			var player entity.Player
			err := rows.Scan(&loot.ID, &loot.Name, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.ID, &player.Name)
			if err != nil {
				return nil, fmt.Errorf("populate loots table with row data: %w", err)
			}
			loot.Raid = &raid
			loot.Player = &player
			loots = append(loots, loot)
		}
		return loots, nil
	}
}

func (pg *PG) CreateLoot(ctx context.Context, loot entity.Loot) (entity.Loot, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Loot/CreateLoot")
	defer span.End()
	span.SetAttributes(
		attribute.String("lootName", loot.Name),
		attribute.String("raidName", loot.Raid.Date.Format("02/01/2006")),
		attribute.String("playerName", loot.Player.Name),
	)

	select {
	case <-ctx.Done():
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.
			Select("name", "raid_id", "player_id").
			From("loots").
			Where("name = $1").
			Where("raid_id = $2").
			Where("player_id = $3").ToSql()
		if err != nil {
			return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, loot.Name, loot.Raid.ID, loot.Player.ID)
		if err != nil {
			return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			return entity.Loot{}, fmt.Errorf("database - CreateLoot - loot already exists")
		}

		sql, args, errInsert := pg.Builder.
			Insert("loots").
			Columns("name", "raid_id", "player_id").
			Values(loot.Name, loot.Raid.ID, loot.Player.ID).ToSql()
		if errInsert != nil {
			return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Builder.Insert: %w", errInsert)
		}
		_, err = pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Pool.Exec: %w", err)
		}
		return loot, nil
	}
}

func (pg *PG) ReadLoot(ctx context.Context, lootID int) (entity.Loot, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Loot/ReadLoot")
	defer span.End()
	span.SetAttributes(
		attribute.Int("lootID", lootID),
	)

	select {
	case <-ctx.Done():
		return entity.Loot{}, fmt.Errorf("database - ReadLoot - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Select("id", "name", "raid_id", "player_id").From("loots").Where("id = $1").ToSql()
		if err != nil {
			return entity.Loot{}, fmt.Errorf("database - ReadLoot - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, lootID)
		if err != nil {
			return entity.Loot{}, fmt.Errorf("database - ReadLoot - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			loot := entity.Loot{}
			var raidID int
			var playerID int
			err := rows.Scan(&loot.ID, &loot.Name, &raidID, &playerID)
			if err != nil {
				return entity.Loot{}, fmt.Errorf("database - ReadLoot - rows.Scan: %w", err)
			}
			raid, err := pg.ReadRaid(ctx, raidID)
			if err != nil {
				return entity.Loot{}, fmt.Errorf("database - ReadLoot - pg.ReadRaid: %w", err)
			}
			loot.Raid = &raid

			player, err := pg.ReadPlayer(ctx, playerID)
			if err != nil {
				return entity.Loot{}, fmt.Errorf("database - ReadLoot - pg.ReadPlayer: %w", err)
			}
			loot.Player = &player

			return loot, nil
		}
		return entity.Loot{}, fmt.Errorf("database - ReadLoot - loot not found")
	}
}

func (pg *PG) UpdateLoot(ctx context.Context, loot entity.Loot) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Loot/UpdateLoot")
	defer span.End()
	span.SetAttributes(
		attribute.Int("lootID", loot.ID),
		attribute.String("lootName", loot.Name),
		attribute.Int("raidID", loot.Raid.ID),
		attribute.Int("playerID", loot.Player.ID),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateLoot - ctx.Done: request took too much time to be proceed")
	default:
		sql, args, err := pg.Builder.
			Update("loots").
			Set("name", loot.Name).
			Set("raid_id", loot.Raid.ID).
			Set("player_id", loot.Player.ID).
			Where("id = $4").ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdateLoot - r.Builder: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, args, loot.ID)
		if err != nil {
			return fmt.Errorf("database - UpdateLoot - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

func (pg *PG) DeleteLoot(ctx context.Context, lootID int) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Loot/DeleteLoot")
	defer span.End()
	span.SetAttributes(
		attribute.Int("lootID", lootID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeleteLoot - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, errInsert := pg.Builder.Delete("loots").Where("id = $1").ToSql()
		if errInsert != nil {
			return fmt.Errorf("database - DeleteLoot - r.Builder: %w", errInsert)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, lootID)
		if err != nil {
			return fmt.Errorf("database - DeleteLoot - r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database - DeleteLoot - loot not found")
		}
		return nil
	}
}
