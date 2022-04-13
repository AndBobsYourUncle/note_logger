package sqlite

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSQLite_Migrate(t *testing.T) {
	ctx := context.Background()

	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)

	mockDB.ExpectExec(regexp.QuoteMeta(createTableIfNotExistsQuery)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectExec(regexp.QuoteMeta(createIndexIfNotExistsQuery)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = migrate(ctx, db)
	assert.NoError(t, err)

	err = mockDB.ExpectationsWereMet()
	assert.NoError(t, err)
}
