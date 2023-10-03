package connector

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"tiktok_api/app/pkg/httpErrors"

	hubspotDelivery "tiktok_api/hubspot/delivery"
	hubspotMiddleware "tiktok_api/hubspot/delivery/http/middleware"

	youtubeDelivery "tiktok_api/youtube/delivery"
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
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, httpErrors.NewRestError(http.StatusInternalServerError, e.Error(), httpErrors.ErrInternalServerError))
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
	// r.Route("/tiktok", tiktokHandler)
	r.Route("/hubspot", hubspotHandler)
	r.Route("/youtube", youtubeHandler)

	return r
}

func youtubeHandler(r chi.Router) {
	r.Method("GET", "/oauth", Handler(youtubeDelivery.GenerateAuthURL))
	r.HandleFunc("/auth/callback", youtubeDelivery.OAuthYoutubeCallback)

	r.Group(func(r chi.Router) {
		// r.Use(youtubeMiddleware.IsTokensValid)
		r.Method("POST", "/video", Handler(youtubeDelivery.YoutubeVideoUpload))
	})
}

func hubspotHandler(r chi.Router) {

	r.HandleFunc("/auth/callback", hubspotDelivery.OAuthHubspotCallback)
	r.Method("POST", "/update", Handler(hubspotDelivery.UpdateToken))

	r.Group(func(r chi.Router) {
		r.Use(hubspotMiddleware.IsTokensValid)
		r.Method("POST", "/call", Handler(hubspotDelivery.ListHubspotObjectFields))
	})
}
