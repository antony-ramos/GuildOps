package entity_test

import (
	"github.com/antony-ramos/guildops/internal/entity"
	"testing"
	"time"
)

func TestStrike_Validate(t *testing.T) {
	type fields struct {
		Date   time.Time
		ID     int
		Season string
		Reason string
		Player *entity.Player
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Strike",
			fields: fields{
				Date:   time.Now(),
				ID:     1,
				Season: "DF/S2",
				Reason: "reason",
				Player: &entity.Player{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Strike - Reason",
			fields: fields{
				Date:   time.Now(),
				ID:     1,
				Season: "DF/S2",
				Reason: "",
				Player: &entity.Player{},
			},
			wantErr: true,
		},
		{
			name: "Invalid Strike - Reason Length",
			fields: fields{
				Date:   time.Now(),
				ID:     1,
				Season: "DF/S2",
				Reason: "ZBNSZVmKQwgZCBU9KjsbEOEewrPl5U1XkH10K4uXYVTuZiZiWzcydA1ISnH7iapcneGpm4CjbdMd1FdDyxuQ4eluwy3jP7kfrLhT" +
					"Wcm6Pbj2DbMnd4J71OzqqPmntmWd5wyiUFoVtcVNthJXFO23rQIg6MrT25DI4V1LLHmZ9dcMJUbcdaGlJ60nLTgmKnBUhYzYC0roBXeCjBCStg16teOgFS23m6j1Yrejjba9Eyro1YOi2ETX6sCesMvKfG2N0",
				Player: &entity.Player{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := entity.Strike{
				Date:   tt.fields.Date,
				ID:     tt.fields.ID,
				Season: tt.fields.Season,
				Reason: tt.fields.Reason,
				Player: tt.fields.Player,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
