package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/out"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Person struct {
	ID         uuid.UUID  `gorm:"column:id;type:uuid"`
	Name       string     `gorm:"column:nome"`
	Nickname   string     `gorm:"column:apelido"`
	Stack      string     `gorm:"column:stack;type:text"`
	Birthday   *time.Time `gorm:"column:nascimento"`
	SearchTrgm string     `gorm:"column:search_trgm"`
}

func (Person) TableName() string {
	return "pessoas"
}

func (p *Person) toDomain() domain.Person {
	var stack []string
	json.Unmarshal([]byte(p.Stack), &stack)
	person := domain.Person{
		ID:       p.ID.String(),
		Name:     p.Name,
		Nickname: p.Nickname,
		Birthday: p.Birthday,
		Stack:    stack,
	}
	return person
}

type personPersistenceAdapter struct {
	db *gorm.DB
}

var once sync.Once
var adapter *personPersistenceAdapter

func NewPersonPersistenceAdapter(db *gorm.DB) out.PersistencePort {
	once.Do(func() {
		adapter = &personPersistenceAdapter{
			db: db,
		}
	})
	return adapter
}
func (a *personPersistenceAdapter) StorePerson(ctx context.Context, person domain.Person) error {
	stackJson, err := json.Marshal(person.Stack)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(person.ID)
	if err != nil {
		return err
	}

	searchTrgm := fmt.Sprintf("%s %s %s", strings.ToLower(person.Nickname), strings.ToLower(person.Name), strings.ToLower(strings.Join(person.Stack, " ")))

	newPerson := Person{
		ID:         id,
		Name:       person.Name,
		Nickname:   person.Nickname,
		Birthday:   person.Birthday,
		Stack:      string(stackJson),
		SearchTrgm: searchTrgm,
	}

	result := a.db.Create(newPerson)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *personPersistenceAdapter) InsertPersonsBatch(ctx context.Context, persons []domain.Person) error {
	chunks := []Person{}

	for _, person := range persons {
		stackJson, err := json.Marshal(person.Stack)
		if err != nil {
			return err
		}
		id, err := uuid.Parse(person.ID)
		if err != nil {
			return err
		}

		searchTrgm := fmt.Sprintf("%s %s %s", strings.ToLower(person.Nickname), strings.ToLower(person.Name), strings.ToLower(strings.Join(person.Stack, " ")))

		newPerson := Person{
			ID:         id,
			Name:       person.Name,
			Nickname:   person.Nickname,
			Birthday:   person.Birthday,
			Stack:      string(stackJson),
			SearchTrgm: searchTrgm,
		}
		chunks = append(chunks, newPerson)
	}

	batchSize := len(chunks)

	// TODO: remove this Clauses later
	result := a.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(chunks, batchSize)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *personPersistenceAdapter) GetPersonByID(ctx context.Context, id string) (person domain.Person, err error) {
	foundPerson := Person{}

	result := a.db.Select("id", "nome", "apelido", "nascimento", "stack").First(&foundPerson, "id = ?", id)

	if result.Error != nil {
		return domain.Person{}, result.Error
	}

	return foundPerson.toDomain(), nil
}

func (a *personPersistenceAdapter) CountPerson(ctx context.Context) int64 {
	var count int64

	a.db.Model(&Person{}).Count(&count)

	return count
}

func (a *personPersistenceAdapter) SearchPersonsByTerm(ctx context.Context, term string) (persons []domain.Person, err error) {
	foundPersons := []Person{}

	result := a.db.Model(&Person{}).Distinct("id", "nome", "apelido", "nascimento", "stack").Limit(50).Find(
		&foundPersons,
		"search_trgm LIKE ?",
		fmt.Sprintf(`%%%s%%`, term),
	)

	if result.Error != nil {
		return []domain.Person{}, result.Error
	}

	for _, foundPerson := range foundPersons {
		persons = append(persons, foundPerson.toDomain())
	}

	return persons, nil
}
