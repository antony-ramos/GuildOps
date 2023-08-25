package usecase

import (
	"context"
	"github.com/coven-discord-bot/internal/entity"
)

type Backend interface {
	AddAbsence(ctx context.Context, absence entity.Absence) error
	RemoveAbsence(ctx context.Context, absence entity.Absence) error
}
