package entity_test

import (
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestLoot_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID     int
		Name   string
		Player *entity.Player
		Raid   *entity.Raid
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Loot",
			fields: fields{
				ID:     1,
				Name:   "LootName",
				Player: &entity.Player{},
				Raid:   &entity.Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Loot - Name",
			fields: fields{
				ID:     1,
				Name:   "",
				Player: &entity.Player{},
				Raid:   &entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Loot - Player",
			fields: fields{
				ID:     1,
				Name:   "LootName",
				Player: nil,
				Raid:   &entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Loot - Raid",
			fields: fields{
				ID:     1,
				Name:   "LootName",
				Player: &entity.Player{},
				Raid:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			loot := entity.Loot{
				ID:     test.fields.ID,
				Name:   test.fields.Name,
				Player: test.fields.Player,
				Raid:   test.fields.Raid,
			}
			if err := loot.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
