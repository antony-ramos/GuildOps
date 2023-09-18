package backend_pg

import (
	"context"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"go.uber.org/zap"
	"time"
)

// SearchStrikeOnID is a function which call backend to search a strike object based on playerID
func (pg *PG) searchStrikeOnID(playerID int) ([]entity.Strike, error) {
	var strikes []entity.Strike
	sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("player_id = $1").ToSql()
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, playerID)
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var strike entity.Strike
		err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - rows.Scan: %w", err)
		}
		strikes = append(strikes, strike)
	}
	return strikes, nil
}

// SearchStrikeOnSeason is a function which call backend to search a strike object based on season
func (pg *PG) searchStrikeOnSeason(season string) ([]entity.Strike, error) {
	var strikes []entity.Strike
	sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("season = $1").ToSql()
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnSeason - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, season)
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnSeason - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var strike entity.Strike
		err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnSeason - rows.Scan: %w", err)
		}
		strikes = append(strikes, strike)
	}
	return strikes, nil
}

// SearchStrikeOnReason is a function which call backend to search a strike object based on reason
func (pg *PG) searchStrikeOnReason(reason string) ([]entity.Strike, error) {
	var strikes []entity.Strike
	sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("reason = $1").ToSql()
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnReason - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, reason)
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnReason - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var strike entity.Strike
		err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnReason - rows.Scan: %w", err)
		}
		strikes = append(strikes, strike)
	}
	return strikes, nil
}

// SearchStrikeOnDate is a function which call backend to search a strike object based on date
func (pg *PG) searchStrikeOnDate(date time.Time) ([]entity.Strike, error) {
	var strikes []entity.Strike
	sql, _, err := pg.Builder.Select("id", "player_id", "season", "reason", "created_at").From("strikes").Where("created_at = $1").ToSql()
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnDate - r.Builder: %w", err)
	}
	rows, err := pg.Pool.Query(context.Background(), sql, date)
	if err != nil {
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnDate - r.Pool.Query: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var strike entity.Strike
		err := rows.Scan(&strike.ID, nil, &strike.Season, &strike.Reason, &strike.Date)
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnDate - rows.Scan: %w", err)
		}
		strikes = append(strikes, strike)
	}
	return strikes, nil
}

// SearchStrike is a function which call backend to Search a Strike Object
// It returns a list of strikes matching the given parameters not combined
func (pg *PG) SearchStrike(ctx context.Context, playerID int, date time.Time, season, reason string) ([]entity.Strike, error) {
	zap.L().Debug("SearchStrike", zap.Int("playerID", playerID), zap.Time("date", date), zap.String("season", season), zap.String("reason", reason))
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var strikes []entity.Strike
		if playerID != -1 {
			s, err := pg.searchStrikeOnID(playerID)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(season) != 0 {
			s, err := pg.searchStrikeOnSeason(season)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(reason) != 0 {
			s, err := pg.searchStrikeOnReason(reason)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if !date.IsZero() {
			s, err := pg.searchStrikeOnDate(date)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		return strikes, nil
	}
}

// CreateStrike is a function which call backend to Create a Strike Object
func (pg *PG) CreateStrike(ctx context.Context, strike entity.Strike, player entity.Player) error {
	zap.L().Debug("CreateStrike", zap.String("season", strike.Season), zap.String("reason", strike.Reason), zap.String("player", player.Name))
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
}

func (pg *PG) ReadStrike(ctx context.Context, strikeID int) (entity.Strike, error) {
	zap.L().Debug("ReadStrike", zap.Int("strikeID", strikeID))
	select {
	case <-ctx.Done():
		return entity.Strike{}, ctx.Err()
	default:
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
}

func (pg *PG) UpdateStrike(ctx context.Context, strike entity.Strike) error {
	zap.L().Debug("UpdateStrike", zap.Int("strikeID", strike.ID), zap.String("season", strike.Season), zap.String("reason", strike.Reason))
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
}

func (pg *PG) DeleteStrike(ctx context.Context, strikeID int) error {
	zap.L().Debug("DeleteStrike", zap.Int("strikeID", strikeID))
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
}
