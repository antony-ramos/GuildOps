package discordhandler_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	discordHandler "github.com/antony-ramos/guildops/internal/controller/discord"
)

func TestHumanReadableError(t *testing.T) {
	t.Parallel()
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Success",
			args{
				err: fmt.Errorf("discord - parseDate - time.Parse: parsing time \"01/01/21\" month out of range"),
			},
			"parsing time \"01/01/21\" month out of range",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			args := tt.args
			t.Parallel()
			if got := discordHandler.HumanReadableError(args.err); got != want {
				t.Errorf("HumanReadableError() = %v, want %v", got, want)
			}
		})
	}
}

func TestGenerateDateList(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		startDate := "01/01/21"
		endDate := "03/01/21"

		expectedDates := []time.Time{
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.January, 3, 0, 0, 0, 0, time.UTC),
		}

		dateList, err := discordHandler.ParseDate(startDate, endDate)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if !reflect.DeepEqual(dateList, expectedDates) {
			t.Errorf("Expected %v, but got %v", expectedDates, dateList)
		}
	})

	t.Run("Success with no endDate", func(t *testing.T) {
		t.Parallel()

		startDate := "01/01/22"
		endDate := ""

		expectedDates := []time.Time{
			time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		}

		dateList, err := discordHandler.ParseDate(startDate, endDate)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if !reflect.DeepEqual(dateList, expectedDates) {
			t.Errorf("Expected %v, but got %v", expectedDates, dateList)
		}
	})

	t.Run("Error with format", func(t *testing.T) {
		t.Parallel()

		startDate := "01-01/22"
		endDate := "01-03/21"

		_, err := discordHandler.ParseDate(startDate, endDate)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Error with format", func(t *testing.T) {
		t.Parallel()

		startDate := "01/01/22"
		endDate := "01-03/22"

		_, err := discordHandler.ParseDate(startDate, endDate)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Error with date order", func(t *testing.T) {
		t.Parallel()
		startDate := "01/03/23"
		endDate := "01/01/23"

		_, err := discordHandler.ParseDate(startDate, endDate)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})
}
