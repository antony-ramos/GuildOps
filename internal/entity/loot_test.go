package entity_test

import (
	"github.com/antony-ramos/guildops/internal/entity"
	"testing"
)

func TestLoot_Validate(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			l := entity.Loot{
				ID:     tt.fields.ID,
				Name:   tt.fields.Name,
				Player: tt.fields.Player,
				Raid:   tt.fields.Raid,
			}
			if err := l.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
