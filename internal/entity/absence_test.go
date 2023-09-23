package entity_test

import (
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestAbsence_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID     int
		Player *entity.Player
		Raid   *entity.Raid
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Absence",
			fields: fields{
				ID: 1,
				Player: &entity.Player{
					ID:   1,
					Name: "playername",
				},
				Raid: &entity.Raid{
					ID:         1,
					Name:       "raidname",
					Difficulty: "normal",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Absence Player",
			fields: fields{
				ID:     1,
				Player: nil,
				Raid:   &entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Absence Raid",
			fields: fields{
				ID:     1,
				Player: &entity.Player{},
				Raid:   nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid Absence Player Validate",
			fields: fields{
				ID: 1,
				Player: &entity.Player{
					ID:   1,
					Name: "",
				},
				Raid: &entity.Raid{
					ID:         1,
					Name:       "raidname",
					Difficulty: "normal",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Absence Raid Validate",
			fields: fields{
				ID: 1,
				Player: &entity.Player{
					ID:   1,
					Name: "playername",
				},
				Raid: &entity.Raid{
					ID:         1,
					Name:       "",
					Difficulty: "normal",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			a := entity.Absence{
				ID:     test.fields.ID,
				Player: test.fields.Player,
				Raid:   test.fields.Raid,
			}
			if err := a.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
