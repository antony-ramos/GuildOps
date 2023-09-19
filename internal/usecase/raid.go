package usecase

import (
	"context"
	"github.com/antony-ramos/guildops/internal/entity"
	"time"
)

type RaidUseCase struct {
	backend Backend
}

func NewRaidUseCase(bk Backend) *RaidUseCase {
	return &RaidUseCase{backend: bk}
}

func (puc RaidUseCase) CreateRaid(ctx context.Context, raidName, difficulty string, date time.Time) (entity.Raid, error) {
	raid := entity.Raid{
		Name:       raidName,
		Difficulty: difficulty,
		Date:       date,
	}
	err := raid.Validate()
	if err != nil {
		return entity.Raid{}, err
	}
	raid, err = puc.backend.CreateRaid(ctx, raid)
	if err != nil {
		return entity.Raid{}, err
	}
	return raid, nil
}

func (puc RaidUseCase) DeleteRaid(ctx context.Context, raidID int) error {
	err := puc.backend.DeleteRaid(ctx, raidID)
	if err != nil {
		return err
	}
	return nil
}
