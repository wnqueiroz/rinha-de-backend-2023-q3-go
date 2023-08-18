package web

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/wnqueiroz/rinha-de-backend-2023-q3-go/internal/application/port/in"
)

type PersonCountHandler struct {
	ctx context.Context
	uc  in.CountPersonUseCase
}

func NewPersonCountHandler(ctx context.Context, uc in.CountPersonUseCase) *PersonCountHandler {
	return &PersonCountHandler{
		ctx: ctx,
		uc:  uc,
	}
}

func (h *PersonCountHandler) HandleCountPerson(w http.ResponseWriter, r *http.Request) {
	count := h.uc.Count(h.ctx)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, fmt.Sprint(count))
}
