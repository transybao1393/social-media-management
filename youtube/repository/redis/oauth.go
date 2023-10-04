package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"tiktok_api/app/logger"
	"tiktok_api/domain"
	"tiktok_api/domain/dbInstance"
	"time"
)

// - youtube require to collect client_id, client_secret, project_id from client
var clientInstance = dbInstance.GetRedisInstance()
var log = logger.NewLogrusLogger()
var ctx = context.Background()

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	b := make([]byte, length)
	randomeKey := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[randomeKey.Intn(len(charset))]
	}
	return string(b)
}

func CreateNewYoutubeClient(clientId string, clientSecret string) (string, error) {
	key := generateRandomString(12) //- same with length of objectId, 12 bytes
	yOAuthInput := &domain.YoutubeOAuthConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	yOAuthInputByte, err := json.Marshal(&yOAuthInput)
	if err != nil {
		handleError(err, fmt.Sprintf("Error when json.Marshal into redis at key %s", key), "error")
		return "", err
	}

	err = clientInstance.Set(ctx, key, string(yOAuthInputByte), 0).Err()
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Error when set into redis at key %s", key), "error")
		return "", err
	}
	return key, nil
}

func IsExist(clientKey string) bool {
	err := clientInstance.Exists(ctx, clientKey).Err()
	return err == nil
}

func GetClientByClientKey(clientKey string) *domain.YoutubeOAuth {
	key := clientKey
	val, err := clientInstance.Get(ctx, key).Result() //- expect this will be single value
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("key %s not exist", key), "error")
		return nil
	}

	yOAuth := &domain.YoutubeOAuth{}
	err = json.Unmarshal([]byte(val), &yOAuth)
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Value of key %s failed to parse to json", key), "error")
		return nil
	}
	return yOAuth
}

func UpdateYoutubeByClientKey(clientKey string, youtubeOAuth *domain.YoutubeOAuth) bool {
	//- check and update data
	byte, err := json.Marshal(&youtubeOAuth)
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Error when json.Marshal into redis at key %s", clientKey), "error")
		return false
	}

	err = clientInstance.Set(ctx, clientKey, string(byte), 0).Err()
	if err != nil {
		//- writing logs and error handling
		handleError(err, fmt.Sprintf("Error when update value at key %s", clientKey), "error")
		return false
	}

	return true
}

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
