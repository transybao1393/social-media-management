package router

import (
	"net/http"
	"tiktok_api/domain"

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
