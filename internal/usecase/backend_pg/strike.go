package backend_pg

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"time"
)

// SearchStrike is a function which call backend to Search a Strike Object
func (pg *PG) SearchStrike(ctx context.Context, playerID int, Date time.Time, Season, Reason string) ([]entity.Strike, error) {
	var strikes []entity.Strike
	if playerID != -1 {
		sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("player_id = $1").ToSql()
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(context.Background(), sql, playerID)
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - r.Pool.Query: %w", err)

		}
		defer rows.Close()
		for rows.Next() {
			var strike entity.Strike
			err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
			if err != nil {
				return nil, fmt.Errorf("database - SearchStrike - rows.Scan: %w", err)
			}
			strikes = append(strikes, strike)
		}
	}

	return strikes, nil

}

// CreateStrike is a function which call backend to Create a Strike Object
func (pg *PG) CreateStrike(ctx context.Context, strike entity.Strike, player entity.Player) error {
	sql, args, errInsert := pg.Builder.Insert("strikes").Columns("player_id", "season", "reason").Values(player.ID, strike.Season, strike.Reason).ToSql()
	if errInsert != nil {
		return fmt.Errorf("database - CreateStrike - r.Builder: %w", errInsert)
	}
	_, err := pg.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("database - CreateStrike - r.Pool.Exec: %w", err)
	}
	return nil
}

func (pg *PG) ReadStrike(ctx context.Context, strikeID int) (entity.Strike, error) {
	// Find Strike from id on database
	sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("id = $1").ToSql()
	if err != nil {
		return entity.Strike{}, fmt.Errorf("database - ReadStrike - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, strikeID)
	if err != nil {
		return entity.Strike{}, fmt.Errorf("database - ReadStrike - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	var strike entity.Strike
	for rows.Next() {
		err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
		if err != nil {
			return entity.Strike{}, fmt.Errorf("database - ReadStrike - rows.Scan: %w", err)
		}

	}
	return strike, nil
}

func (pg *PG) UpdateStrike(ctx context.Context, strike entity.Strike) error {
	// Update strike on database
	sql, _, err := pg.Builder.Update("strikes").Set("season", strike.Season).Set("reason", strike.Reason).Where("id = $1").ToSql()
	if err != nil {
		return fmt.Errorf("database - UpdateStrike - r.Builder: %w", err)
	}
	_, err = pg.Pool.Exec(context.Background(), sql, strike.ID)
	if err != nil {
		return fmt.Errorf("database - UpdateStrike - r.Pool.Exec: %w", err)
	}
	return nil
}

func (pg *PG) DeleteStrike(ctx context.Context, strikeID int) error {
	sql, _, errInsert := pg.Builder.Delete("strikes").Where("id = $1").ToSql()
	if errInsert != nil {
		return fmt.Errorf("database - DeleteStrike - r.Builder: %w", errInsert)
	}
	_, err := pg.Pool.Exec(context.Background(), sql, strikeID)
	if err != nil {
		return fmt.Errorf("database - DeleteStrike - r.Pool.Exec: %w", err)
	}
	return nil
}
