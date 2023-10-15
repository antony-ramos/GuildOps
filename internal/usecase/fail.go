package usecase

import (
	"context"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/pkg/logger"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type FailUseCase struct {
	backend Backend
}

func NewFailUseCase(bk Backend) *FailUseCase {
	return &FailUseCase{backend: bk}
}

func (fuc FailUseCase) CreateFail(ctx context.Context, failReason string, date time.Time, playerName string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/CreateFail")
	span.SetAttributes(
		attribute.String("failReason", failReason),
		attribute.String("playerName", playerName),
		attribute.String("date", date.String()))
	defer span.End()
	logger.FromContext(ctx).Debug("create fail use case")

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "create fail")
	default:
		player, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return errors.Wrap(err, "create fail search player")
		}
		if len(player) == 0 {
			return errors.New("player not found")
		}

		raid, err := fuc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return errors.Wrap(err, "create fail search raid")
		}
		if len(raid) == 0 {
			return errors.New("raid not found")
		}

		fail, err := entity.NewFail(-1, failReason, &player[0], &raid[0])
		if err != nil {
			return errors.Wrap(err, "create a fail object")
		}

		_, err = fuc.backend.CreateFail(ctx, fail)
		if err != nil {
			return errors.Wrap(err, "create fail")
		}
		logger.FromContext(ctx).Debug("create fail use case success")
		return nil
	}
}

func (fuc FailUseCase) ListFailOnPLayer(ctx context.Context, playerName string) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/ListFailOnPLayer")
	span.SetAttributes(attribute.String("playerName", playerName))
	defer span.End()
	logger.FromContext(ctx).Debug("create fail use case")

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "list fail on player")
	default:
		player, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search player")
		}
		if len(player) == 0 {
			return nil, errors.New("player not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", player[0].ID, -1, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search fail")
		}
		for k, fail := range fails {
			r, err := fuc.backend.ReadRaid(ctx, fail.Raid.ID)
			if err != nil {
				return nil, errors.Wrap(err, "list fail on player read raid")
			}
			fails[k].Raid = &r
		}

		return fails, nil
	}
}

func (fuc FailUseCase) ListFailOnRaid(ctx context.Context, date time.Time) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/ListFailOnRaid")
	span.SetAttributes(attribute.String("date", date.String()))
	defer span.End()
	logger.FromContext(ctx).Debug("list fail on raid use case")

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "list fail on raid")
	default:
		raid, err := fuc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search raid")
		}
		if len(raid) == 0 {
			return nil, errors.New("raid not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", -1, raid[0].ID, "")
		for k, fail := range fails {
			p, err := fuc.backend.ReadPlayer(ctx, fail.Player.ID)
			if err != nil {
				return nil, errors.Wrap(err, "list fail on raid read player")
			}
			fails[k].Player = &p
		}

		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search fail")
		}
		return fails, nil
	}
}

func (fuc FailUseCase) ListFailOnRaidAndPlayer(
	ctx context.Context, raidName string, playerName string,
) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/ListFailOnRaidAndPlayer")
	span.SetAttributes(attribute.String("raidName", raidName), attribute.String("playerName", playerName))
	defer span.End()
	logger.FromContext(ctx).Debug("list fail on raid and player use case")

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "list fail on raid and player")
	default:
		player, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search player")
		}
		if len(player) == 0 {
			return nil, errors.New("player not found")
		}

		raid, err := fuc.backend.SearchRaid(ctx, raidName, time.Time{}, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search raid")
		}
		if len(raid) == 0 {
			return nil, errors.New("raid not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", player[0].ID, raid[0].ID, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search fail")
		}
		return fails, nil
	}
}

func (fuc FailUseCase) DeleteFail(ctx context.Context, failID int) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/DeleteFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()
	logger.FromContext(ctx).Debug("delete fail use case")

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "delete fail")
	default:
		err := fuc.backend.DeleteFail(ctx, failID)
		if err != nil {
			return errors.Wrap(err, "delete fail")
		}
		return nil
	}
}

func (fuc FailUseCase) UpdateFail(ctx context.Context, failID int, failReason string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/UpdateFail")
	span.SetAttributes(attribute.Int("failID", failID), attribute.String("failReason", failReason))
	defer span.End()
	logger.FromContext(ctx).Debug("update fail use case")

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "update fail")
	default:
		fail, err := fuc.backend.ReadFail(ctx, failID)
		if err != nil {
			return errors.Wrap(err, "update fail read fail")
		}

		fail, err = entity.NewFail(fail.ID, failReason, fail.Player, fail.Raid)
		if err != nil {
			return errors.Wrap(err, "create a new fail object")
		}

		err = fuc.backend.UpdateFail(ctx, fail)
		if err != nil {
			return errors.Wrap(err, "update fail")
		}
		return nil
	}
}

func (fuc FailUseCase) ReadFail(ctx context.Context, failID int) (entity.Fail, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Fail/ReadFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()
	logger.FromContext(ctx).Debug("read fail use case")

	select {
	case <-ctx.Done():
		return entity.Fail{}, errors.Wrap(ctx.Err(), "read fail")
	default:
		fail, err := fuc.backend.ReadFail(ctx, failID)
		if err != nil {
			return entity.Fail{}, errors.Wrap(err, "read fail")
		}
		return fail, nil
	}
}
