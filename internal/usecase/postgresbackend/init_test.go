package postgresbackend_test

import (
	"testing"

	"github.com/antony-ramos/guildops/internal/usecase/postgresbackend"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPG_Init(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		database, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock database: %v", err)
		}
		defer database.Close()

		pgBackend := &postgresbackend.PG{nil}

		mock.ExpectExec(".*CREATE TABLE IF NOT EXISTS players.*").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(".*CREATE TABLE IF NOT EXISTS raids.*").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(".*CREATE TABLE IF NOT EXISTS strikes.*").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(".*CREATE TABLE IF NOT EXISTS loots.*").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(".*CREATE TABLE IF NOT EXISTS absences.*").WillReturnResult(sqlmock.NewResult(0, 0))

		err = pgBackend.Init("mock_conn_string", database)
		assert.NoError(t, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("Failed ping database", func(t *testing.T) {
		t.Parallel()

		pgBackend := &postgresbackend.PG{nil}
		err := pgBackend.Init("mock_conn_string", nil)
		assert.Error(t, err)
	})
}
