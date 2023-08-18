package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID       string     `json:"id"`
	Name     string     `json:"nome"`
	Nickname string     `json:"apelido"`
	Birthday *time.Time `json:"nascimento"`
	Stack    []string   `json:"stack"`
}

type PrettyPerson struct {
	ID       string   `json:"id"`
	Name     string   `json:"nome"`
	Nickname string   `json:"apelido"`
	Birthday string   `json:"nascimento"`
	Stack    []string `json:"stack"`
}

var ErrNickNameTooLong = errors.New("nickname must be equal to or less than 32 characters")
var ErrNameTooLong = errors.New("name must be equal to or less than 100 characters")

func NewPerson(name, nickname string, birthday *time.Time, stack []string) (*Person, error) {
	if name == "" {
		return &Person{}, errors.New("name cannot be empty")
	}
	if nickname == "" {
		return &Person{}, errors.New("nickname cannot be empty")
	}
	if len(nickname) > 32 {
		return &Person{}, ErrNickNameTooLong
	}
	if len(name) > 100 {
		return &Person{}, ErrNameTooLong
	}

	return &Person{
		ID:       (uuid.New()).String(),
		Name:     name,
		Nickname: nickname,
		Birthday: birthday,
		Stack:    stack,
	}, nil
}

func (p *Person) Pretty() *PrettyPerson {
	return &PrettyPerson{
		ID:       p.ID,
		Name:     p.Name,
		Nickname: p.Nickname,
		Stack:    p.Stack,
		Birthday: p.Birthday.Format("2006-01-02"),
	}
}
