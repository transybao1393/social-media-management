package connector

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	delivery "tiktok_api/tiktok/delivery"
)

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
	r.Get("/api/call", delivery.TiktokAPISampleCall)

	// r.Post("/hubspot/update", s.UpdateToken)
	r.HandleFunc("/auth/callback", delivery.OAuthTiktokCallback)
}
