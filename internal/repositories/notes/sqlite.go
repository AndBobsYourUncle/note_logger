package notes

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"note-logger/internal/clock"
	"note-logger/internal/entities"

	_ "github.com/mattn/go-sqlite3"
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

const insertNoteQuery string = `
INSERT INTO notes VALUES(NULL,?,?);
`

const listBetweenQuery string = `
SELECT * FROM notes WHERE created_at >= ? AND created_at <= ? ORDER BY created_at ASC
`

const deleteNoteQuery string = `
DELETE FROM notes WHERE id = ?
`

//go:generate mockgen -destination=mock_sql/mock.go -package=mock_sql -source=sqlite.go

type sqliteRepo struct {
	dbConn *sql.DB
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

	return &sqliteRepo{
		dbConn: db,
		clock:  clock.NewClock(),
	}, nil
}

func (repo *sqliteRepo) Migrate(ctx context.Context) error {
	if _, err := repo.dbConn.Exec(createTableIfNotExistsQuery); err != nil {
		return err
	}

	if _, err := repo.dbConn.Exec(createIndexIfNotExistsQuery); err != nil {
		return err
	}

	return nil
}

func (repo *sqliteRepo) Create(ctx context.Context, note *entities.Note) (*entities.Note, error) {
	note.CreatedAt = repo.clock.Now()

	res, err := repo.dbConn.Exec(insertNoteQuery, note.CreatedAt, note.Content)
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

	rows, err := repo.dbConn.Query(listBetweenQuery, startTime, endTime)
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

func (repo *sqliteRepo) Delete(ctx context.Context, noteID int64) error {
	_, err := repo.dbConn.Exec(deleteNoteQuery, noteID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *sqliteRepo) Close(ctx context.Context) error {
	return repo.dbConn.Close()
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
