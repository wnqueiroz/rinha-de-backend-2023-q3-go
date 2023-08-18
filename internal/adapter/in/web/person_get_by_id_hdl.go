package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/in"
)

type PersonGetByIdHandler struct {
	ctx context.Context
	uc  in.GetPersonByIdUseCase
}

func NewPersonGetByIdHandler(ctx context.Context, uc in.GetPersonByIdUseCase) *PersonGetByIdHandler {
	return &PersonGetByIdHandler{
		ctx: ctx,
		uc:  uc,
	}
}

func (h *PersonGetByIdHandler) HandleGetPersonById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var id string
	if strings.Contains(path, "/pessoas/") {
		id = path[len("/pessoas/"):]
	}

	person, err := h.uc.GetByID(h.ctx, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}
