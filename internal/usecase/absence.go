package usecase

import (
	"context"
	"github.com/coven-discord-bot/internal/entity"
)

type AbsenceUseCase struct {
	backend Backend
}

func NewAbsenceUseCase(bk Backend) *AbsenceUseCase {
	return &AbsenceUseCase{backend: bk}
}

func (a AbsenceUseCase) AddAbsence(ctx context.Context, absence entity.Absence) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := absence.Validate()
		if err != nil {
			return err
		}
		err = a.backend.AddAbsence(ctx, absence)
		if err != nil {
			return err
		}
		return nil
	}
}

func (a AbsenceUseCase) RemoveAbsence(ctx context.Context, absence entity.Absence) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		err := absence.Validate()
		if err != nil {
			return err
		}
		err = a.backend.RemoveAbsence(ctx, absence)
		if err != nil {
			return err
		}
		return nil
	}
}
