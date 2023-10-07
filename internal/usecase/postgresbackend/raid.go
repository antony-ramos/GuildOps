package postgresbackend

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

// SearchRaid searches a raid in the database.
func (pg *PG) SearchRaid(
	ctx context.Context, raidName string, date time.Time, difficulty string,
) ([]entity.Raid, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/SearchRaid")
	defer span.End()
	span.SetAttributes(
		attribute.String("raidName", raidName),
		attribute.String("date", date.String()),
		attribute.String("difficulty", difficulty),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchRaid - ctx.Done: request took too much time to be proceed")
	default:
		var raids []entity.Raid
		if raidName != "" {
			sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("name = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Builder: %w", err)
			}
			rows, err := pg.Pool.Query(ctx, sql, raidName)
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
		if difficulty != "" {
			sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("difficulty = $1").ToSql()
			if err != nil {
				return nil, fmt.Errorf("database - SearchRaid - r.Builder: %w", err)
			}
			rows, err := pg.Pool.Query(ctx, sql, difficulty)
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
			rows, err := pg.Pool.Query(ctx, sql, date)
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

// CreateRaid creates a raid in the database.
func (pg *PG) CreateRaid(ctx context.Context, raid entity.Raid) (entity.Raid, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/CreateRaid")
	defer span.End()
	span.SetAttributes(
		attribute.String("raidName", raid.Name),
		attribute.String("difficulty", raid.Difficulty),
		attribute.String("date", raid.Date.String()),
	)
	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("database - CreateRaid - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.
			Select("name", "date", "difficulty").
			From("raids").
			Where("date = $1 AND difficulty = $2").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, raid.Date, raid.Difficulty)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			return entity.Raid{}, fmt.Errorf("raid already exists")
		}
		sql, _, errInsert := pg.Builder.
			Insert("raids").
			Columns("name", "date", "difficulty").
			Values(raid.Name, raid.Date, raid.Difficulty).ToSql()
		if errInsert != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder.Insert: %w", errInsert)
		}

		_, err = pg.Pool.Exec(ctx, sql, raid.Name, raid.Date, raid.Difficulty)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Pool.Exec: %w", err)
		}
		// get raid ID
		sql, _, err = pg.Builder.Select("id").From("raids").Where("name = $1 AND date = $2 AND difficulty = $3").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Builder: %w", err)
		}
		rows, err = pg.Pool.Query(ctx, sql, raid.Name, raid.Date, raid.Difficulty)
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

// ReadRaid returns a raid from the database.
func (pg *PG) ReadRaid(ctx context.Context, raidID int) (entity.Raid, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/ReadRaid")
	defer span.End()
	span.SetAttributes(
		attribute.Int("raidID", raidID),
	)

	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("database - ReadRaid - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("id = $1").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - ReadRaid - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, raidID)
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

// ReadRaidOnDate returns a raid from the database.
func (pg *PG) ReadRaidOnDate(ctx context.Context, date time.Time) (entity.Raid, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/ReadRaidOnDate")
	defer span.End()
	span.SetAttributes(
		attribute.String("date", date.Format("02/01/2006")),
	)
	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("database - ReadRaid - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Select("id", "name", "date", "difficulty").From("raids").Where("date = $1").ToSql()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - ReadRaid - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, date)
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

// UpdateRaid updates a raid in the database.
func (pg *PG) UpdateRaid(ctx context.Context, raid entity.Raid) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/UpdateRaid")
	defer span.End()
	span.SetAttributes(
		attribute.Int("raidID", raid.ID),
		attribute.String("raidName", raid.Name),
		attribute.String("difficulty", raid.Difficulty),
		attribute.String("date", raid.Date.Format("02/01/2006")),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateRaid - ctx.Done: request took too much time to be proceed")
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
		sql, _, err := pg.Builder.
			Update("raids").
			Set("name", oldRaid.Name).
			Set("date", oldRaid.Date).
			Set("difficulty", oldRaid.Difficulty).
			Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdateRaid - r.Builder.Update: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, oldRaid.ID)
		if err != nil {
			return fmt.Errorf("database - UpdateRaid - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

func (pg *PG) DeleteRaid(ctx context.Context, raidID int) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Raid/DeleteRaid")
	defer span.End()
	span.SetAttributes(
		attribute.Int("raidID", raidID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeleteRaid - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Delete("raids").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeleteRaid - r.Builder.Delete: %w", err)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, raidID)
		if err != nil {
			return fmt.Errorf("database - DeleteRaid - r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database - DeleteRaid - raid not found")
		}
		return nil
	}
}
