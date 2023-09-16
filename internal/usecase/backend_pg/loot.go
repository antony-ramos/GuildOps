package backend_pg

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"time"
)

func (pg *PG) SearchLoot(ctx context.Context, name string, date time.Time, difficulty string) ([]entity.Loot, error) {
	var loots []entity.Loot
	sql, _, err := pg.Builder.Select("loots.id", "loots.name", "loots.raid_id", "raids.name", "raids.difficulty", "raids.date", "loots.player_id", "players.name").From("loots").Join("raids ON raids.id = loots.raid_id").Join("players ON players.id = loots.player_id").Where("loots.name = $1").Where("raids.date = $2").Where("raids.difficulty = $3").ToSql()
	if err != nil {
		return nil, fmt.Errorf("database - SearchLoot - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, name, date, difficulty)
	if err != nil {
		return nil, fmt.Errorf("database - SearchLoot - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var loot entity.Loot
		var raid entity.Raid
		var player entity.Player
		err := rows.Scan(&loot.ID, &loot.Name, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.ID, &player.Name)
		if err != nil {
			return nil, fmt.Errorf("database - SearchLoot - rows.Scan: %w", err)
		}
		loot.Raid = &raid
		loot.Player = &player
		loots = append(loots, loot)
	}
	return loots, nil
}

func (pg *PG) CreateLoot(ctx context.Context, loot entity.Loot) (entity.Loot, error) {
	// Verify if loot already exists
	sql, _, err := pg.Builder.Select("name", "raid_id", "player_id").From("loots").Where("name = $1").Where("raid_id = $2").Where("player_id = $3").ToSql()
	if err != nil {
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, loot.Name, loot.Raid.ID, loot.Player.ID)
	if err != nil {
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Pool.Query: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - loot already exists")
	}

	sql, args, errInsert := pg.Builder.Insert("loots").Columns("name", "raid_id", "player_id").Values(loot.Name, loot.Raid.ID, loot.Player.ID).ToSql()
	if errInsert != nil {
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Builder.Insert: %w", errInsert)
	}
	_, err = pg.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return entity.Loot{}, fmt.Errorf("database - CreateLoot - r.Pool.Exec: %w", err)
	}
	return loot, nil
}

func (pg *PG) ReadLoot(ctx context.Context, lootID int) (entity.Loot, error) {
	sql, _, err := pg.Builder.Select("id", "name", "raid_id", "player_id").From("loots").Where("id = $1").ToSql()
	if err != nil {
		return entity.Loot{}, fmt.Errorf("database - ReadLoot - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, lootID)
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

func (pg *PG) UpdateLoot(ctx context.Context, loot entity.Loot) error {
	// Update loot in database based on loot parameter
	sql, args, err := pg.Builder.Update("loots").Set("name", loot.Name).Set("raid_id", loot.Raid.ID).Set("player_id", loot.Player.ID).Where("id = $1").ToSql()
	if err != nil {
		return fmt.Errorf("database - UpdateLoot - r.Builder: %w", err)
	}
	_, err = pg.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("database - UpdateLoot - r.Pool.Exec: %w", err)
	}
	return nil
}

func (pg *PG) DeleteLoot(ctx context.Context, lootID int) error {
	sql, _, errInsert := pg.Builder.Delete("loots").Where("id = $1").ToSql()
	if errInsert != nil {
		return fmt.Errorf("database - DeleteLoot - r.Builder: %w", errInsert)
	}
	_, err := pg.Pool.Exec(context.Background(), sql, lootID)
	if err != nil {
		return fmt.Errorf("database - DeleteLoot - r.Pool.Exec: %w", err)
	}
	return nil
}
