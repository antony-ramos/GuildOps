package entity_test

import (
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestRaid_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID         int
		Name       string
		Date       time.Time
		Difficulty string
		Absences   []*entity.Player
		Players    []*entity.Player
		Bench      []*entity.Player
		Loots      []*entity.Loot
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Raid",
			fields: fields{
				ID:         1,
				Name:       "raid name",
				Date:       time.Now(),
				Difficulty: "normal",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Raid - Name",
			fields: fields{
				ID:         1,
				Name:       "",
				Date:       time.Now(),
				Difficulty: "normal",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Name Length",
			fields: fields{
				ID:         1,
				Name:       "raidnameraidname",
				Date:       time.Now(),
				Difficulty: "normal",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Name Characters",
			fields: fields{
				ID:         1,
				Name:       "raidname123",
				Date:       time.Now(),
				Difficulty: "normal",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Name Lowercase",
			fields: fields{
				ID:         1,
				Name:       "RaidName",
				Date:       time.Now(),
				Difficulty: "normal",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Difficulty",
			fields: fields{
				ID:         1,
				Name:       "raidname",
				Date:       time.Now(),
				Difficulty: "invalid",
				Absences:   []*entity.Player{},
				Players:    []*entity.Player{},
				Bench:      []*entity.Player{},
				Loots:      []*entity.Loot{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			raid := entity.Raid{
				ID:         test.fields.ID,
				Name:       test.fields.Name,
				Date:       test.fields.Date,
				Difficulty: test.fields.Difficulty,
				Absences:   test.fields.Absences,
				Players:    test.fields.Players,
				Bench:      test.fields.Bench,
				Loots:      test.fields.Loots,
			}
			if err := raid.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
