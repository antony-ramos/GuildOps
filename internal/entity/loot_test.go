package entity_test

import (
	"reflect"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestNewLoot(t *testing.T) {
	t.Parallel()

	type args struct {
		id     int
		name   string
		player *entity.Player
		raid   *entity.Raid
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Loot
		wantErr bool
	}{
		{
			name: "Valid Loot",
			args: args{
				id:     1,
				name:   "loot name",
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want: entity.Loot{
				ID:     1,
				Name:   "loot name",
				Player: &entity.Player{},
				Raid:   &entity.Raid{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Loot - Name",
			args: args{
				id:     1,
				name:   "",
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want:    entity.Loot{},
			wantErr: true,
		},
		{
			name: "Invalid Loot - Name Length",
			args: args{
				id:     1,
				name:   "lootnameadzqdqzdklqzdkqzdlqzkdqlzdkqzmdkqzdlqzdkqzd",
				player: &entity.Player{},
				raid:   &entity.Raid{},
			},
			want:    entity.Loot{},
			wantErr: true,
		},
		{
			name: "Invalid Loot - Player",
			args: args{
				id:     1,
				name:   "loot name",
				player: nil,
				raid:   &entity.Raid{},
			},
			want:    entity.Loot{},
			wantErr: true,
		},
		{
			name: "Invalid Loot - Raid",
			args: args{
				id:     1,
				name:   "loot name",
				player: &entity.Player{},
				raid:   nil,
			},
			want:    entity.Loot{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := entity.NewLoot(test.args.id, test.args.name, test.args.player, test.args.raid)
			if (err != nil) != test.wantErr {
				t.Errorf("NewLoot() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewLoot() got = %v, want %v", got, test.want)
			}
		})
	}
}
