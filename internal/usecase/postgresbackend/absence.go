package postgresbackend

import (
	"context"
	"fmt"
	"time"


	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

func (pg *PG) searchAbsenceOnParam(ctx context.Context, paramName string, param interface{}) ([]entity.Absence, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/searchAbsenceOnParam")
	defer span.End()
	span.SetAttributes(
		attribute.String("paramName", paramName),
		attribute.String("param", fmt.Sprintf("%v", param)),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchAbsence - searchAbsenceOnParam - " +
			"ctx.Done: request took too much time to be proceed")
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/SearchAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.Int("playerID", playerID),
		attribute.String("date", date.Format("02/01/2006")),
	)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchAbsence - ctx.Done: request took too much time to be proceed")
	default:
		var absences []entity.Absence
		switch {
		case playerID != -1 && !date.IsZero():
			a, err := pg.searchAbsenceOnParam(ctx, "date", date)
			if err != nil {
				return nil, errors.Wrap(err, "database - SearchAbsence - searchAbsenceOnParam")
			}
			var absencesOnPlayer []entity.Absence
			for _, absence := range a {
				if absence.Player.ID == playerID {
					absencesOnPlayer = append(absencesOnPlayer, absence)
				}
			}
			return append(absences, absencesOnPlayer...), nil
		case playerID != -1 && playerName == "":
			a, err := pg.searchAbsenceOnParam(ctx, "player_id", playerID)
			if err != nil {
				return nil, err
			}
			return append(absences, a...), nil

		case playerID == -1 && playerName != "":
			a, err := pg.searchAbsenceOnParam(ctx, "players.name", playerName)
			if err != nil {
				return nil, err
			}
			return append(absences, a...), nil

		case playerID == -1 && playerName == "" && !date.IsZero():
			a, err := pg.searchAbsenceOnParam(ctx, "date", date)
			if err != nil {
				return nil, err
			}
			return append(absences, a...), nil
		}
		return nil, errors.New("database - SearchAbsence: no param given")
	}
}

func (pg *PG) CreateAbsence(ctx context.Context, absence entity.Absence) (entity.Absence, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/CreateAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.Int("absenceID", absence.ID),
		attribute.Int("playerID", absence.Player.ID),
		attribute.Int("raidID", absence.Raid.ID),
	)
	select {
	case <-ctx.Done():
		return entity.Absence{}, fmt.Errorf("database - CreateAbsence:  ctx.Done: request took too much time to be proceed")
	default:
		// Search if the absence already exists
		sql, _, err := pg.Builder.
			Select("id", "player_id", "raid_id").
			From("absences").
			Where("player_id = $1 AND raid_id = $2").ToSql()
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence:  r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, absence.Player.ID, absence.Raid.ID)
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence:  r.Pool.Query: %w", err)
		}
		defer rows.Close()
		if rows.Next() {
			return absence, fmt.Errorf("absence already exists")
		}

		sql, args, errInsert := pg.Builder.
			Insert("absences").
			Columns("player_id", "raid_id").
			Values(absence.Player.ID, absence.Raid.ID).ToSql()
		if errInsert != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence:  r.Builder: %w", errInsert)
		}
		_, err = pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return entity.Absence{}, fmt.Errorf("database - CreateAbsence:  r.Pool.Exec: %w", err)
		}
		return absence, nil
	}
}

func (pg *PG) ReadAbsence(ctx context.Context, absenceID int) (entity.Absence, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/ReadAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.Int("absenceID", absenceID),
	)
	select {
	case <-ctx.Done():
		return entity.Absence{}, fmt.Errorf("database - ReadAbsence - ctx.Done: request took too much time to be proceed")
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
		return entity.Absence{}, fmt.Errorf("absence not found")
	}
}

func (pg *PG) UpdateAbsence(ctx context.Context, absence entity.Absence) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/UpdateAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.Int("absenceID", absence.ID),
		attribute.Int("playerID", absence.Player.ID),
		attribute.Int("raidID", absence.Raid.ID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateAbsence - ctx.Done: request took too much time to be proceed")
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Absence/DeleteAbsence")
	defer span.End()
	span.SetAttributes(
		attribute.Int("absenceID", absenceID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - DeleteAbsence - ctx.Done: request took too much time to be proceed")
	default:
		sql, _, err := pg.Builder.Delete("absences").Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - DeleteAbsence - r.Builder: %w", err)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, absenceID)
		if err != nil {
			return fmt.Errorf("database - DeleteAbsence - r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database - DeleteAbsence - absence not found")
		}
		return nil
	}
}
