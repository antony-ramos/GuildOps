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
	select {
	case <-ctx.Done():
		return entity.Raid{}, fmt.Errorf("RaidUseCase - CreateRaid - ctx.Done: request took too much time to be proceed")
	default:
		raid := entity.Raid{
			Name:       raidName,
			Difficulty: difficulty,
			Date:       date,
		}
		err := raid.Validate()
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.Validate: %w", err)
		}
		raid, err = puc.backend.CreateRaid(ctx, raid)
		if err != nil {
			return entity.Raid{}, fmt.Errorf("database - CreateRaid - r.CreateRaid: %w", err)
		}
		return raid, nil
	}
}

func (puc RaidUseCase) DeleteRaid(ctx context.Context, raidID int) error {
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

func (puc RaidUseCase) ReadRaid(ctx context.Context, date time.Time) (entity.Raid, error) {
	ctx, span := otel.Tracer("usecase").Start(ctx, "ReadRaid")
	span.SetAttributes(
		attribute.String("date", date.String()),
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
