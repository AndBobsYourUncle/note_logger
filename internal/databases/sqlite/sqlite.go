package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
)

const dbFile string = "notes.sqlite"

const createTableIfNotExistsQuery string = `
CREATE TABLE IF NOT EXISTS notes (
id INTEGER NOT NULL PRIMARY KEY,
created_at DATETIME NOT NULL,
content TEXT
);`

const createIndexIfNotExistsQuery string = `
CREATE INDEX IF NOT EXISTS created_at_index 
ON notes(created_at);
`

func New() (*sql.DB, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	exPath := filepath.Dir(ex)

	filename := exPath + "/" + dbFile

	err = touchDBFile(filename)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	if _, err := db.Exec(createTableIfNotExistsQuery); err != nil {
		return err
	}

	if _, err := db.Exec(createIndexIfNotExistsQuery); err != nil {
		return err
	}

	return nil
}

func touchDBFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			return createErr
		}

		closeErr := file.Close()
		if closeErr != nil {
			return closeErr
		}
	}

	return nil
}
