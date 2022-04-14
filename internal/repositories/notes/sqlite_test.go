package notes

import (
	"context"
	"regexp"
	"testing"
	"time"

	mock_clock "note-logger/internal/clock/mock"
	"note-logger/internal/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	ctx         context.Context
	mockDB      sqlmock.Sqlmock
	mockClock   *mock_clock.MockClock
	repoFixture *sqliteRepo
}

func (s *testSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()

	db, mockDB, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.mockDB = mockDB
	s.mockClock = mock_clock.NewMockClock(s.ctrl)

	s.repoFixture = &sqliteRepo{
		dbConn: db,
		clock:  s.mockClock,
	}
}

func (s *testSuite) AfterTest(_, _ string) {
	s.ctrl.Finish()

	err := s.mockDB.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

func (s *testSuite) TestNotesRepo_Create_Success() {
	createdAt := time.Unix(1649707678, 0).UTC()

	newNote := &entities.Note{
		Content: "This is a new note!",
	}

	expectedNote := &entities.Note{
		ID:        5,
		Content:   "This is a new note!",
		CreatedAt: createdAt,
	}

	s.mockClock.EXPECT().Now().Return(createdAt)

	s.mockDB.ExpectExec(regexp.QuoteMeta(insertNoteQuery)).
		WithArgs(expectedNote.Content, createdAt).WillReturnResult(sqlmock.NewResult(5, 1))

	res, err := s.repoFixture.Create(s.ctx, newNote)

	assert.Equal(s.T(), expectedNote, res)
	assert.NoError(s.T(), err)
}

func (s *testSuite) TestNotesRepo_ListBetween_Success() {
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

	rows := sqlmock.NewRows([]string{"id", "content", "created_at"}).
		AddRow(expectedNotes[0].ID, expectedNotes[0].Content, expectedNotes[0].CreatedAt).
		AddRow(expectedNotes[1].ID, expectedNotes[1].Content, expectedNotes[1].CreatedAt).
		AddRow(expectedNotes[2].ID, expectedNotes[2].Content, expectedNotes[2].CreatedAt)

	startTime := time.Unix(1649707678, 0).UTC()
	endTime := time.Unix(1649807678, 0).UTC()

	s.mockDB.ExpectQuery(regexp.QuoteMeta(listBetweenQuery)).WithArgs(startTime, endTime).WillReturnRows(rows)

	res, err := s.repoFixture.ListBetween(s.ctx, startTime, endTime)

	assert.Equal(s.T(), expectedNotes, res)
	assert.NoError(s.T(), err)
}

func (s *testSuite) TestNotesRepo_Delete_Success() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(100)

	s.mockDB.ExpectQuery(regexp.QuoteMeta(noteExistsQuery)).WithArgs(int64(100)).WillReturnRows(rows)

	s.mockDB.ExpectExec(regexp.QuoteMeta(deleteNoteQuery)).WithArgs(int64(100)).WillReturnResult(sqlmock.NewResult(100, 1))

	err := s.repoFixture.Delete(s.ctx, int64(100))

	assert.NoError(s.T(), err)
}

func TestSuites(t *testing.T) {
	suite.Run(t, new(testSuite))
}
