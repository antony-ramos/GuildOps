package postgresbackend

import (
	"context"
	"fmt"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// SearchFailOnParam is a function which call backend to Search an entity.Fail on a given parameter.
func (pg *PG) SearchFailOnParam(ctx context.Context, paramName string, param interface{}) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/SearchFailOnParam")
	span.SetAttributes(
		attribute.String("paramName", paramName),
		attribute.String("param", fmt.Sprintf("%v", param)))
	defer span.End()

	logger.FromContext(ctx).Debug("SearchFailOnParam",
		zap.String("paramName", paramName),
		zap.String("param", fmt.Sprintf("%v", param)))

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "search fail from pg database with param")
	default:
		var fails []entity.Fail
		condition := fmt.Sprintf("%s = $1", paramName)
		sql, _, err := pg.Builder.
			Select("id", "player_id", "raid_id", "reason").
			From("fails").Where(condition).ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "create query to search fail with param")
		}
		rows, err := pg.Pool.Query(ctx, sql, param)
		if err != nil {
			return nil, errors.Wrap(err, "exec query to search fail with param")
		}
		defer rows.Close()
		for rows.Next() {
			var fail entity.Fail
			err := rows.Scan(&fail.ID, nil, nil, &fail.Reason)
			if err != nil {
				return nil, errors.Wrap(err, "scan query to search fail with param")
			}
			fails = append(fails, fail)
		}
		return fails, nil
	}
}

// SearchFail is a function which call backend to Search an entity.Fail.
func (pg *PG) SearchFail(
	ctx context.Context, playerName string, playerID int, raidID int, reason string,
) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/SearchFail")
	span.SetAttributes(
		attribute.String("playerName", playerName),
		attribute.Int("playerID", playerID),
		attribute.Int("raidID", raidID),
		attribute.String("reason", reason))
	defer span.End()

	logger.FromContext(ctx).Debug("SearchFail",
		zap.String("playerName", playerName),
		zap.Int("playerID", playerID),
		zap.Int("raidID", raidID),
		zap.String("reason", reason))

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "search fail from pg database")

	default:
		var fails []entity.Fail
		if playerID != -1 {
			s, err := pg.SearchFailOnParam(ctx, "player_id", playerID)
			if err != nil {
				return nil, err
			}
			fails = append(fails, s...)
		}
		if raidID != -1 {
			s, err := pg.SearchFailOnParam(ctx, "raid_ID", playerID)
			if err != nil {
				return nil, err
			}
			fails = append(fails, s...)
		}
		if len(reason) != 0 {
			s, err := pg.SearchFailOnParam(ctx, "reason", reason)
			if err != nil {
				return nil, err
			}
			fails = append(fails, s...)
		}
		return fails, nil
	}
}

// CreateFail create an entity.Fail in database.
func (pg *PG) CreateFail(ctx context.Context, fail entity.Fail) (entity.Fail, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/CreateFail")
	span.SetAttributes(
		attribute.String("failReason", fail.Reason),
		attribute.String("playerName", fail.Player.Name),
		attribute.String("date", fail.Raid.Date.String()))
	defer span.End()

	logger.FromContext(ctx).Debug("CreateFail",
		zap.String("failReason", fail.Reason),
		zap.String("playerName", fail.Player.Name),
		zap.String("date", fail.Raid.Date.String()))

	select {
	case <-ctx.Done():
		return entity.Fail{}, errors.Wrap(ctx.Err(), "create fail from pg database")
	default:
		sql, args, errInsert := pg.Builder.
			Insert("fails").
			Columns("player_id", "raid_id", "reason").
			Values(fail.Player.ID, fail.Raid.ID, fail.Reason).
			ToSql()
		if errInsert != nil {
			return entity.Fail{}, errors.Wrap(errInsert, "create query to create fail")
		}
		_, err := pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return entity.Fail{}, errors.Wrap(err, "exec query to create fail")
		}
		return entity.Fail{}, nil
	}
}

// ReadFail read an entity.Fail from database.
func (pg *PG) ReadFail(ctx context.Context, failID int) (entity.Fail, error) {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/ReadFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()
	logger.FromContext(ctx).Debug("ReadFail", zap.Int("failID", failID))
	select {
	case <-ctx.Done():
		return entity.Fail{}, errors.Wrap(ctx.Err(), "read fail from pg database")
	default:
		sql, _, err := pg.Builder.
			Select("id", "player_id", "raid_id", "reason").
			From("fails").
			Where("id = $1").ToSql()
		if err != nil {
			return entity.Fail{}, errors.Wrap(err, "create query to read fail")
		}
		rows, err := pg.Pool.Query(ctx, sql, failID)
		if err != nil {
			return entity.Fail{}, errors.Wrap(err, "exec query to read fail")
		}
		defer rows.Close()
		var fail entity.Fail
		for rows.Next() {
			err := rows.Scan(&fail.ID, nil, nil, &fail.Reason)
			if err != nil {
				return entity.Fail{}, fmt.Errorf("database - ReadFail - rows.Scan: %w", err)
			}
		}
		return fail, nil
	}
}

// UpdateFail update an entity.Fail from database.
func (pg *PG) UpdateFail(ctx context.Context, fail entity.Fail) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/UpdateFail")
	span.SetAttributes(attribute.Int("failID", fail.ID), attribute.String("failReason", fail.Reason))
	defer span.End()

	logger.FromContext(ctx).Debug(" UpdateFail",
		zap.Int("failID", fail.ID),
		zap.String("reason", fail.Reason))
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "update fail from pg database")
	default:
		sql, _, err := pg.Builder.Update("fails").
			Set("reason", fail.Reason).
			Where("id = $2").ToSql()
		if err != nil {
			return errors.Wrap(err, "create query to update fail")
		}
		_, err = pg.Pool.Exec(ctx, sql, fail.Reason, fail.ID)
		if err != nil {
			return errors.Wrap(err, "exec query to update fail")
		}
		return nil
	}
}

// DeleteFail delete an entity.Fail from database based on entity.Fail ID.
func (pg *PG) DeleteFail(ctx context.Context, failID int) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "FailBackend/DeleteFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()

	logger.FromContext(ctx).Debug("DeleteFail", zap.Int("failID", failID))
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "delete fail from pg database")
	default:
		sql, _, errInsert := pg.Builder.Delete("fails").Where("id = $1").ToSql()
		if errInsert != nil {
			return errors.Wrap(errInsert, "create query to delete fail")
		}
		isDelete, err := pg.Pool.Exec(ctx, sql, failID)
		if err != nil {
			return errors.Wrap(err, "exec query to delete fail")
		}
		if isDelete.String() == isNotDeleted {
			return fmt.Errorf("fail not found from pg database")
		}
		return nil
	}
}
