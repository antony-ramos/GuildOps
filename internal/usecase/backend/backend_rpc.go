package backend

import (
	"context"
	"errors"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type RPC struct {
	url     string
	circuit *gobreaker.CircuitBreaker
}

func NewRPC(url string) *RPC {
	settings := gobreaker.Settings{
		Name:    "backend-rpc",
		Timeout: 2 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 2
		},
	}
	circuit := gobreaker.NewCircuitBreaker(settings)
	return &RPC{
		url:     url,
		circuit: circuit,
	}
}

func (r RPC) AddAbsence(ctx context.Context, absence entity.Absence) error {
	_, span := otel.Tracer("").Start(ctx, "Backend Connector send add absence request to backend", trace.WithTimestamp(time.Now()))
	defer span.End(trace.WithTimestamp(time.Now()))

	select {
	case <-ctx.Done():
		span.RecordError(errors.New("context time has exceeded"))
		return ctx.Err()
	default:
		conn, err := grpc.DialContext(ctx, r.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			err := fmt.Errorf("cannot dial backend %e", err)
			span.RecordError(err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				span.RecordError(err)
			}
		}(conn)

		absenceClient := NewAbsenceServiceClient(conn)

		addResponse, err := absenceClient.Add(context.Background(), &AbsenceRequest{
			Pseudo: absence.Name,
			Date:   []int64{absence.Date.Unix()}, // Exemple de dates Unix timestamp
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		if addResponse.Success {
			return nil
		} else {
			err := errors.New(addResponse.Message)
			span.RecordError(err)
			return err
		}
	}
}

func (r RPC) RemoveAbsence(ctx context.Context, absence entity.Absence) error {
	_, span := otel.Tracer("").Start(ctx, "Backend Connector send cancel absence request to backend", trace.WithTimestamp(time.Now()))
	defer span.End(trace.WithTimestamp(time.Now()))

	select {
	case <-ctx.Done():
		span.RecordError(errors.New("context time has exceeded"))
		return ctx.Err()
	default:
		conn, err := grpc.DialContext(ctx, r.url)
		if err != nil {
			err := fmt.Errorf("cannot dial backend %e", err)
			span.RecordError(err)
			return err
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				span.RecordError(err)
			}
		}(conn)

		absenceClient := NewAbsenceServiceClient(conn)

		removeResponse, err := absenceClient.Remove(context.Background(), &AbsenceRequest{
			Pseudo: absence.Name,
			Date:   []int64{absence.Date.Unix()},
		})
		if err != nil {
			span.RecordError(err)
			return err
		}

		if removeResponse.Success {
			return nil
		} else {
			err := errors.New(removeResponse.Message)
			span.RecordError(err)
			return err
		}
	}
}
