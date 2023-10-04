package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tiktok_api/domain"
	youtubeUsecase "tiktok_api/youtube/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	MB = 1 << 20 //- 1MB
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

func YoutubeVideoEngagement(w http.ResponseWriter, r *http.Request) error {
	clientKey := chi.URLParam(r, "clientKey")
	videoId := chi.URLParam(r, "videoId")
	videoStats, _ := youtubeUsecase.YoutubeVideoEngagement(clientKey, videoId)
	//- Call logic from use case or repository
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       videoStats,
		StatusCode: 200,
	})

	return nil
}

// - using form
func YoutubeVideoUploadFile(w http.ResponseWriter, r *http.Request) error {
	//- limit to 5mb per file
	if err := r.ParseMultipartForm(5 * MB); err != nil {
		return err
	}

	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, 5*MB) // 5 Mb

	// Get handler for filename, size and headers
	clientKey := r.FormValue("client_key")
	fmt.Printf("clientKey %s\n", clientKey)
	file, handler, err := r.FormFile("file_upload")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return err
	}
	defer file.Close()

	//- save information to redis
	ytbFileUploadInfo := &domain.YoutubeFileUploadInfo{
		FileName:        handler.Filename,
		FileSize:        handler.Size,
		FileContentType: handler.Header.Get("Content-Type"),
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	fmt.Printf("Content type: %+v\n", handler.Header.Get("Content-Type"))

	//- convert to buffer
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		fmt.Printf("Error when copy file to buffer %v\n", err)
		return err
	}

	videoId, err := youtubeUsecase.YoutubeVideoUploadFile(buf, clientKey, ytbFileUploadInfo)

	dataResponse := map[string]string{}
	statusCode := 200
	message := "Video upload success"
	if err != nil {
		dataResponse["error"] = err.Error()
		statusCode = 403
		message = "Video upload failed"
	} else {
		dataResponse["youtube_channel"] = fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
		dataResponse["precaution"] = "Google Quota for video upload is 6 videos per day for free account"
	}
	//- Call logic from use case or repository
	render.Status(r, statusCode)
	render.JSON(w, r, domain.Response{
		Message:    message,
		Data:       dataResponse,
		StatusCode: statusCode,
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

	videoId, err := youtubeUsecase.YoutubeVideoUpload(vp.ClientKey, vp.VideoPath)

	dataResponse := map[string]string{}
	statusCode := 200
	message := "Video upload success"
	if err != nil {
		dataResponse["error"] = err.Error()
		statusCode = 403
		message = "Video upload failed"
	} else {
		dataResponse["youtube_channel"] = fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoId)
		dataResponse["precaution"] = "Google Quota for video upload is 6 videos per day for free account"
	}
	//- Call logic from use case or repository
	render.Status(r, statusCode)
	render.JSON(w, r, domain.Response{
		Message:    message,
		Data:       dataResponse,
		StatusCode: statusCode,
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
