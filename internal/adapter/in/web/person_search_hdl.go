package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/in"
)

type PersonSearchHandler struct {
	ctx context.Context
	uc  in.SearchPersonsByTermUseCase
}

func NewPersonSearchHandler(ctx context.Context, uc in.SearchPersonsByTermUseCase) *PersonSearchHandler {
	return &PersonSearchHandler{
		ctx: ctx,
		uc:  uc,
	}
}

func (h *PersonSearchHandler) HandleSearchPersonsByTerm(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")

	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "the query parameter \"t\" is mandatory")
		return
	}

	persons, err := h.uc.Search(h.ctx, term)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	data, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
