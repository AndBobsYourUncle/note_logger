package notes

import (
	"context"
	"note-logger/internal/entities"
)

//go:generate mockgen -destination=mock/mock.go -package=mock_notes -source=interface.go

type Repository interface {
	Create(ctx context.Context, note *entities.Note) (*entities.Note, error)
}
