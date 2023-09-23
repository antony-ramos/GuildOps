package entity_test

import (
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestStrike_Validate(t *testing.T) {
	t.Parallel()
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
					"Wcm6Pbj2DbMnd4J71OzqqPmntmWd5wyiUFoVtcVNthJXFO23rQIg6MrT25DI4V1LLHmZ9dcMJUbcdaGlJ60nLT" +
					"gmKnBUhYzYC0roBXeCjBCStg16teOgFS23m6j1Yrejjba9Eyro1YOi2ETX6sCesMvKfG2N0",
				Player: &entity.Player{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			strike := entity.Strike{
				Date:   test.fields.Date,
				ID:     test.fields.ID,
				Season: test.fields.Season,
				Reason: test.fields.Reason,
				Player: test.fields.Player,
			}
			if err := strike.Validate(); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
