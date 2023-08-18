package out

import (
	"context"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
)

type PersistencePort interface {
	StorePerson(ctx context.Context, person domain.Person) error
	InsertPersonsBatch(ctx context.Context, persons []domain.Person) error
	GetPersonByID(ctx context.Context, id string) (person domain.Person, err error)
	SearchPersonsByTerm(ctx context.Context, term string) (persons []domain.Person, err error)
	CountPerson(ctx context.Context) int64
}

type CachePort interface {
	CheckIfNicknameExists(ctx context.Context, nickname string) bool
	StoreNickname(ctx context.Context, nickname string, id string) error
	StorePerson(ctx context.Context, person domain.Person) error
	GetPersonByID(ctx context.Context, id string) (person domain.Person, err error)
}
