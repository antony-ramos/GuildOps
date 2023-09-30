package usecase

import (
	"context"
	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/pkg/logger"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type FailUseCase struct {
	backend Backend
}

func NewFailUseCase(bk Backend) *FailUseCase {
	return &FailUseCase{backend: bk}
}

func (fuc FailUseCase) CreateFail(ctx context.Context, failReason string, date time.Time, playerName string) error {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/CreateFail")
	span.SetAttributes(attribute.String("failReason", failReason), attribute.String("playerName", playerName), attribute.String("date", date.String()))
	defer span.End()
	logger.FromContext(ctx).Debug("create fail use case")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		p, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return errors.Wrap(err, "create fail search player")
		}
		if len(p) == 0 {
			return errors.New("player not found")
		}

		r, err := fuc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return errors.Wrap(err, "create fail search raid")
		}
		if len(r) == 0 {
			return errors.New("raid not found")
		}

		fail := entity.Fail{
			Reason: failReason,
			Player: &p[0],
			Raid:   &r[0],
		}
		err = fail.Validate()
		if err != nil {
			return errors.Wrap(err, "create fail validate")
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
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/ListFailOnPLayer")
	span.SetAttributes(attribute.String("playerName", playerName))
	defer span.End()
	logger.FromContext(ctx).Debug("create fail use case")

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search player")
		}
		if len(p) == 0 {
			return nil, errors.New("player not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", p[0].ID, -1, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search fail")
		}
		return fails, nil
	}
}

func (fuc FailUseCase) ListFailOnRaid(ctx context.Context, raidName string) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/ListFailOnRaid")
	span.SetAttributes(attribute.String("raidName", raidName))
	defer span.End()
	logger.FromContext(ctx).Debug("list fail on raid use case")

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

		r, err := fuc.backend.SearchRaid(ctx, raidName, time.Time{}, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search raid")
		}
		if len(r) == 0 {
			return nil, errors.New("raid not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", -1, r[0].ID, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search fail")
		}
		return fails, nil
	}
}

func (fuc FailUseCase) ListFailOnRaidAndPlayer(ctx context.Context, raidName string, playerName string) ([]entity.Fail, error) {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/ListFailOnRaidAndPlayer")
	span.SetAttributes(attribute.String("raidName", raidName), attribute.String("playerName", playerName))
	defer span.End()
	logger.FromContext(ctx).Debug("list fail on raid and player use case")

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p, err := fuc.backend.SearchPlayer(ctx, -1, playerName, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on player search player")
		}
		if len(p) == 0 {
			return nil, errors.New("player not found")
		}

		r, err := fuc.backend.SearchRaid(ctx, raidName, time.Time{}, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search raid")
		}
		if len(r) == 0 {
			return nil, errors.New("raid not found")
		}

		fails, err := fuc.backend.SearchFail(ctx, "", p[0].ID, r[0].ID, "")
		if err != nil {
			return nil, errors.Wrap(err, "list fail on raid search fail")
		}
		return fails, nil
	}
}

func (fuc FailUseCase) DeleteFail(ctx context.Context, failID int) error {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/DeleteFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()
	logger.FromContext(ctx).Debug("delete fail use case")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := fuc.backend.DeleteFail(ctx, failID)
		if err != nil {
			return errors.Wrap(err, "delete fail")
		}
		return nil
	}
}

func (fuc FailUseCase) UpdateFail(ctx context.Context, failID int, failReason string) error {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/UpdateFail")
	span.SetAttributes(attribute.Int("failID", failID), attribute.String("failReason", failReason))
	defer span.End()
	logger.FromContext(ctx).Debug("update fail use case")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		fail, err := fuc.backend.ReadFail(ctx, failID)
		fail.Reason = failReason
		err = fail.Validate()
		if err != nil {
			return errors.Wrap(err, "update fail validate")
		}

		err = fuc.backend.UpdateFail(ctx, fail)
		if err != nil {
			return errors.Wrap(err, "update fail")
		}
		return nil
	}
}

func (fuc FailUseCase) ReadFail(ctx context.Context, failID int) (entity.Fail, error) {
	ctx, span := otel.Tracer("UseCase").Start(ctx, "FailUseCase/ReadFail")
	span.SetAttributes(attribute.Int("failID", failID))
	defer span.End()
	logger.FromContext(ctx).Debug("read fail use case")

	select {
	case <-ctx.Done():
		return entity.Fail{}, ctx.Err()
	default:
		fail, err := fuc.backend.ReadFail(ctx, failID)
		if err != nil {
			return entity.Fail{}, errors.Wrap(err, "read fail")

		}
		return fail, nil
	}
}
