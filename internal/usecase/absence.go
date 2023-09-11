package usecase

import (
	"context"
	"github.com/coven-discord-bot/internal/entity"
	logger "github.com/coven-discord-bot/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type AbsenceUseCase struct {
	backend Backend
}

func NewAbsenceUseCase(bk Backend) *AbsenceUseCase {
	return &AbsenceUseCase{backend: bk}
}

func (a AbsenceUseCase) AddAbsence(ctx context.Context, log logger.Logger, absence entity.Absence) error {
	_, span := otel.Tracer("").Start(ctx, "Usecase AbsenceUseCase/AddAbsence is processing", trace.WithTimestamp(time.Now()))
	defer span.End(trace.WithTimestamp(time.Now()))
	select {
	case <-ctx.Done():
		log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
		log.Debug("UseCase/AddAbsence Context Exceeded")
		return ctx.Err()
	default:
		err := absence.Validate()
		span.RecordError(err)
		if err != nil {
			log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
			log.Debug("UseCase/AddAbsence Validate Absence failed")
			return err
		}
		err = a.backend.AddAbsence(ctx, absence)
		if err != nil {
			log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
			log.Debug("UseCase/AddAbsence a.backend.AddAbsence failed")
			return err
		}
		return nil
	}
}

func (a AbsenceUseCase) RemoveAbsence(ctx context.Context, log logger.Logger, absence entity.Absence) error {
	_, span := otel.Tracer("").Start(ctx, "Usecase AbsenceUseCase/RemoveAbsence is processing", trace.WithTimestamp(time.Now()))
	defer span.End(trace.WithTimestamp(time.Now()))
	select {
	case <-ctx.Done():
		log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
		log.Debug("UseCase/AddAbsence Context Exceeded")
		return ctx.Err()
	default:
		err := absence.Validate()
		span.RecordError(err)
		if err != nil {
			log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
			log.Debug("UseCase/AddAbsence Validate Absence failed")
			return err
		}
		err = a.backend.RemoveAbsence(ctx, absence)
		if err != nil {
			log = log.With(zap.String("name", absence.Name)).With(zap.String("date", absence.Date.String()))
			log.Debug("UseCase/AddAbsence a.backend.RemoveAbsence failed")
			return err
		}
		return nil
	}
}
