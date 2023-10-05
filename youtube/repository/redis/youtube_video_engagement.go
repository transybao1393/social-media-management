package redis

import (
	"encoding/json"
	"fmt"
	"tiktok_api/domain"
	"time"

	"google.golang.org/api/youtube/v3"
)

const (
	HSET_KEY = "youtube"
)

func SaveYoutubeFileUploadInfo(clientKey string, ytbFileUploadInfo *domain.YoutubeFileUploadInfo) (bool, error) {
	byte, err := json.Marshal(&ytbFileUploadInfo)
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Error when json.Marshal into redis at key %s", clientKey), "error")
		return false, err
	}

	err = clientInstance.HSet(ctx, HSET_KEY, clientKey, string(byte)).Err()
	if err != nil {
		handleError(err, "Error when save youtube file upload info into redis", "error")
		return false, err
	}
	return true, nil
}

func GetYoutubeFileUploadInfo(clientKey string, fileName string) (*domain.YoutubeFileUploadInfo, error) {
	val, err := clientInstance.HGet(ctx, HSET_KEY, clientKey).Result()
	if err != nil {
		handleError(err, "Error when get youtube file upload info from redis", "error")
		return nil, err
	}

	youtubeFileUploadInfo := &domain.YoutubeFileUploadInfo{}
	err = json.Unmarshal([]byte(val), &youtubeFileUploadInfo)
	if err != nil {
		handleError(err, "Error when unmarshal youtube file upload info from redis", "error")
		return nil, err
	}
	return youtubeFileUploadInfo, nil
}

func SaveVideoEngagementInfo(clientKey string, videoId string, videoEngagement *youtube.VideoStatistics) (bool, error) {
	videoClientKey := fmt.Sprintf("%s_%s", clientKey, videoId)
	byte, err := json.Marshal(&videoEngagement)
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Error when json.Marshal into redis at key %s", clientKey), "error")
		return false, err
	}

	expirationHour := 24 * time.Hour
	err = clientInstance.Set(ctx, videoClientKey, string(byte), expirationHour).Err()
	if err != nil {
		handleError(err, "Error when set TTL info into redis", "error")
		return false, err
	}
	return true, nil
}

func GetVideoEngagementInfo(clientKey string, videoId string) (*youtube.VideoStatistics, error) {
	videoClientKey := fmt.Sprintf("%s_%s", clientKey, videoId)
	val, err := clientInstance.Get(ctx, videoClientKey).Result()
	if err != nil {
		handleError(err, "Error when get video engagement from redis", "error")
		return nil, err
	}

	videoEngagementInfo := &youtube.VideoStatistics{}
	err = json.Unmarshal([]byte(val), &videoEngagementInfo)
	if err != nil {
		handleError(err, "Error when unmarshal video engagement info from redis", "error")
		return nil, err
	}
	return videoEngagementInfo, nil
}

func IsYoutubeVideoEngagementExist(clientKey string, videoId string) (bool, bool, error) {
	isExpireSoon := false
	//- also check expiration?
	//- redis hget check expiration
	videoClientKey := fmt.Sprintf("%s_%s", clientKey, videoId)

	remainingTime, err := clientInstance.TTL(ctx, videoClientKey).Result()
	if err != nil {
		handleError(err, "Error when check TTL of a key from redis", "error")
		return false, isExpireSoon, err
	}

	isExist, err := clientInstance.Exists(ctx, videoClientKey).Result()
	if err != nil {
		handleError(err, "Error when get video engagement from redis", "error")
		return false, isExpireSoon, err
	}
	fmt.Printf("isExist %v\n", isExist)
	//- if not exist
	if isExist == 0 {
		return false, isExpireSoon, nil
	}

	//- exist but expired soon
	//- FIXME: Improve this logic
	if remainingTime == 30*time.Second {
		fmt.Printf("Remaining time %v\n", remainingTime)
		isExpireSoon = true
	}

	//- if exist
	return true, isExpireSoon, nil
}
