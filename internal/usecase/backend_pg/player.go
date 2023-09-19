package backend_pg

import (
	"context"
	"fmt"
	"strconv"

	"github.com/coven-discord-bot/internal/entity"
)

// SearchPlayer is a function which call backend to Search a Player Object.
func (pg *PG) SearchPlayer(ctx context.Context, id int, name string) ([]entity.Player, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var players []entity.Player

		var sql string
		var err error
		switch {
		case id != -1:
			sql, _, err = pg.Builder.Select("id", "name").From("players").Where("name = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
			}
		case name != "":
			sql, _, err = pg.Builder.Select("id", "name").From("players").Where("name = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchPlayer - r.Builder: %w", err)
			}
		}

		rows, err := pg.Pool.Query(context.Background(), sql, name)
		if err != nil {
			return nil, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			player, err := func() (entity.Player, error) {
				var player entity.Player
				err := rows.Scan(&player.ID, &player.Name)
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
					return entity.Player{}, fmt.Errorf("database - SearchPlayer - r.Builder.Select: %w", err)
				}
				r, err := pg.Pool.Query(context.Background(), sql, strconv.FormatInt(int64(player.ID), 10))
				if err != nil {
					return entity.Player{}, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
				}
				defer r.Close()
				for r.Next() {
					loot := entity.Loot{}
					raid := entity.Raid{}

					err := r.Scan(&loot.ID, &loot.Name, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date)
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
					}
					loot.Raid = &raid
					loot.Player = &player
					player.Loots = append(player.Loots, loot)
				}

				// populate player.MissedRaids list
				sql, _, err = pg.Builder.
					Select("absences.id", "absences.player_id", "absences.raid_id", "raids.name", "raids.difficulty", "raids.date").
					From("absences").
					Join("raids ON raids.id = absences.raid_id").
					Where("absences.player_id = $1").ToSql()
				if err != nil {
					return entity.Player{}, fmt.Errorf("database - SearchPlayer - r.Builder.Select: %w", err)
				}
				r, err = pg.Pool.Query(context.Background(), sql, strconv.FormatInt(int64(player.ID), 10))
				if err != nil {
					return entity.Player{}, fmt.Errorf("database - SearchPlayer - r.Pool.Query: %w", err)
				}
				defer r.Close()
				for r.Next() {
					raid := entity.Raid{}
					err := r.Scan(nil, nil, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date)
					if err != nil {
						return entity.Player{}, fmt.Errorf("database - SearchPlayer - rows.Scan: %w", err)
					}

					player.MissedRaids = append(player.MissedRaids, raid)
				}
				return player, nil
			}()
			if err != nil {
				return nil, err
			}
			players = append(players, player)
		}
		return players, nil
	}
}

// CreatePlayer is a function which call backend to Create a Player Object.
func (pg *PG) CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error) {
	select {
	case <-ctx.Done():
		return entity.Player{}, ctx.Err()
	default:
		sql, _, err := pg.Builder.Select("name").From("players").Where("name = $1").ToSql()
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(context.Background(), sql, player.Name)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - player already exists")
		}
		sql, args, errInsert := pg.Builder.Insert("players").Columns("name").Values(player.Name).ToSql()
		if errInsert != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Builder.Insert: %w", errInsert)
		}
		_, err = pg.Pool.Exec(context.Background(), sql, args...)
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - CreatePlayer - r.Pool.Exec: %w", err)
		}
		// player.ID = id
		return player, nil
	}
}

// ReadPlayer is a function which call backend to Read a Player Object.
func (pg *PG) ReadPlayer(ctx context.Context, playerID int) (entity.Player, error) {
	select {
	case <-ctx.Done():
		return entity.Player{}, ctx.Err()
	default:
		sql, _, err := pg.Builder.Select("id", "name").From("players").Where("name = $1").ToSql()
		if err != nil {
			return entity.Player{}, fmt.Errorf("database - ReadPlayer - r.Builder.Select: %w", err)
		}
		rows, err := pg.Pool.Query(context.Background(), sql, strconv.FormatInt(int64(playerID), 10))
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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		sql, args, err := pg.Builder.Update("players").Set("name", player.Name).Where("id = ?", player.ID).ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdatePlayer - r.Builder.Update: %w", err)
		}
		_, err = pg.Pool.Exec(context.Background(), sql, args...)
		if err != nil {
			return fmt.Errorf("database - UpdatePlayer - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

// DeletePlayer is a function which call backend to Delete a Player Object.
func (pg *PG) DeletePlayer(ctx context.Context, playerID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		sql, _, err := pg.Builder.Delete("players").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Builder.Delete: %w", err)
		}
		_, err = pg.Pool.Exec(context.Background(), sql, playerID)
		if err != nil {
			return fmt.Errorf("database - DeletePlayer - r.Pool.Exec: %w", err)
		}
		return nil
	}
}
