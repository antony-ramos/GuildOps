package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
	"github.com/antony-ramos/guildops/internal/usecase"
	"github.com/antony-ramos/guildops/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLootUseCase_SelectPlayerToAssign(t *testing.T) {
	t.Parallel()

	t.Run("context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := LootUseCase.SelectPlayerToAssign(ctx, []string{"playerone", "playertwo"}, "mythic")
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		players := []entity.Player{
			{
				ID:   1,
				Name: "playerone",
				Loots: []entity.Loot{
					{
						ID:   1,
						Name: "lootone",
						Raid: &entity.Raid{
							ID:         1,
							Name:       "castle nathria",
							Difficulty: "mythic",
						},
					},
					{
						ID:   1,
						Name: "loottwo",
						Raid: &entity.Raid{
							ID:         1,
							Name:       "castle nathria",
							Difficulty: "mythic",
						},
					},
				},
			},
			{
				ID:   2,
				Name: "playertwo",
				Loots: []entity.Loot{
					{
						ID:   1,
						Name: "lootone",
						Raid: &entity.Raid{
							ID:         1,
							Name:       "castle nathria",
							Difficulty: "mythic",
						},
					},
				},
			},
		}
		playersNames := []string{"playerone", "playertwo"}

		for _, player := range players {
			mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, player.Name, mock.Anything).
				Return([]entity.Player{player}, nil)
		}

		p, err := LootUseCase.SelectPlayerToAssign(context.Background(), playersNames, "mythic")
		assert.NoError(t, err)
		assert.Equal(t, players[1], p)
	})

	t.Run("Player List empty", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		p, err := LootUseCase.SelectPlayerToAssign(context.Background(), []string{}, "mythic")
		assert.Error(t, err)
		assert.Equal(t, entity.Player{}, p)
	})

	t.Run("Players doesnt exists", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		playersNames := []string{"playerone", "playertwo"}

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, playersNames[0], mock.Anything).
			Return([]entity.Player{}, nil)

		p, err := LootUseCase.SelectPlayerToAssign(context.Background(), playersNames, "mythic")
		assert.Error(t, err)
		assert.Equal(t, entity.Player{}, p)
	})
}

func TestLootUseCase_DeleteLoot(t *testing.T) {
	t.Parallel()

	t.Run("context is done", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := LootUseCase.DeleteLoot(ctx, 1)
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		mockBackend.On("DeleteLoot", mock.Anything, mock.Anything).Return(nil)

		err := LootUseCase.DeleteLoot(context.Background(), 1)
		assert.NoError(t, err)
		mockBackend.AssertExpectations(t)
	})

	t.Run("Backend Error", func(t *testing.T) {
		t.Parallel()

		mockBackend := mocks.NewBackend(t)

		LootUseCase := usecase.NewLootUseCase(mockBackend)

		mockBackend.On("DeleteLoot", mock.Anything, mock.Anything).Return(errors.New("Backend Error"))

		err := LootUseCase.DeleteLoot(context.Background(), 1)
		assert.Error(t, err)
		mockBackend.AssertExpectations(t)
	})
}

func TestLootUseCase_CreateLoot(t *testing.T) {

}
