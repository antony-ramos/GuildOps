package entity_test

import (
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestPlayer_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID          int
		Name        string
		Strikes     []entity.Strike
		Loots       []entity.Loot
		MissedRaids []entity.Raid
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Player",
			fields: fields{
				ID:          1,
				Name:        "playername",
				Strikes:     []entity.Strike{},
				Loots:       []entity.Loot{},
				MissedRaids: []entity.Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Player - Name",
			fields: fields{
				ID:          1,
				Name:        "",
				Strikes:     []entity.Strike{},
				Loots:       []entity.Loot{},
				MissedRaids: []entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Length",
			fields: fields{
				ID:          1,
				Name:        "playernameplayername",
				Strikes:     []entity.Strike{},
				Loots:       []entity.Loot{},
				MissedRaids: []entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Characters",
			fields: fields{
				ID:          1,
				Name:        "playername123",
				Strikes:     []entity.Strike{},
				Loots:       []entity.Loot{},
				MissedRaids: []entity.Raid{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			player := entity.Player{
				ID:          test.fields.ID,
				Name:        test.fields.Name,
				Strikes:     test.fields.Strikes,
				Loots:       test.fields.Loots,
				MissedRaids: test.fields.MissedRaids,
			}
			if err := player.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
