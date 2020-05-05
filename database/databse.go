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

func (db *Database) InsertCommand(name, response string, security bool) {
	dbBool := 0
	if security {
		dbBool = 1
	}

	stmt, err := db.conn.Prepare(`INSERT INTO command(name, response, security) VALUES (?, ?, ?)`)
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not prepare insert command: %v", err)
	}

	_, err = stmt.Exec(name, response, dbBool)
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not execute insert command: %v", err)
	}
	defer stmt.Close()
}

func (db *Database) GetCommandByName(name string) string {
	var response string

	stmt, err := db.conn.Prepare(`SELECT response FROM command WHERE name = ?`)
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not prepare select command: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		db.logger.Fatalf("[SQL ERROR] could not execute select command: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&response)
		if err != nil {
			db.logger.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		db.logger.Fatal(err)
	}

	return response
}
