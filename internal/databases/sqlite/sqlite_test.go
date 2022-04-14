package sqlite

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
)

func TestSQLite_Migrate(t *testing.T) {
	t.Run("from scratch", func(t *testing.T) {
		ctx := context.Background()

		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"user_version"}).AddRow(0)

		mockDB.ExpectQuery(regexp.QuoteMeta(getCurrentMigration)).WillReturnRows(rows)

		for i := 1; i <= len(migrations); i++ {
			mockDB.ExpectBegin()

			mockDB.ExpectExec(regexp.QuoteMeta(migrations[i-1].migrationQuery)).
				WillReturnResult(sqlmock.NewResult(1, 1))

			setQuery := strings.Replace(setCurrentMigration, "?", strconv.Itoa(i), 1)

			mockDB.ExpectExec(regexp.QuoteMeta(setQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

			mockDB.ExpectCommit()
		}

		err = migrate(ctx, db)
		assert.NoError(t, err)

		err = mockDB.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("has a rollback", func(t *testing.T) {
		ctx := context.Background()

		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"user_version"}).AddRow(0)

		mockDB.ExpectQuery(regexp.QuoteMeta(getCurrentMigration)).WillReturnRows(rows)

		mockDB.ExpectBegin()

		mockDB.ExpectExec(regexp.QuoteMeta(createTableIfNotExistsQuery)).WillReturnError(errors.New("some sql error"))

		mockDB.ExpectRollback()

		err = migrate(ctx, db)
		assert.Equal(t, errors.New("some sql error"), err)

		err = mockDB.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("behind by one", func(t *testing.T) {
		ctx := context.Background()

		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"user_version"}).AddRow(len(migrations) - 1)

		mockDB.ExpectQuery(regexp.QuoteMeta(getCurrentMigration)).WillReturnRows(rows)

		mockDB.ExpectBegin()

		mockDB.ExpectExec(regexp.QuoteMeta(migrations[len(migrations)-1].migrationQuery)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		setQuery := strings.Replace(setCurrentMigration, "?", strconv.Itoa(len(migrations)), 1)

		mockDB.ExpectExec(regexp.QuoteMeta(setQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

		mockDB.ExpectCommit()

		err = migrate(ctx, db)
		assert.NoError(t, err)

		err = mockDB.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("up to date", func(t *testing.T) {
		ctx := context.Background()

		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"user_version"}).AddRow(len(migrations))

		mockDB.ExpectQuery(regexp.QuoteMeta(getCurrentMigration)).WillReturnRows(rows)

		err = migrate(ctx, db)
		assert.NoError(t, err)

		err = mockDB.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
