package backend_pg

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"time"
)

// SearchRaid searches a raid in the database
func (pg *PG) SearchRaid(ctx context.Context, raidName string, date time.Time, difficulty string) ([]entity.Raid, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var raids []entity.Raid
		if raidName != "" {
			sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("name = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Builder: %w", err)
			}
			rows, err := pg.Pool.Query(context.Background(), sql, raidName)
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Pool.Query: %w", err)

			}
			defer rows.Close()
			for rows.Next() {
				var raid entity.Raid
				err := rows.Scan(&raid.ID, &raid.Name, &raid.Date, &raid.Difficulty)
				if err != nil {
					return nil, fmt.Errorf("database - SearchRaid - rows.Scan: %w", err)
				}
				raids = append(raids, raid)
			}
		}
		if !date.IsZero() {
			sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("date = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Builder: %w", err)
			}
			rows, err := pg.Pool.Query(context.Background(), sql, date)
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Pool.Query: %w", err)

			}
			defer rows.Close()
			for rows.Next() {
				var raid entity.Raid
				err := rows.Scan(&raid.ID, &raid.Name, &raid.Date, &raid.Difficulty)
				if err != nil {
					return nil, fmt.Errorf("database - SearchRaid - rows.Scan: %w", err)
				}
				raids = append(raids, raid)
			}
		}
		return raids, nil
	}
}

// CreateRaid creates a raid in the database
func (pg *PG) CreateRaid(ctx context.Context, raid entity.Raid) (entity.Raid, error) {
	select {
	case <-ctx.Done():
		return entity.Raid{}, ctx.Err()
	default:
		sql, _, err := pg.Builder.Select("name", "date", "difficulty").From("raids").Where("name = $1 AND date = $2 AND difficulty = $3").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(context.Background(), sql, raid.Name, raid.Date, raid.Difficulty)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Pool.Query: %w", err)

		}
		defer rows.Close()
		if rows.Next() {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - raid already exists")
		}
		sql, _, errInsert := pg.Builder.Insert("raids").Columns("name", "date", "difficulty").Values(raid.Name, raid.Date, raid.Difficulty).ToSql()
		if errInsert != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder.Insert: %w", errInsert)
		}

		_, err = pg.Pool.Exec(context.Background(), sql, raid.Name, raid.Date, raid.Difficulty)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Pool.Exec: %w", err)
		}
		// get raid ID
		sql, _, err = pg.Builder.Select("id").From("raids").Where("name = $1 AND date = $2 AND difficulty = $3").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder: %w", err)
		}
		rows, err = pg.Pool.Query(context.Background(), sql, raid.Name, raid.Date, raid.Difficulty)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&raid.ID)
			if err != nil {
				return entity.Raid{}, fmt.Errorf("database - CreateRaid - rows.Scan: %w", err)
			}
		}
		return raid, nil
	}
}

// ReadRaid returns a raid from the database
func (pg *PG) ReadRaid(ctx context.Context, raidID int) (entity.Raid, error) {
	select {
	case <-ctx.Done():
		return entity.Raid{}, ctx.Err()
	default:
		sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("id = $1").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - ReadRaid - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(context.Background(), sql, raidID)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - ReadRaid - r.Pool.Query: %w", err)

		}
		defer rows.Close()
		var raid entity.Raid
		for rows.Next() {
			err := rows.Scan(&raid.ID, &raid.Name, &raid.Date, &raid.Difficulty)
			if err != nil {
				return entity.Raid{}, fmt.Errorf("database - ReadRaid - rows.Scan: %w", err)
			}
		}
		return raid, nil
	}
}

// UpdateRaid updates a raid in the database
func (pg *PG) UpdateRaid(ctx context.Context, raid entity.Raid) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		oldRaid, err := pg.ReadRaid(ctx, raid.ID)
		if err != nil {
			return fmt.Errorf("database - UpdateRaid - r.ReadRaid: %w", err)
		}

		// Check if fields from raid in parameters. If different, updates them. Else, do nothing
		if raid.Name != "" {
			oldRaid.Name = raid.Name
		}
		if !raid.Date.IsZero() {
			oldRaid.Date = raid.Date
		}
		if raid.Difficulty != "" {
			oldRaid.Difficulty = raid.Difficulty
		}

		// Update raid in database
		sql, _, err := pg.Builder.Update("raids").Set("name", oldRaid.Name).Set("date", oldRaid.Date).Set("difficulty", oldRaid.Difficulty).Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdateRaid - r.Builder.Update: %w", err)
		}
		_, err = pg.Pool.Exec(context.Background(), sql, oldRaid.ID)
		if err != nil {
			return fmt.Errorf("database - UpdateRaid - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

func (pg *PG) DeleteRaid(ctx context.Context, raidID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		sql, _, err := pg.Builder.Delete("raids").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeleteRaid - r.Builder.Delete: %w", err)
		}
		_, err = pg.Pool.Exec(context.Background(), sql, raidID)
		if err != nil {
			return fmt.Errorf("database - DeleteRaid - r.Pool.Exec: %w", err)
		}
		return nil
	}
}
