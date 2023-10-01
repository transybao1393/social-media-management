package connector

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"tiktok_api/app/pkg/httpErrors"

	delivery "tiktok_api/tiktok/delivery"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := h(w, r)
	if err != nil {
		switch e := err.(type) {
		case httpErrors.Error:
			w.WriteHeader(e.Status())
			render.JSON(w, r, httpErrors.NewRestError(e.Status(), e.Error(), e.Causes()))
		default:
			render.JSON(w, r, httpErrors.NewInternalServerError(e.Error()))
		}
	}
}

// - TODO: Need to add CORS and appropriate rate limiting
func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	//- General middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routing
	r.Route("/tiktok", tiktokHandler)

	return r
}

// - FIXME: seperate handler also into a seperate package
func tiktokHandler(r chi.Router) {
	//- Not applying token validation middleware
	r.Method("GET", "/api/call", Handler(delivery.TiktokAPISampleCall))

	// r.Post("/hubspot/update", s.UpdateToken)
	// r.HandleFunc("/auth/callback", Handler(delivery.OAuthTiktokCallback))
	r.Method("POST", "/auth/callback", Handler(delivery.OAuthTiktokCallback))
	r.Method("GET", "/redis/test", Handler(delivery.RedisTest))
}
