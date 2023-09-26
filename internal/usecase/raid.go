package usecase

import (
	"context"
	"fmt"
	"time"

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
