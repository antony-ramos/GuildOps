package entity_test

import (
	"testing"
	"time"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestStrike_NewStrike(t *testing.T) {
	t.Parallel()
	type fields struct {
		Reason string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid Strike",
			fields: fields{
				Reason: "reason",
			},
			wantErr: false,
		},
		{
			name: "Invalid Strike - Reason",
			fields: fields{
				Reason: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid Strike - Reason Length",
			fields: fields{
				Reason: "ZBNSZVmKQwgZCBU9KjsbEOEewrPl5U1XkH10K4uXYVTuZiZiWzcydA1ISnH7iapcneGpm4CjbdMd1FdDyxuQ4eluwy3jP7kfrLhT" +
					"Wcm6Pbj2DbMnd4J71OzqqPmntmWd5wyiUFoVtcVNthJXFO23rQIg6MrT25DI4V1LLHmZ9dcMJUbcdaGlJ60nLT" +
					"gmKnBUhYzYC0roBXeCjBCStg16teOgFS23m6j1Yrejjba9Eyro1YOi2ETX6sCesMvKfG2N0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if _, err := entity.NewStrike(test.fields.Reason); (err != nil) != test.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestSeasonCalculator(t *testing.T) {
	t.Parallel()
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Season 2",
			args: args{
				date: time.Date(2023, 5, 2, 0, 0, 0, 0, time.UTC),
			},
			want: "DF/S2",
		},
		{
			name: "Unknown Season",
			args: args{
				date: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1),
			},
			want: "Unknown",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := entity.SeasonCalculator(tt.args.date); got != tt.want {
				t.Errorf("SeasonCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}
