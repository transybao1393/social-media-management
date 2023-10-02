package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tiktok_api/app/pkg/httpErrors"
	"tiktok_api/domain"

	redisRepository "tiktok_api/tiktok/repository/redis"

	usecase "tiktok_api/tiktok/usecase"

	"github.com/go-chi/render"
)

func TiktokAPISampleCall(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case or repository
	//- call to another service
	//- ex: https://jsonplaceholder.typicode.com/
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Tiktok API sample call 5",
		StatusCode: 200,
	})

	return nil
}

func RedisTest(w http.ResponseWriter, r *http.Request) error {
	val := ""
	if redisRepository.AddSimpleUser("username", "Johnathan2") { //- if add success
		val = redisRepository.GetSimpleUser("username")
	}

	render.JSON(w, r, domain.Response{
		Message:    http.StatusText(http.StatusOK),
		Data:       val,
		StatusCode: 200,
	})
	return nil
}

func OAuthTiktokCallback(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "OAuth Tiktok callback success",
		StatusCode: 200,
	})
}

func MediaUpdate(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Media update success",
		StatusCode: 200,
	})

	return nil
}

func DataRetrieval(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Data retrieval success",
		StatusCode: 200,
	})

	return nil
}

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
	var config domain.OAuth
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}
	cor, err := usecase.ListHubspotObjectFieldsUseCase(&config)
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
	var config domain.OAuth
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		return httpErrors.NewBadRequestError(err.Error())
	}

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
