package handler

import (
	"net/http"
	"tiktok_api/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GeneralHandler(r chi.Router) {
	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, domain.Response{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       "sample response",
		})
	})
}
