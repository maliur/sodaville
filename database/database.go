package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
	mu   sync.RWMutex
}

func OpenDB(name string) (*Database, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	return &Database{conn: db}, nil
}

func (db *Database) InsertCommand(name, response string, security bool) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbBool := 0
	if security {
		dbBool = 1
	}

	_, err := db.conn.Exec(`INSERT INTO command(name, response, security) VALUES (?, ?, ?)`, name, response, dbBool)
	if err != nil {
		return fmt.Errorf("[SQL ERROR] could not execute insert command: %v", err)
	}

	return nil
}

func (db *Database) DeleteCommand(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec(`DELETE FROM command WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("[SQL ERROR] could not execute delete command for name '%s' cause of '%v'", name, err)
	}

	return nil
}

func (db *Database) GetCommandByName(name string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var response string

	row := db.conn.QueryRow(`SELECT response FROM command WHERE name = ?`, name)
	if err := row.Scan(&response); err != nil {
		return "", fmt.Errorf("[SQL ERROR] could not execute select command for name '%s' cause of %v", name, err)
	}

	return response, nil
}

func (db *Database) GetAllCommands() (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	response := "available commands: "

	rows, err := db.conn.Query(`SELECT name FROM command`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return "", fmt.Errorf("[SQL ERROR] could not execute select command: %v", err)
		}

		response += fmt.Sprintf("$%s ", name)
	}

	return response, nil
}

func (db *Database) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.Close()
}
