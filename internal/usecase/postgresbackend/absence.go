package postgresbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

func (pg *PG) searchAbsenceOnParam(ctx context.Context, paramName string, param interface{}) ([]entity.Absence, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchAbsence - searchAbsenceOnParam - ctx.Done: %w", ctx.Err())
	default:
		sql, _, err := pg.Builder.Select("absences.id", "absences.player_id", "absences.raid_id",
			"raids.name", "raids.difficulty", "raids.date", "players.name").
			From("absences").
			Join("raids ON raids.id = absences.raid_id").
			Join("players ON players.id = absences.player_id").
			Where(fmt.Sprintf("%s = $1", paramName)).ToSql()
		if err != nil {
			return nil, fmt.Errorf("database - SearchAbsence - searchAbsenceOnParam - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, param)
		if err != nil {
			return nil, fmt.Errorf("database - SearchAbsence - searchAbsenceOnParam - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		var absences []entity.Absence

		for rows.Next() {
			var absence entity.Absence
			var raid entity.Raid
			var player entity.Player
			err := rows.Scan(&absence.ID, &player.ID, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.Name)
			if err != nil {
				return nil, fmt.Errorf("database - SearchAbsence - searchAbsenceOnParam - rows.Scan: %w", err)
			}
			absence.Player = &player
			absence.Raid = &raid
			absences = append(absences, absence)
		}
		return absences, nil
	}
}

func (pg *PG) SearchAbsence(
	ctx context.Context, playerName string, playerID int, date time.Time,
) ([]entity.Absence, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchAbsence - ctx.Done: %w", ctx.Err())
	default:
		var absences []entity.Absence
		switch {
		case playerID != -1 && playerName == "":
			a, err := pg.searchAbsenceOnParam(ctx, "player_id", playerID)
			if err != nil {
				return nil, err
			}
			absences = append(absences, a...)
		case playerID == -1 && playerName != "":
			a, err := pg.searchAbsenceOnParam(ctx, "players.name", playerName)
			if err != nil {
				return nil, err
			}
			absences = append(absences, a...)
		case playerID != -1 && playerName != "" && !date.IsZero():
			a, err := pg.searchAbsenceOnParam(ctx, "date", date)
			if err != nil {
				return nil, err
			}
			absences = append(absences, a...)
		}
		return absences, nil
	}
}

func (pg *PG) CreateAbsence(ctx context.Context, absence entity.Absence) (entity.Absence, error) {
	select {
	case <-ctx.Done():
		return entity.Absence{}, fmt.Errorf("database - CreateAbsence - ctx.Done: %w", ctx.Err())
	default:
		sql, args, errInsert := pg.Builder.
			Insert("absences").
			Columns("player_id", "raid_id").
			Values(absence.Player.ID, absence.Raid.ID).ToSql()
		if errInsert != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence - r.Builder: %w", errInsert)
		}
		_, err := pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence - r.Pool.Exec: %w", err)
		}
		return absence, nil
	}
}

func (pg *PG) ReadAbsence(ctx context.Context, absenceID int) (entity.Absence, error) {
	select {
	case <-ctx.Done():
		return entity.Absence{}, fmt.Errorf("database - ReadAbsence - ctx.Done: %w", ctx.Err())
	default:
		sql, _, err := pg.Builder.Select("id", "player_id", "raid_id").From("absences").Where("id = $1").ToSql()
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - ReadAbsence - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, absenceID)
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - ReadAbsence - r.Pool.Query: %w", err)
		}
		defer rows.Close()
		var absence entity.Absence
		if rows.Next() {
			err := rows.Scan(&absence.ID, &absence.Player.ID, &absence.Raid.ID)
			if err != nil {
				return entity.Absence{}, fmt.Errorf("database - ReadAbsence - rows.Scan: %w", err)
			}
			return absence, nil
		}
		return entity.Absence{}, fmt.Errorf("no absence found")
	}
}

func (pg *PG) UpdateAbsence(ctx context.Context, absence entity.Absence) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateAbsence - ctx.Done: %w", ctx.Err())
	default:
		sql, args, err := pg.Builder.
			Update("absences").
			Set("player_id", absence.Player.ID).
			Set("raid_id", absence.Raid.ID).
			Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdateAbsence - r.Builder: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("database - UpdateAbsence - r.Pool.Exec: %w", err)
		}
		return nil
	}
}

func (pg *PG) DeleteAbsence(ctx context.Context, absenceID int) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeleteAbsence - ctx.Done: %w", ctx.Err())
	default:
		sql, _, err := pg.Builder.Delete("absences").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeleteAbsence - r.Builder: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, absenceID)
		if err != nil {
			return fmt.Errorf("database - DeleteAbsence - r.Pool.Exec: %w", err)
		}
		return nil
	}
}
