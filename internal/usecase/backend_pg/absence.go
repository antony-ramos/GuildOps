package backend_pg

import (
	"context"
	"fmt"
	"time"

	"github.com/coven-discord-bot/internal/entity"
)

func (pg *PG) SearchAbsence(ctx context.Context, playerName string, playerID int, date time.Time) ([]entity.Absence, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var absences []entity.Absence
		switch {
		case playerID == -1 && playerName == "":
			sql, _, err := pg.Builder.Select("absences.id", "absences.player_id", "absences.raid_id", "raids.name", "raids.difficulty", "raids.date", "players.name").From("absences").Join("raids ON raids.id = absences.raid_id").Join("players ON players.id = absences.player_id").Where("raids.date = $1").ToSql()
			if err != nil {
				return nil, err
			}
			rows, err := pg.Pool.Query(context.Background(), sql, date)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var absence entity.Absence
				var raid entity.Raid
				var player entity.Player
				err := rows.Scan(&absence.ID, &player.ID, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.Name)
				if err != nil {
					return nil, err
				}
				absence.Player = &player
				absence.Raid = &raid
				absences = append(absences, absence)
			}
		case playerID != -1 && playerName == "":
			sql, _, err := pg.Builder.Select("absences.id", "absences.player_id", "absences.raid_id", "raids.name", "raids.difficulty", "raids.date", "players.name").From("absences").Join("raids ON raids.id = absences.raid_id").Join("players ON players.id = absences.player_id").Where("absences.player_id = $1").Where("raids.date = $2").ToSql()
			if err != nil {
				return nil, err
			}
			rows, err := pg.Pool.Query(context.Background(), sql, playerID, date)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var absence entity.Absence
				var raid entity.Raid
				var player entity.Player
				err := rows.Scan(&absence.ID, &player.ID, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.Name)
				if err != nil {
					return nil, err
				}
				absence.Player = &player
				absence.Raid = &raid
				absences = append(absences, absence)
			}
		case playerID == -1 && playerName != "":
			sql, _, err := pg.Builder.Select("absences.id", "absences.player_id", "absences.raid_id", "raids.name", "raids.difficulty", "raids.date", "players.name").From("absences").Join("raids ON raids.id = absences.raid_id").Join("players ON players.id = absences.player_id").Where("players.name = $1").Where("raids.date = $2").ToSql()
			if err != nil {
				return nil, err
			}
			rows, err := pg.Pool.Query(context.Background(), sql, playerName, date)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var absence entity.Absence
				var raid entity.Raid
				var player entity.Player
				err := rows.Scan(&absence.ID, &player.ID, &raid.ID, &raid.Name, &raid.Difficulty, &raid.Date, &player.Name)
				if err != nil {
					return nil, err
				}
				absence.Player = &player
				absence.Raid = &raid
				absences = append(absences, absence)
			}
		case playerID != -1 && playerName != "":
			return nil, fmt.Errorf("database - SearchAbsence - cannot have both playerID and playerName")
		}

		return absences, nil
	}
}

func (pg *PG) CreateAbsence(ctx context.Context, absence entity.Absence) (entity.Absence, error) {
	select {
	case <-ctx.Done():
		return entity.Absence{}, ctx.Err()
	default:
		sql, args, errInsert := pg.Builder.Insert("absences").Columns("player_id", "raid_id").Values(absence.Player.ID, absence.Raid.ID).ToSql()
		if errInsert != nil {
			return entity.Absence{}, errInsert
		}
		_, err := pg.Pool.Exec(context.Background(), sql, args...)
		if err != nil {
			return entity.Absence{}, err
		}
		return absence, nil
	}
}

func (pg *PG) ReadAbsence(ctx context.Context, absenceID int) (entity.Absence, error) {
	select {
	case <-ctx.Done():
		return entity.Absence{}, ctx.Err()
	default:
		sql, _, err := pg.Builder.Select("id", "player_id", "raid_id").From("absences").Where("id = $1").ToSql()
		if err != nil {
			return entity.Absence{}, err
		}
		rows, err := pg.Pool.Query(context.Background(), sql, absenceID)
		if err != nil {
			return entity.Absence{}, err
		}
		defer rows.Close()
		var absence entity.Absence
		if rows.Next() {
			err := rows.Scan(&absence.ID, &absence.Player.ID, &absence.Raid.ID)
			if err != nil {
				return entity.Absence{}, err
			}
			return absence, nil
		}
		return entity.Absence{}, err
	}
}

func (pg *PG) UpdateAbsence(ctx context.Context, absence entity.Absence) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		sql, args, err := pg.Builder.Update("absences").Set("player_id", absence.Player.ID).Set("raid_id", absence.Raid.ID).Where("id = $1").ToSql()
		if err != nil {
			return err
		}
		_, err = pg.Pool.Exec(context.Background(), sql, args...)
		if err != nil {
			return err
		}
		return nil
	}
}

func (pg *PG) DeleteAbsence(ctx context.Context, absenceID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		sql, _, err := pg.Builder.Delete("absences").Where("id = $1").ToSql()
		if err != nil {
			return err
		}
		_, err = pg.Pool.Exec(context.Background(), sql, absenceID)
		if err != nil {
			return err
		}
		return nil
	}
}
