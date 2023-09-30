package entity_test

import (
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestFail_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		ID     int
		Reason string
		Player *entity.Player
		Raid   *entity.Raid
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Fail",
			fields: fields{
				ID:     1,
				Reason: "reason",
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
			name: "Invalid Fail Player",
			fields: fields{
				ID:     1,
				Reason: "reason",
				Player: nil,
				Raid:   &entity.Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Fail Raid",
			fields: fields{
				ID:     1,
				Reason: "reason",
				Player: &entity.Player{},
				Raid:   nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid Reason",
			fields: fields{
				ID:     1,
				Reason: "",
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
			fail := entity.Fail{
				ID:     test.fields.ID,
				Reason: test.fields.Reason,
				Player: test.fields.Player,
				Raid:   test.fields.Raid,
			}
			if err := fail.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
