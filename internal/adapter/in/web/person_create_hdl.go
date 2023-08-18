package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/domain"
	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/in"
)

type Payload struct {
	Name     string   `json:"nome"`
	Nickname string   `json:"apelido"`
	Birthday string   `json:"nascimento"`
	Stack    []string `json:"stack"`
}

type PersonCreateHandler struct {
	ctx context.Context
	uc  in.CreatePersonUseCase
}

func NewPersonCreateHandler(ctx context.Context, uc in.CreatePersonUseCase) *PersonCreateHandler {
	return &PersonCreateHandler{
		ctx: ctx,
		uc:  uc,
	}
}

func (h *PersonCreateHandler) HandleCreatePerson(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var payload Payload
	var unmarshalErr *json.UnmarshalTypeError

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong type provided for field: "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	birthday, _ := time.Parse("2006-01-02", payload.Birthday)

	person, err := h.uc.Create(h.ctx, payload.Name, payload.Nickname, &birthday, payload.Stack)

	if err != nil {
		statusCode := http.StatusUnprocessableEntity
		switch {
		case errors.Is(err, domain.ErrNickNameTooLong):
			statusCode = http.StatusBadRequest
		case errors.Is(err, domain.ErrNameTooLong):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusUnprocessableEntity
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/pessoas/%s", person.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person.Pretty())
}
