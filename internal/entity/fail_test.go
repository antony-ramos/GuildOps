package entity

import "testing"

func TestFail_Validate(t *testing.T) {
	type fields struct {
		ID     int
		Reason string
		Player *Player
		Raid   *Raid
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
				Player: &Player{
					ID:   1,
					Name: "playername",
				},
				Raid: &Raid{
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
				Raid:   &Raid{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Fail Raid",
			fields: fields{
				ID:     1,
				Reason: "reason",
				Player: &Player{},
				Raid:   nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid Reason",
			fields: fields{
				ID:     1,
				Reason: "",
				Player: &Player{},
				Raid:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fail{
				ID:     tt.fields.ID,
				Reason: tt.fields.Reason,
				Player: tt.fields.Player,
				Raid:   tt.fields.Raid,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
