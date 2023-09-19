package postgresbackend

import (
	"context"
	"fmt"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"go.uber.org/zap"
)

// searchStrikeOnID is a function which call backend to search a strike object based on a parameter.
func (pg *PG) searchStrikeOnParam(ctx context.Context, paramName string, param interface{}) ([]entity.Strike, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - ctx.Done: %w", ctx.Err())
	default:
		var strikes []entity.Strike
		sql, _, err := pg.Builder.
			Select("id", "player_id", "season", "reason", "created_at").
			From("strikes").Where("$1 = $2").ToSql()
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, paramName, param)
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
}

// SearchStrike is a function which call backend to Search a Strike Object
// It returns a list of strikes matching the given parameters not combined.
func (pg *PG) SearchStrike(
	ctx context.Context, playerID int, date time.Time, season, reason string,
) ([]entity.Strike, error) {
	zap.L().Debug("SearchStrike",
		zap.Int("playerID", playerID),
		zap.Time("date", date),
		zap.String("season", season),
		zap.String("reason", reason))
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchStrike - ctx.Done: %w", ctx.Err())
	default:
		var strikes []entity.Strike
		if playerID != -1 {
			s, err := pg.searchStrikeOnParam(ctx, "player_id", playerID)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(season) != 0 {
			s, err := pg.searchStrikeOnParam(ctx, "season", season)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(reason) != 0 {
			s, err := pg.searchStrikeOnParam(ctx, "reason", reason)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if !date.IsZero() {
			s, err := pg.searchStrikeOnParam(ctx, "created_at", date)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		return strikes, nil
	}
}

// CreateStrike is a function which call backend to Create a Strike Object.
func (pg *PG) CreateStrike(ctx context.Context, strike entity.Strike, player entity.Player) error {
	zap.L().Debug("CreateStrike",
		zap.String("season", strike.Season),
		zap.String("reason", strike.Reason),
		zap.String("player", player.Name))
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - CreateStrike - ctx.Done: %w", ctx.Err())
	default:
		sql, args, errInsert := pg.Builder.
			Insert("strikes").
			Columns("player_id", "season", "reason").
			Values(player.ID, strike.Season, strike.Reason).ToSql()
		if errInsert != nil {
			return fmt.Errorf("database - CreateStrike - r.Builder: %w", errInsert)
		}
		_, err := pg.Pool.Exec(ctx, sql, args...)
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
		return entity.Strike{}, fmt.Errorf("database - ReadStrike - ctx.Done: %w", ctx.Err())
	default:
		// Find Strike from id on database
		sql, _, err := pg.Builder.
			Select("id", "player_id", "season", "reason", "created_at").
			From("strikes").
			Where("id = $1").ToSql()
		if err != nil {
			return entity.Strike{}, fmt.Errorf("database - ReadStrike - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, strikeID)
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
	zap.L().Debug("UpdateStrike",
		zap.Int("strikeID", strike.ID),
		zap.String("season", strike.Season),
		zap.String("reason", strike.Reason))
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateStrike - ctx.Done: %w", ctx.Err())
	default:
		// Update strike on database
		sql, _, err := pg.Builder.Update("strikes").
			Set("season", strike.Season).
			Set("reason", strike.Reason).
			Where("id = $1").ToSql()
		if err != nil {
			return fmt.Errorf("database - UpdateStrike - r.Builder: %w", err)
		}
		_, err = pg.Pool.Exec(ctx, sql, strike.ID)
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
		return fmt.Errorf("database - DeleteStrike - ctx.Done: %w", ctx.Err())
	default:
		sql, _, errInsert := pg.Builder.Delete("strikes").Where("id = $1").ToSql()
		if errInsert != nil {
			return fmt.Errorf("database - DeleteStrike - r.Builder: %w", errInsert)
		}
		_, err := pg.Pool.Exec(ctx, sql, strikeID)
		if err != nil {
			return fmt.Errorf("database - DeleteStrike - r.Pool.Exec: %w", err)
		}
		return nil
	}
}
