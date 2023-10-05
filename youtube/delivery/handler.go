package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tiktok_api/app/logger"
	"tiktok_api/app/utils"
	"tiktok_api/domain"
	youtubeUsecase "tiktok_api/youtube/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slices"
)

var log = logger.NewLogrusLogger()

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// - using form
func YoutubeVideoUploadFile(w http.ResponseWriter, r *http.Request) error {
	message := "Video upload success"
	//- limit to 5mb per file
	if err := r.ParseMultipartForm(5 * MB); err != nil {
		fields := logger.Fields{
			"service": "Youtube",
			"message": "Error when parse multipart form",
		}
		log.Fields(fields).Errorf(err, "Error when parse multipart form")
		return err
	}

	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, 5*MB) // 5 Mb

	// Get handler for filename, size and headers
	clientKey := r.FormValue("client_key")
	file, handler, err := r.FormFile("file_upload")
	if err != nil {
		fields := logger.Fields{
			"service": "Youtube",
			"message": fmt.Sprintf("Error when receive form file from request with client key %s", clientKey),
		}
		log.Fields(fields).Errorf(err, "Error when receive form file from request")
		return err
	}
	defer file.Close()

	// validation media type is video
	if !slices.Contains(utils.VideoContentType, handler.Header.Get("Content-Type")) {
		fields := logger.Fields{
			"service": "Youtube",
			"message": fmt.Sprintf("This file is not video content type with client key %s", clientKey),
		}
		log.Fields(fields).Errorf(err, "This file is not video content type")
		return fmt.Errorf("invalid media type, the file is not a video")
	}

	//- save information to redis
	ytbFileUploadInfo := &domain.YoutubeFileUploadInfo{
		FileName:        handler.Filename,
		FileSize:        handler.Size,
		FileContentType: handler.Header.Get("Content-Type"),
		CreatedAt:       time.Now(),
	}

	//- convert to buffer
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		fields := logger.Fields{
			"service": "Youtube",
			"message": "Error when copy file to buffer",
		}
		log.Fields(fields).Errorf(err, "Error when copy file to buffer")
		return err
	}

	videoId, err := youtubeUsecase.YoutubeVideoUploadFile(buf, clientKey, ytbFileUploadInfo)

	dataResponse := map[string]string{}
	statusCode := 200

	if err != nil {
		fields := logger.Fields{
			"service": "Youtube",
			"message": "Error when Video upload failed",
		}
		log.Fields(fields).Errorf(err, "Error when Video upload failed")

		dataResponse["error"] = err.Error()
		statusCode = 403
		message = "Video upload failed"
	} else {
		dataResponse["video_id"] = videoId
		dataResponse["client_key"] = clientKey
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
