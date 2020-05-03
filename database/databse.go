package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	logger *log.Logger
	conn   *sql.DB
}

func NewDatabase(logger *log.Logger, conn *sql.DB) *Database {
	return &Database{logger, conn}
}

func (db *Database) InsertCommand(name string) {
	stmt, err := db.conn.Prepare(`INSERT INTO command(name) VALUES (?)`) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not prepare insert command: %v", err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not execute insert command: %v", err)
	}
	defer stmt.Close()
}
