package notes

import (
	"context"
	"note-logger/internal/entities"
	"time"
)

//go:generate mockgen -destination=mock/mock.go -package=mock_notes -source=interface.go

type Repository interface {
	Create(ctx context.Context, note *entities.Note) (*entities.Note, error)
	ListBetween(ctx context.Context, startTime time.Time, endTime time.Time) ([]*entities.Note, error)
	Delete(ctx context.Context, noteID int64) error
}
