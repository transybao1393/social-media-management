package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tiktok_api/app/pkg/httpErrors"
	"tiktok_api/domain"

	usecase "tiktok_api/hubspot/usecase"

	"github.com/go-chi/render"
)

// - Hubspot OAuth
func OAuthHubspotCallback(w http.ResponseWriter, r *http.Request) {
	//- get info from database by using state - _id
	id := r.FormValue("state")

	//- Get responsed code
	code := r.FormValue("code")

	url, err := usecase.OAuthHubspotCallbackUseCase(id, code)
	if err != nil {
		fmt.Printf("Bad request: CSRF violation errors %s", err.Error())
		w.WriteHeader(http.StatusForbidden)
		render.JSON(w, r, domain.Response{
			StatusCode: http.StatusForbidden,
			Message:    http.StatusText(http.StatusForbidden),
			Data:       fmt.Sprintf("Bad request: CSRF violation errors %s", err.Error()),
		})
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func ListHubspotObjectFields(w http.ResponseWriter, r *http.Request) error {
	// var config domain.OAuth
	// err := json.NewDecoder(r.Body).Decode(&config)
	// if err != nil {
	// 	return httpErrors.NewBadRequestError(err.Error())
	// }
	accessToken := r.Context().Value("access_token").(string)
	cor, err := usecase.ListHubspotObjectFieldsUseCase(accessToken)
	if err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}

	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       cor,
		StatusCode: 200,
	})

	return nil
}

func UpdateToken(w http.ResponseWriter, r *http.Request) error {
	// json new Decode body
	var config domain.OAuth
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}

	// Update hubspot Token UseCase
	finalOAuth2URL, err := usecase.UpdateTokenUseCase(&config)
	if err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, domain.Response{
		StatusCode: http.StatusCreated,
		Message:    http.StatusText(http.StatusCreated),
		Data:       finalOAuth2URL,
	})
	return nil
}
