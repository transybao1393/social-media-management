package router

import (
	"net/http"
	"tiktok_api/domain"
	redisRepository "tiktok_api/tiktok/repository/redis"

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
	if redisRepository.AddSimpleUser("username", "johnathan") { //- if add success
		val = redisRepository.GetSimpleUser("username")
	}

	render.JSON(w, r, domain.Response{
		Message:    http.StatusText(http.StatusOK),
		Data:       val,
		StatusCode: 200,
	})
	return nil
}

func OAuthTiktokCallback(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "OAuth tiktok api success",
		StatusCode: 200,
	})

	return nil
}

// - Update user's token
func UpdateToken(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Update token success",
		StatusCode: 200,
	})

	return nil
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
