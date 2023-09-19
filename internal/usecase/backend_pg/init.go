package backend_pg

import (
	"database/sql"
	"log"

	"github.com/coven-discord-bot/pkg/postgres"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PG struct {
	*postgres.Postgres
}

// Init Database Tables.
func (pg *PG) Init(connStr string) {
	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(db)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create a player table if it doesn't exist
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS players (
            id serial PRIMARY KEY,
            name VARCHAR(255) UNIQUE
        );
    `

	// Create a table for raid
	createTableSQL += `
		CREATE TABLE IF NOT EXISTS raids (
			id serial PRIMARY KEY,
			name VARCHAR(255),
			date TIMESTAMP,
			difficulty VARCHAR(50),
		    CONSTRAINT unique_raid_entry UNIQUE (name, date, difficulty)
		);
	`

	// Create a table for strikes
	createTableSQL += `
		CREATE TABLE IF NOT EXISTS strikes (
			id serial PRIMARY KEY,
			player_id INTEGER REFERENCES players(id),
			season VARCHAR(50),
			reason VARCHAR(100), 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	// Create a table for loots
	createTableSQL += `
		CREATE TABLE IF NOT EXISTS loots (
			id serial PRIMARY KEY,
			name VARCHAR(255),
			raid_id INTEGER REFERENCES raids(id),
			player_id INTEGER REFERENCES players(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_loot_entry UNIQUE (name, raid_id, player_id)
		);
	`

	// Create a table for absences
	createTableSQL += `
		CREATE TABLE IF NOT EXISTS absences (
		    			id serial PRIMARY KEY,
		    			player_id INTEGER REFERENCES players(id),
		    			raid_id INTEGER REFERENCES raids(id),
		    			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    			CONSTRAINT unique_absence_entry UNIQUE (player_id, raid_id)
		    		);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	zap.L().With(zap.String("table", "raids"))
	zap.L().Info("Player table created (if it didn't exist)")
}
