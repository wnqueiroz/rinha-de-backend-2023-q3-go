package in

import (
	"context"
	"time"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
)

type CreatePersonUseCase interface {
	Create(ctx context.Context, name, nickname string, birthday *time.Time, stack []string) (person domain.Person, err error)
}

type GetPersonByIdUseCase interface {
	GetByID(ctx context.Context, id string) (person domain.Person, err error)
}

type SearchPersonsByTermUseCase interface {
	Search(ctx context.Context, term string) (persons []domain.Person, err error)
}

type CountPersonUseCase interface {
	Count(ctx context.Context) int64
}
