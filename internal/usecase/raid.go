package usecase

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/antony-ramos/guildops/internal/entity"
)

type RaidUseCase struct {
	backend Backend
}

func NewRaidUseCase(bk Backend) *RaidUseCase {
	return &RaidUseCase{backend: bk}
}

func (puc RaidUseCase) CreateRaid(
	ctx context.Context, raidName, difficulty string, date time.Time,
) (entity.Raid, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Raid/CreateRaid")
	defer span.End()
	span.SetAttributes(
		attribute.String("raidName", raidName),
		attribute.String("difficulty", difficulty),
		attribute.String("date", date.String()),
	)

	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("RaidUseCase - CreateRaid - ctx.Done: request took too much time to be proceed")
	default:
		raid, err := entity.NewRaid(raidName, difficulty, date)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("create entity raid for backend: %w", err)
		}
		raid, err = puc.backend.CreateRaid(ctx, raid)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.CreateRaid: %w", err)
		}
		return raid, nil
	}
}

func (puc RaidUseCase) DeleteRaidWithID(ctx context.Context, raidID int) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Raid/DeleteRaidWithID")
	defer span.End()
	span.SetAttributes(
		attribute.Int("raidID", raidID),
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("RaidUseCase - DeleteRaid - ctx.Done: request took too much time to be proceed")
	default:
		err := puc.backend.DeleteRaid(ctx, raidID)
		if err != nil {
			return fmt.Errorf("database - DeleteRaid - r.DeleteRaid: %w", err)
		}
		return nil
	}
}

func (puc RaidUseCase) DeleteRaidOnDate(ctx context.Context, date time.Time, difficulty string) error {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Raid/DeleteRaidOnDate")
	defer span.End()
	span.SetAttributes(
		attribute.String("difficulty", difficulty),
		attribute.String("date", date.String()),
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("check if context is valid: %w", ctx.Err())
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", date, difficulty)
		if err != nil {
			return fmt.Errorf("search raid with this date and difficulty combination: %w", err)
		}
		if len(raids) == 0 {
			return fmt.Errorf("check if there is a raid with this date/difficulty combination: %w", err)
		}
		err = puc.backend.DeleteRaid(ctx, raids[0].ID)
		if err != nil {
			return fmt.Errorf("delete raid previously found with date/difficulty combination: %w", err)
		}
	}
	return nil
}

func (puc RaidUseCase) ReadRaid(ctx context.Context, date time.Time) (entity.Raid, error) {
	ctx, span := otel.Tracer("Usecase").Start(ctx, "Raid/ReadRaid")
	span.SetAttributes(
		attribute.String("date", date.Format("02/01/06")),
	)
	defer span.End()

	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("RaidUseCase - DeleteRaid - ctx.Done: request took too much time to be proceed")
	default:
		raids, err := puc.backend.SearchRaid(ctx, "", date, "")
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - DeleteRaid - r.DeleteRaid: %w", err)
		}
		if len(raids) == 0 {
			return entity.Raid{}, fmt.Errorf("database - DeleteRaid - r.DeleteRaid: %w", err)
		}
		return raids[0], nil
	}
}
