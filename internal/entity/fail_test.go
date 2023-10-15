package entity_test

import (
	"reflect"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestNewFail(t *testing.T) {
	t.Parallel()

	type args struct {
		id     int
		reason string
		player *entity.Player
		raid   *entity.Raid
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Fail
		wantErr bool
	}{
		{
			name: "Valid Fail",
			args: args{
				id:     1,
				reason: "fail reason",
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want: entity.Fail{
				ID:     1,
				Reason: "fail reason",
				Player: &entity.Player{},
				Raid:   &entity.Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Fail - Reason",
			args: args{
				id:     1,
				reason: "",
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want:    entity.Fail{},
			wantErr: true,
		},
		{
			name: "Invalid Fail - Raid Nil",
			args: args{
				id:     1,
				reason: "fail reason",
				player: &entity.Player{},
				raid:   nil,
			},
			want:    entity.Fail{},
			wantErr: true,
		},
		{
			name: "Invalid Fail - Player Nil",
			args: args{
				id:     1,
				reason: "fail reason",
				player: nil,
				raid:   &entity.Raid{},
			},

			want:    entity.Fail{},
			wantErr: true,
		},
		{
			name: "Invalid Fail - Player and Raid Nil",
			args: args{
				id:     1,
				reason: "fail reason",
				player: nil,

				raid: nil,
			},
			want:    entity.Fail{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := entity.NewFail(test.args.id, test.args.reason, test.args.player, test.args.raid)
			if (err != nil) != test.wantErr {
				t.Errorf("NewFail() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewFail() got = %v, want %v", got, test.want)
			}
		})
	}
}
