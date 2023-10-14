package entity_test

import (
	"reflect"
	"testing"

	"github.com/antony-ramos/guildops/internal/entity"
)

func TestNewPlayer(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewPlayer(tt.args.id, tt.args.name, tt.args.discordName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlayer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
