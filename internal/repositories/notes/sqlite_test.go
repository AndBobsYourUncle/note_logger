package notes

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	mock_clock "note-logger/internal/clock/mock"
	"note-logger/internal/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestMigrate(t *testing.T) {
	t.Run("successfully runs the migrations", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		db, mock := newMock()
		mockClock := mock_clock.NewMockClock(ctrl)

		repo := &sqliteRepo{
			dbConn: db,
			clock:  mockClock,
		}

		defer func() {
			repo.Close(ctx)
		}()

		mock.ExpectExec(regexp.QuoteMeta(createTableIfNotExistsQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(createIndexIfNotExistsQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Migrate(ctx)

		assert.NoError(t, err)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		ctrl.Finish()
	})
}

func TestCreate(t *testing.T) {
	t.Run("successfully creates a new note", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		db, mock := newMock()
		mockClock := mock_clock.NewMockClock(ctrl)

		repo := &sqliteRepo{
			dbConn: db,
			clock:  mockClock,
		}

		defer func() {
			repo.Close(ctx)
		}()

		createdAt := time.Unix(1649707678, 0).UTC()

		newNote := &entities.Note{
			Content: "This is a new note!",
		}

		expectedNote := &entities.Note{
			ID:        5,
			Content:   "This is a new note!",
			CreatedAt: createdAt,
		}

		mockClock.EXPECT().Now().Return(createdAt)

		mock.ExpectExec(regexp.QuoteMeta(insertNoteQuery)).
			WithArgs(createdAt, expectedNote.Content).WillReturnResult(sqlmock.NewResult(5, 1))

		res, err := repo.Create(ctx, newNote)

		assert.Equal(t, expectedNote, res)
		assert.NoError(t, err)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		ctrl.Finish()
	})
}

func TestListBetween(t *testing.T) {
	t.Run("successfully lists notes between two dates", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		db, mock := newMock()
		mockClock := mock_clock.NewMockClock(ctrl)

		repo := &sqliteRepo{
			dbConn: db,
			clock:  mockClock,
		}

		defer func() {
			repo.Close(ctx)
		}()

		expectedNotes := []*entities.Note{
			{
				ID:        1,
				Content:   "Some first note!",
				CreatedAt: time.Unix(1649707678, 0).UTC(),
			},
			{
				ID:        2,
				Content:   "Some second note!",
				CreatedAt: time.Unix(1649717678, 0).UTC(),
			},
			{
				ID:        3,
				Content:   "Some third note!",
				CreatedAt: time.Unix(1649727678, 0).UTC(),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "content"}).
			AddRow(expectedNotes[0].ID, expectedNotes[0].CreatedAt, expectedNotes[0].Content).
			AddRow(expectedNotes[1].ID, expectedNotes[1].CreatedAt, expectedNotes[1].Content).
			AddRow(expectedNotes[2].ID, expectedNotes[2].CreatedAt, expectedNotes[2].Content)

		startTime := time.Unix(1649707678, 0).UTC()
		endTime := time.Unix(1649807678, 0).UTC()

		mock.ExpectQuery(regexp.QuoteMeta(listBetweenQuery)).WithArgs(startTime, endTime).WillReturnRows(rows)

		res, err := repo.ListBetween(ctx, startTime, endTime)

		assert.Equal(t, expectedNotes, res)
		assert.NoError(t, err)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		ctrl.Finish()
	})
}

func TestDelete(t *testing.T) {
	t.Run("successfully deletes a note", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		db, mock := newMock()
		mockClock := mock_clock.NewMockClock(ctrl)

		repo := &sqliteRepo{
			dbConn: db,
			clock:  mockClock,
		}

		defer func() {
			repo.Close(ctx)
		}()

		mock.ExpectExec(regexp.QuoteMeta(deleteNoteQuery)).WithArgs(int64(100)).WillReturnResult(sqlmock.NewResult(100, 1))

		err := repo.Delete(ctx, int64(100))

		assert.NoError(t, err)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		ctrl.Finish()
	})
}
