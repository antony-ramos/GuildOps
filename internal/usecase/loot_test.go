package usecase_test

import (
	"context"
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
			mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, player.Name).Return([]entity.Player{player}, nil)
		}

		p, err := LootUseCase.SelectPlayerToAssign(context.Background(), playersNames, "mythic")
		assert.NoError(t, err)
		assert.Equal(t, players[1], p)
	})

	t.Run("Test that winning is random", func(t *testing.T) {
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

		playersWinners := make([]entity.Player, 0)
		for i := 0; i < 10; i++ {
			for _, player := range players {
				mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, player.Name).Return([]entity.Player{player}, nil)
			}
			p, err := LootUseCase.SelectPlayerToAssign(context.Background(), playersNames, "mythic")
			assert.NoError(t, err)
			playersWinners = append(playersWinners, p)
		}

		mockBackend.AssertExpectations(t)

		for i := 0; i < len(playersWinners); i++ {
			for j := i + 1; j < len(players); j++ {
				if playersWinners[i].Name != playersWinners[j].Name {
					// random algorithm is random
					return
				}
			}
		}

		t.Errorf("random algorithm is not random")
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

		mockBackend.On("SearchPlayer", mock.Anything, mock.Anything, playersNames[0]).Return([]entity.Player{}, nil)

		p, err := LootUseCase.SelectPlayerToAssign(context.Background(), playersNames, "mythic")
		assert.Error(t, err)
		assert.Equal(t, entity.Player{}, p)
	})
}
