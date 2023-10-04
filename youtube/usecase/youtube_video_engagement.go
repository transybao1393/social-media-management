package usecase

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"tiktok_api/app/logger"
	"tiktok_api/domain"
	"tiktok_api/youtube/repository/redis"

	"google.golang.org/api/youtube/v3"
)

func handleError(err error, message string, errorType string) {
	fields := logger.Fields{
		"service": "Youtube",
		"message": message,
	}
	switch errorType {
	case "fatal":
		log.Fields(fields).Fatalf(err, message)
	case "error":
		log.Fields(fields).Errorf(err, message)
	case "warn":
		log.Fields(fields).Warnf(message)
	case "info":
		log.Fields(fields).Infof(message)
	case "debug":
		log.Fields(fields).Debugf(message)
	}
}

func YoutubeVideoUploadFile(fileBuffer *bytes.Buffer, clientKey string, ytbFileUploadInfo *domain.YoutubeFileUploadInfo) (string, error) {
	//- default params
	//- default params
	title := "This is a test video"
	description := "This is a test video from Johnathan using Youtube API"
	category := "22"
	keywords := "video, test"
	privacy := "unlisted"

	service := BuildServiceFromToken(clientKey)

	//- video uploading using youtube service
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	//- The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}
	parts := []string{
		"snippet",
		"status",
	}
	call := service.Videos.Insert(parts, upload)

	response, err := call.Media(fileBuffer).Do()
	if err != nil {
		log.Printf("Cannot upload video with error %v", err)
		return "", err
	}

	//- when success, then save youtube file upload info into redis
	_, err = redis.SaveYoutubeFileUploadInfo(clientKey, ytbFileUploadInfo)
	if err != nil {
		log.Printf("Save youtube upload file failed %v", err)
		return "", err
	}
	log.Printf("Upload successful! Video ID: %v\n", response.Id)
	return response.Id, nil
}

func YoutubeVideoUpload(clientKey string, videoFilePath string) (string, error) {
	//- default params
	//- default params
	filename := videoFilePath //- upload file path
	title := "This is a test video"
	description := "This is a test video from Johnathan using Youtube API"
	category := "22"
	keywords := "video, test"
	privacy := "unlisted"

	service := BuildServiceFromToken(clientKey)

	//- video uploading using youtube service
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}
	parts := []string{
		"snippet",
		"status",
	}
	call := service.Videos.Insert(parts, upload)

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %v: %v", filename, err)
		return "", err
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		log.Printf("Cannot upload video with error %v", err)
		return "", err
	}
	log.Printf("Upload successful! Video ID: %v\n", response.Id)
	return response.Id, nil
}

// - get current video engagement
func YoutubeVideoEngagement(clientKey string, videoId string) (*youtube.VideoStatistics, error) {
	service := BuildServiceFromToken(clientKey)
	parts := []string{
		"id",
		"snippet",
		"statistics",
	}

	//- priority to show on redis rather than call to youtube api for 24 hours
	//- call to redis to get video engagement
	isVideoEngagementExist, isExpireSoon, _ := redis.IsYoutubeVideoEngagementExist(clientKey, videoId)
	fmt.Printf("isVideoEngagementExist %v\n", isVideoEngagementExist)
	if isVideoEngagementExist && !isExpireSoon {
		videoEngagement, _ := redis.GetVideoEngagementInfo(clientKey, videoId)
		fmt.Println("Get video engagement from redis")
		return videoEngagement, nil
	}

	//- if not exist => call to youtube api to get video engagement
	call := service.Videos.List(parts)
	call = call.Id(videoId)
	response, err := call.Do()

	if err != nil {
		handleError(err, "Error when call service.Channels.List()", "error")
		return nil, err
	}

	log.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))

	fmt.Println("Get video engagement from calling Youtube API")

	//- update or cache video engagement to redis
	_, err = redis.SaveVideoEngagementInfo(clientKey, videoId, response.Items[0].Statistics)
	if err != nil {
		handleError(err, "Error when save video engagement information", "error")
		return nil, err
	}

	return response.Items[0].Statistics, nil
}

func ChannelsListByUsername(service *youtube.Service, part string, forUsername string) {
	var parts []string
	parts = append(parts, part)
	call := service.Channels.List(parts)
	call = call.ForUsername(forUsername)
	response, err := call.Do()

	if err != nil {
		handleError(err, "Error when call service.Channels.List()", "error")
	}

	log.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

// Retrieve playlistItems in the specified playlist
func PlaylistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	var parts []string
	parts = append(parts, part)
	call := service.PlaylistItems.List(parts)

	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	handleError(err, "Cannot get playlistItemsList", "error")
	return response
}

// Retrieve resource for the authenticated user's channel
func ChannelsListMine(service *youtube.Service, part string) *youtube.ChannelListResponse {
	var parts []string
	parts = append(parts, part)
	call := service.Channels.List(parts)
	call = call.Mine(true)
	response, err := call.Do()
	handleError(err, "Cannot get channelsListMine", "error")
	return response
}
