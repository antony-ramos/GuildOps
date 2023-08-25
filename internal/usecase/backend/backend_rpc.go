package backend

import (
	"context"
	"errors"
	"fmt"
	"github.com/coven-discord-bot/internal/entity"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		conn, err := grpc.DialContext(ctx, r.url)
		if err != nil {
			return err
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				fmt.Print(err)
			}
		}(conn)

		absenceClient := NewAbsenceServiceClient(conn)

		addResponse, err := absenceClient.Add(context.Background(), &AbsenceRequest{
			Pseudo: absence.Name,
			Date:   []int64{absence.Date.Unix()}, // Exemple de dates Unix timestamp
		})
		if err != nil {
			return err
		}

		if addResponse.Success {
			return nil
		} else {
			return errors.New(addResponse.Message)
		}
	}
}

func (r RPC) RemoveAbsence(ctx context.Context, absence entity.Absence) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		conn, err := grpc.DialContext(ctx, r.url)
		if err != nil {
			return err
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				fmt.Print(err)
			}
		}(conn)

		absenceClient := NewAbsenceServiceClient(conn)

		removeResponse, err := absenceClient.Remove(context.Background(), &AbsenceRequest{
			Pseudo: absence.Name,
			Date:   []int64{absence.Date.Unix()}, // Exemple de dates Unix timestamp
		})
		if err != nil {
			return err
		}

		if removeResponse.Success {
			return nil
		} else {
			return errors.New(removeResponse.Message)
		}
	}
}
