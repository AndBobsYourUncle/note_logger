package notes

import (
	"context"
	"database/sql"
	"note-logger/internal/clock"
	"note-logger/internal/entities"
	"os"
	"path/filepath"
	"time"
)

const dbFile string = "notes.sqlite"

const createTable string = `
CREATE TABLE IF NOT EXISTS notes (
id INTEGER NOT NULL PRIMARY KEY,
created_at DATETIME NOT NULL,
content TEXT
);`

const createIndex string = `
CREATE INDEX IF NOT EXISTS created_at_index 
ON notes(created_at);
`

const insertNote string = `
INSERT INTO notes VALUES(NULL,?,?);
`

const listBetween string = `
SELECT * FROM notes WHERE created_at >= ? AND created_at <= ? ORDER BY created_at ASC
`

//go:generate mockgen -destination=mock_sql/mock.go -package=mock_sql -source=sqlite.go

type sqliteDBConn interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
}

type sqliteRepo struct {
	dbConn sqliteDBConn
	clock  clock.Clock
}

type Config struct {
	Filename string
}

func NewRepository() (Repository, error) {
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

	// ensure that the DB is set up properly with the table/indices needed
	if _, err = db.Exec(createTable); err != nil {
		return nil, err
	}

	if _, err = db.Exec(createIndex); err != nil {
		return nil, err
	}

	return &sqliteRepo{
		dbConn: db,
		clock:  clock.NewClock(),
	}, nil
}

func (repo *sqliteRepo) Create(ctx context.Context, note *entities.Note) (*entities.Note, error) {
	note.CreatedAt = repo.clock.Now()

	res, err := repo.dbConn.Exec(insertNote, time.Now(), note.Content)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	note.ID = lastID

	return note, nil
}

func (repo *sqliteRepo) ListBetween(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entities.Note, error) {
	retNotes := make([]*entities.Note, 0)

	rows, err := repo.dbConn.Query(listBetween, startTime, endTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var createdAt time.Time
		var content string

		err = rows.Scan(&id, &createdAt, &content)
		if err != nil {
			return nil, err
		}

		retNotes = append(retNotes, &entities.Note{
			ID:        id,
			Content:   content,
			CreatedAt: createdAt,
		})
	}

	return retNotes, nil
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
