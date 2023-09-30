package router

import (
	"net/http"
	"tiktok_api/domain"

	"github.com/go-chi/render"
)

func TiktokAPISampleCall(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Tiktok API sample call 5",
		StatusCode: 200,
	})
}

func OAuthTiktokCallback(w http.ResponseWriter, r *http.Request) {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "OAuth tiktok api success",
		StatusCode: 200,
	})
}

// - Update user's token
func UpdateToken(w http.ResponseWriter, r *http.Request) {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Update token success",
		StatusCode: 200,
	})
}
