package entity

import "testing"

func TestPlayer_Validate(t *testing.T) {
	type fields struct {
		ID          int
		Name        string
		Strikes     []Strike
		Loots       []Loot
		MissedRaids []Raid
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
				Strikes:     []Strike{},
				Loots:       []Loot{},
				MissedRaids: []Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Player - Name",
			fields: fields{
				ID:          1,
				Name:        "",
				Strikes:     []Strike{},
				Loots:       []Loot{},
				MissedRaids: []Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Length",
			fields: fields{
				ID:          1,
				Name:        "playernameplayername",
				Strikes:     []Strike{},
				Loots:       []Loot{},
				MissedRaids: []Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Characters",
			fields: fields{
				ID:          1,
				Name:        "playername123",
				Strikes:     []Strike{},
				Loots:       []Loot{},
				MissedRaids: []Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Lowercase",
			fields: fields{
				ID:          1,
				Name:        "PlayerName",
				Strikes:     []Strike{},
				Loots:       []Loot{},
				MissedRaids: []Raid{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Strikes:     tt.fields.Strikes,
				Loots:       tt.fields.Loots,
				MissedRaids: tt.fields.MissedRaids,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
