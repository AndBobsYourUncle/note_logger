package sqlite

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dbFile string = "notes.sqlite"

const getCurrentMigration string = `PRAGMA user_version;`
const setCurrentMigration string = `PRAGMA user_version = ?;`

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

type migration struct {
	migrationName  string
	migrationQuery string
}

var migrations = []migration{
	{migrationName: "create notes table", migrationQuery: createTableIfNotExistsQuery},
	{migrationName: "add notes created_at index", migrationQuery: createIndexIfNotExistsQuery},
}

func New(ctx context.Context) (*sql.DB, error) {
	filename, err := DBFilename()
	if err != nil {
		return nil, err
	}

	err = touchDBFile(filename)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	err = migrate(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	var currentMigration int

	row := db.QueryRowContext(ctx, getCurrentMigration)

	err := row.Scan(&currentMigration)
	if err != nil {
		return err
	}

	requiredMigration := len(migrations)

	if currentMigration < requiredMigration {
		log.Printf("Current DB version: %v, required DB version: %v\n", currentMigration, requiredMigration)

		for migrationNum := currentMigration + 1; migrationNum <= requiredMigration; migrationNum++ {
			err = execMigration(ctx, db, migrationNum)
			if err != nil {
				log.Printf("Error running migration %v '%v'\n", migrationNum, migrations[migrationNum-1].migrationName)

				return err
			}
		}
	}

	return nil
}

func execMigration(ctx context.Context, db *sql.DB, migrationNum int) (err error) {
	log.Printf("Running migration %v '%v'\n", migrationNum, migrations[migrationNum-1].migrationName)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, migrations[migrationNum-1].migrationQuery)
	if err != nil {
		return err
	}

	setQuery := strings.Replace(setCurrentMigration, "?", strconv.Itoa(migrationNum), 1)

	_, err = tx.ExecContext(ctx, setQuery)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func DBFilename() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	exPath := filepath.Dir(ex)

	return exPath + "/" + dbFile, nil
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
