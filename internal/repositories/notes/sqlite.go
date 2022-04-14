package notes

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"note-logger/internal/clock"
	"note-logger/internal/entities"

	_ "github.com/mattn/go-sqlite3"
)

const insertNoteQuery string = `
INSERT INTO notes (content, created_at) VALUES(?,?);
`

const listBetweenQuery string = `
SELECT id, content, created_at FROM notes WHERE created_at >= ? AND created_at <= ? ORDER BY created_at ASC
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
	DB *sql.DB
}

func NewRepository(cfg *Config) (Repository, error) {
	if cfg.DB == nil {
		return nil, errors.New("missing DB parameter")
	}

	newRepo := &sqliteRepo{
		dbConn: cfg.DB,
		clock:  clock.NewClock(),
	}

	return newRepo, nil
}

func (repo *sqliteRepo) Create(ctx context.Context, note *entities.Note) (*entities.Note, error) {
	note.CreatedAt = repo.clock.Now()

	res, err := repo.dbConn.ExecContext(ctx, insertNoteQuery, note.Content, note.CreatedAt)
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

	rows, err := repo.dbConn.QueryContext(ctx, listBetweenQuery, startTime, endTime)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var content string
		var createdAt time.Time

		err = rows.Scan(&id, &content, &createdAt)
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
	_, err := repo.dbConn.ExecContext(ctx, deleteNoteQuery, noteID)
	if err != nil {
		return err
	}

	return nil
}
