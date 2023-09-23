package entity

import (
	"testing"
	"time"
)

func TestRaid_Validate(t *testing.T) {
	type fields struct {
		ID         int
		Name       string
		Date       time.Time
		Difficulty string
		Absences   []*Player
		Players    []*Player
		Bench      []*Player
		Loots      []*Loot
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
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
				Absences:   []*Player{},
				Players:    []*Player{},
				Bench:      []*Player{},
				Loots:      []*Loot{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Raid{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				Date:       tt.fields.Date,
				Difficulty: tt.fields.Difficulty,
				Absences:   tt.fields.Absences,
				Players:    tt.fields.Players,
				Bench:      tt.fields.Bench,
				Loots:      tt.fields.Loots,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
