package postgresbackend

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"

	"github.com/antony-ramos/guildops/pkg/logger"

	"github.com/antony-ramos/guildops/pkg/postgres"
	_ "github.com/lib/pq"
)

type PG struct {
	*postgres.Postgres
}

var isNotDeleted = "DELETE 0"

// Init Database Tables.
func (pg *PG) Init(ctx context.Context, connStr string, database *sql.DB) error {
	ctx, span := otel.Tracer("Backend").Start(ctx, "Init")
	defer span.End()

	select {
	case <-ctx.Done():
		return fmt.Errorf("database - Init - ctx.Done: request took too much time to be proceed")
	default:
		var err error
		if database == nil {
			database, err = sql.Open("postgres", connStr)
			if err != nil {
				return errors.Wrap(err, "database - Init - sql.Open")
			}
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				logger.FromContext(ctx).Error(err.Error())
			}
		}(database)

		// Test the connection
		err = database.Ping()
		if err != nil {
			return fmt.Errorf("database - Init - database.Ping: %w", err)
		}

		// Create a player table if it doesn't exist
		createTableSQL := `
        CREATE TABLE IF NOT EXISTS players (
            id serial PRIMARY KEY,
            name VARCHAR(255) UNIQUE,
            discord_id VARCHAR(255) UNIQUE
        );
    `

		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		// Create a table for raid
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS raids (
			id serial PRIMARY KEY,
			name VARCHAR(255),
			date TIMESTAMP,
			difficulty VARCHAR(50),
		    CONSTRAINT unique_raid_entry UNIQUE (date, difficulty)
		);
	`

		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		// Create a table for strikes
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS strikes (
			id serial PRIMARY KEY,
			player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
			season VARCHAR(50),
			reason VARCHAR(255), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		// Create a table for loots
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS loots (
			id serial PRIMARY KEY,
			name VARCHAR(255),
			raid_id INTEGER REFERENCES raids(id) ON DELETE CASCADE,
			player_id INTEGER 
				REFERENCES players(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_loot_entry UNIQUE (name, raid_id, player_id)
		);
	`

		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		// Create a table for absences
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS absences (
			id serial PRIMARY KEY,
			player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
			raid_id INTEGER REFERENCES raids(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_absence_entry UNIQUE (player_id, raid_id)
		);
	`

		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		// Create a table for fails
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS fails (
			id serial PRIMARY KEY,
			player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
			raid_id INTEGER REFERENCES raids(id) ON DELETE CASCADE,
			reason VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
		_, err = database.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("database - Init - database.Exec: %w", err)
		}

		return nil
	}
}
