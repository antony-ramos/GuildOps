package entity_test

import (
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestRaid_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name       string
		Date       time.Time
		Difficulty string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Raid",
			fields: fields{
				Name:       "raid name",
				Date:       time.Now(),
				Difficulty: "normal",
			},
			wantErr: false,
		},
		{
			name: "Invalid Raid - Name",
			fields: fields{
				Name:       "",
				Date:       time.Now(),
				Difficulty: "normal",
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Name Length",
			fields: fields{
				Name:       "raidnameraidname",
				Date:       time.Now(),
				Difficulty: "normal",
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Name Characters",
			fields: fields{
				Name:       "raidname123",
				Date:       time.Now(),
				Difficulty: "normal",
			},
			wantErr: true,
		},
		{
			name: "Invalid Raid - Difficulty",
			fields: fields{
				Name:       "raidname",
				Date:       time.Now(),
				Difficulty: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if _, err := entity.NewRaid(test.fields.Name,
				test.fields.Difficulty, test.fields.Date); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
