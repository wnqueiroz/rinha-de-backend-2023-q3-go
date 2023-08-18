package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/in"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/out"
)

type personService struct {
	personChan      chan domain.Person
	persistencePort out.PersistencePort
	cachePort       out.CachePort
}

type PersonService interface {
	in.CreatePersonUseCase
	in.CountPersonUseCase
	in.SearchPersonsByTermUseCase
	in.GetPersonByIdUseCase
}

var once sync.Once
var psvc *personService

func NewPersonService(personChan chan domain.Person, persistencePort out.PersistencePort, cachePort out.CachePort) PersonService {
	once.Do(func() {
		psvc = &personService{
			persistencePort: persistencePort,
			cachePort:       cachePort,
			personChan:      personChan,
		}
	})
	return psvc
}

func (s *personService) Count(ctx context.Context) int64 {
	return s.persistencePort.CountPerson(ctx)
}

func (s *personService) Create(ctx context.Context, name, nickname string, birthday *time.Time, stack []string) (person domain.Person, err error) {
	newPerson, err := domain.NewPerson(name, nickname, birthday, stack)

	if err != nil {
		return domain.Person{}, err
	}
	if s.cachePort.CheckIfNicknameExists(ctx, nickname) {
		return domain.Person{}, errors.New("nickname already exists, use another")
	}

	s.persistencePort.StorePerson(ctx, *newPerson)
	s.cachePort.StoreNickname(ctx, newPerson.Nickname, newPerson.ID)
	s.cachePort.StorePerson(ctx, *newPerson)

	return *newPerson, nil
}

func (s *personService) GetByID(ctx context.Context, id string) (person domain.Person, err error) {
	cachedPerson, err := s.cachePort.GetPersonByID(ctx, id)
	if err == nil {
		return cachedPerson, nil
	}

	person, err = s.persistencePort.GetPersonByID(ctx, id)
	if err != nil {
		return domain.Person{}, err
	}

	return person, nil
}

func (s *personService) Search(ctx context.Context, term string) (persons []domain.Person, err error) {
	persons, err = s.persistencePort.SearchPersonsByTerm(ctx, strings.ToLower(term))
	return
}
