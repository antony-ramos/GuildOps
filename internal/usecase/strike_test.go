package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSeasonCalculator(t *testing.T) {
	t.Parallel()
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Season 2",
			args: args{
				date: time.Date(2023, 5, 2, 0, 0, 0, 0, time.UTC),
			},
			want: "DF/S2",
		},
		{
			name: "Unknown",
			args: args{
				date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, test.want, usecase.SeasonCalculator(test.args.date), "SeasonCalculator(%v)", test.args.date)
		})
	}
}

func TestStrikeUseCase_CreateStrike(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("invalid strike", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		err := strikeUseCase.CreateStrike(context.Background(),
			"ZBNSZVmKQwgZCBU9KjsbEOEewrPl5U1XkH10K4uXYVTuZiZiWzcydA1ISnH7iapcneGp"+
				"m4CjbdMd1FdDyxuQ4eluwy3jP7kfrLhT"+
				"Wcm6Pbj2DbMnd4J71OzqqPmntmWd5wyiUFoVtcVNthJXFO23rQIg6MrT25DI4V1LLHmZ9dcMJUbcdaGlJ60nLTgmKnBUhYzYC0roBXeC"+
				"jBCStg16teOgFS23m6j1Yrejjba9Eyro1YOi2ETX6sCesMvKfG2N0", "playername")

		assert.Error(t, err, "invalid strike")
		mockBackend.AssertExpectations(t)
	})

	t.Run("bug SearchPlayer", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("bug SearchPlayer"))

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "bug SearchPlayer")
		mockBackend.AssertExpectations(t)
	})

	t.Run("player not found", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "player not found")
		mockBackend.AssertExpectations(t)
	})

	t.Run("bug Create Strike", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		player := entity.Player{
			ID:   1,
			Name: "playername",
		}
		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]entity.Player{player}, nil)

		mockBackend.On("CreateStrike", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("bug Create Strike"))

		err := strikeUseCase.CreateStrike(context.Background(), "valid reason", "playername")

		assert.Error(t, err, "bug Create Strike")
		mockBackend.AssertExpectations(t)
	})

	t.Run("Context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		strikeUseCase := usecase.NewStrikeUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := strikeUseCase.CreateStrike(ctx, "valid reason", "playername")

		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}
