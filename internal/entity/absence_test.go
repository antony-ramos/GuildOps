package entity_test

import (
	"reflect"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestNewAbsence(t *testing.T) {
	t.Parallel()

	type args struct {
		id     int
		player *entity.Player
		raid   *entity.Raid
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Absence
		wantErr bool
	}{
		{
			name: "Valid Absence",
			args: args{
				id:     1,
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want: entity.Absence{
				ID:     1,
				Player: &entity.Player{},
				Raid:   &entity.Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Absence - Player Nil",
			args: args{
				id:     1,
				player: nil,
				raid:   &entity.Raid{},
			},
			want:    entity.Absence{},
			wantErr: true,
		},
		{
			name: "Invalid Absence - Raid Nil",
			args: args{
				id:     1,
				player: &entity.Player{},
				raid:   nil,
			},
			want:    entity.Absence{},
			wantErr: true,
		},
		{
			name: "Invalid Absence - Player and Raid Nil",
			args: args{
				id:     1,
				player: nil,
				raid:   nil,
			},
			want:    entity.Absence{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := entity.NewAbsence(test.args.id, test.args.player, test.args.raid)
			if (err != nil) != test.wantErr {
				t.Errorf("NewAbsence() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewAbsence() got = %v, want %v", got, test.want)
			}
		})
	}
}
