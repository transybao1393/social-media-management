package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tiktok_api/domain"
	youtubeUsecase "tiktok_api/youtube/usecase"

	"github.com/go-chi/render"
)

// - generate AuthURL from config file
func GenerateAuthURL(w http.ResponseWriter, r *http.Request) error {
	authURL, clientKey := youtubeUsecase.GetAuthURL()
	render.JSON(w, r, domain.Response{
		Message: "Success",
		Data: map[string]string{
			"client_key":      clientKey,
			"google_auth_url": authURL,
		},
		StatusCode: 200,
	})
	return nil
}

func OAuthYoutubeCallback(w http.ResponseWriter, r *http.Request) {
	//- get info from database by using state - _id
	clientKey := r.FormValue("state")

	//- Get responsed code
	code := r.FormValue("code")
	fmt.Printf("code from callback %s, clientKey from callback %s", code, clientKey)
	url := youtubeUsecase.YoutubeOAuthCodeExchange(clientKey, code)

	http.Redirect(w, r, url, http.StatusMovedPermanently)
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

func YoutubeVideoUpload(w http.ResponseWriter, r *http.Request) error {
	var vp *domain.YoutubeVideoUploadPayload
	err := json.NewDecoder(r.Body).Decode(&vp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	videoId := youtubeUsecase.YoutubeVideoUpload(vp.ClientKey, vp.VideoPath)
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message: "Video upload success",
		Data: map[string]string{
			"youtube_channel": fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId),
			"precaution":      "Google Quota for video upload is 6 videos per day for free account",
		},
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
