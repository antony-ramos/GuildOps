package entity_test

import (
	"reflect"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestNewPlayer(t *testing.T) {
	t.Parallel()
	type args struct {
		id          int
		name        string
		discordName string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Player
		wantErr bool
	}{
		{
			name: "Valid Player",
			args: args{
				id:          1,
				name:        "playername",
				discordName: "discordname",
			},
			want: entity.Player{
				ID:          1,
				Name:        "playername",
				DiscordName: "discordname",
			},
			wantErr: false,
		},
		{
			name: "Invalid Player - Name",
			args: args{
				id:          1,
				name:        "",
				discordName: "discordname",
			},
			want:    entity.Player{},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Length",
			args: args{
				id:          1,
				name:        "playernameplayer",
				discordName: "discordname",
			},
			want:    entity.Player{},
			wantErr: true,
		},
		{
			name: "Invalid Player - Name Characters",
			args: args{
				id:          1,
				name:        "player123",
				discordName: "discordname",
			},
			want:    entity.Player{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := entity.NewPlayer(test.args.id, test.args.name, test.args.discordName)
			if (err != nil) != test.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("NewPlayer() got = %v, want %v", got, test.want)
			}
		})
	}
}
