package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRaidUseCase_CreateRaid(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		raidUseCase := usecase.NewRaidUseCase(mockBackend)

		raid := entity.Raid{
			Name:       "raid name",
			Difficulty: "normal",
			Date:       time.Now(),
		}

		mockBackend.On("CreateRaid", mock.Anything, mock.Anything).
			Return(raid, nil)

		r, err := raidUseCase.CreateRaid(context.Background(), raid.Name, raid.Difficulty, raid.Date)

		assert.NoError(t, err)
		assert.Equal(t, raid, r)
		mockBackend.AssertExpectations(t)
	})
}
