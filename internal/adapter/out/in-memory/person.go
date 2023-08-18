package inmemory

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/patrickmn/go-cache"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/out"
)

type personMemoryAdapter struct {
	db *cache.Cache
}

var adapter *personMemoryAdapter
var once sync.Once

func NewPersonMemoryAdapter() out.CachePort {
	once.Do(func() {
		adapter = &personMemoryAdapter{
			db: cache.New(-1, -1),
		}
	})

	return adapter
}

func (a *personMemoryAdapter) StoreNickname(ctx context.Context, nickname string, id string) error {
	key := generateNicknameKey(nickname)
	a.db.Set(key, id, -1)
	return nil
}

func (a *personMemoryAdapter) CheckIfNicknameExists(ctx context.Context, nickname string) bool {
	key := generateNicknameKey(nickname)
	_, found := a.db.Get(key)
	return found
}

func (a *personMemoryAdapter) StorePerson(ctx context.Context, person domain.Person) error {
	key := generatePersonKey(person.ID)
	a.db.Set(key, person, -1)
	return nil
}

func (a *personMemoryAdapter) GetPersonByID(ctx context.Context, id string) (person domain.Person, err error) {
	key := generatePersonKey(id)
	if x, found := a.db.Get(key); found {
		person := x.(domain.Person)
		return person, nil
	}
	return domain.Person{}, errors.New("não encontrou a pessoa na memória")
}

func generateNicknameKey(nickname string) string {
	h := md5.New()
	io.WriteString(h, nickname)
	return fmt.Sprintf("nickname:%x", h.Sum(nil))
}
func generatePersonKey(id string) string {
	return fmt.Sprintf("person:%s", id)
}
