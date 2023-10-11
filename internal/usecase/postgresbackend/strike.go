package postgresbackend

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/pkg/logger"

	"go.uber.org/zap"
)

// SearchStrikeOnParam is a function which call backend to Search a Strike Object on a given parameter.
func (pg *PG) SearchStrikeOnParam(ctx context.Context, paramName string, param interface{}) ([]entity.Strike, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/SearchStrikeOnParam")
	span.SetAttributes(
		attribute.String("paramName", paramName),
		attribute.String("param", fmt.Sprintf("%v", param)))
	defer span.End()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - " +
			"ctx.Done: request took too much time to be proceed")
	default:
		var strikes []entity.Strike
		condition := fmt.Sprintf("%s = $1", paramName)
		sql, _, err := pg.Builder.
			Select("id", "player_id", "season", "reason", "created_at").
			From("strikes").Where(condition).ToSql()
		if err != nil {
			return nil, fmt.Errorf("database - SearchStrike - searchStrikeOnID - r.Builder: %w", err)
		}
		rows, err := pg.Pool.Query(ctx, sql, param)
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/SearchStrike")
	span.SetAttributes(
		attribute.Int("playerID", playerID),
		attribute.String("season", season),
		attribute.String("reason", reason),
		attribute.String("date", date.String()))

	defer span.End()
	logger.FromContext(ctx).Debug("SearchStrike",
		zap.Int("playerID", playerID),
		zap.Time("date", date),
		zap.String("season", season),
		zap.String("reason", reason))
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("database - SearchStrike - ctx.Done: request took too much time to be proceed")
	default:
		var strikes []entity.Strike
		if playerID != -1 {
			s, err := pg.SearchStrikeOnParam(ctx, "player_id", playerID)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(season) != 0 {
			s, err := pg.SearchStrikeOnParam(ctx, "season", season)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if len(reason) != 0 {
			s, err := pg.SearchStrikeOnParam(ctx, "reason", reason)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		if !date.IsZero() {
			s, err := pg.SearchStrikeOnParam(ctx, "created_at", date)
			if err != nil {
				return nil, err
			}
			strikes = append(strikes, s...)
		}
		return strikes, nil
	}
}

// CreateStrike is a function which call backend to Create a Strike Object.
func (pg *PG) CreateStrike(ctx context.Context, strike entity.Strike, playerID int) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/CreateStrike")
	span.SetAttributes(
		attribute.String("season", strike.Season),
		attribute.String("reason", strike.Reason),
		attribute.Int("playerID", playerID))

	defer span.End()

	logger.FromContext(ctx).Debug("CreateStrike",
		zap.String("season", strike.Season),
		zap.String("reason", strike.Reason),
		zap.Int("player", playerID))
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - CreateStrike - ctx.Done: request took too much time to be proceed")
	default:
		sql, args, errInsert := pg.Builder.
			Insert("strikes").
			Columns("player_id", "season", "reason").
			Values(playerID, strike.Season, strike.Reason).ToSql()
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/ReadStrike")
	span.SetAttributes(
		attribute.Int("strikeID", strikeID))

	defer span.End()

	logger.FromContext(ctx).Debug("ReadStrike", zap.Int("strikeID", strikeID))
	select {
	case <-ctx.Done():
		return entity.Strike{}, fmt.Errorf("database - ReadStrike - ctx.Done: request took too much time to be proceed")
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/UpdateStrike")
	span.SetAttributes(
		attribute.Int("strikeID", strike.ID),
		attribute.String("season", strike.Season),
		attribute.String("reason", strike.Reason))

	defer span.End()

	logger.FromContext(ctx).Debug("UpdateStrike",
		zap.Int("strikeID", strike.ID),
		zap.String("season", strike.Season),
		zap.String("reason", strike.Reason))
	select {
	case <-ctx.Done():
		return fmt.Errorf("database - UpdateStrike - ctx.Done: request took too much time to be proceed")
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
	ctx, span := otel.Tracer("Backend").Start(ctx, "Strike/DeleteStrike")
	span.SetAttributes(
		attribute.Int("strikeID", strikeID))
	defer span.End()

	logger.FromContext(ctx).Debug("DeleteStrike", zap.Int("strikeID", strikeID))
	select {
	case <-ctx.Done():
		return fmt.Errorf("database DeleteStrike: ctx.Done: request took too much time to be proceed")
	default:
		sql, _, errInsert := pg.Builder.Delete("strikes").Where("id = $1").ToSql()
		if errInsert != nil {
			return fmt.Errorf("database DeleteStrike: r.Builder: %w", errInsert)
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, strikeID)
		if err != nil {
			return fmt.Errorf("database DeleteStrike: r.Pool.Exec: %w", err)
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("database DeleteStrike: strike not found")
		}
		return nil
	}
}
