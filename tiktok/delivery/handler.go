package router

import (
	"context"
	"net/http"
	"tiktok_api/app/logger"
	"tiktok_api/domain"
	"tiktok_api/domain/dbInstance"

	"github.com/go-chi/render"
)

var log = logger.NewLogrusLogger()
var ctx = context.Background()

func TiktokAPISampleCall(w http.ResponseWriter, r *http.Request) error {
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Tiktok API sample call 5",
		StatusCode: 200,
	})

	return nil
}

func RedisTest(w http.ResponseWriter, r *http.Request) error {
	clientInstance := dbInstance.GetRedisInstance()
	err := clientInstance.Set(ctx, "username", "johnathan", 0).Err() //- never expire

	if err != nil {
		log.Fatal(err, "Error when set new key-value into redis") //- only print error
	}

	val, err := clientInstance.Get(ctx, "username").Result()
	if err != nil {
		panic(err)
	}

	render.JSON(w, r, domain.Response{
		Message:    http.StatusText(http.StatusOK),
		Data:       val,
		StatusCode: 200,
	})
	return nil
}

func OAuthTiktokCallback(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "OAuth tiktok api success",
		StatusCode: 200,
	})

	return nil
}

// - Update user's token
func UpdateToken(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Update token success",
		StatusCode: 200,
	})

	return nil
}

func MediaUpdate(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Media update success",
		StatusCode: 200,
	})

	return nil
}

func DataRetrieval(w http.ResponseWriter, r *http.Request) error {
	//- Call logic from use case
	render.JSON(w, r, domain.Response{
		Message:    "Success",
		Data:       "Data retrieval success",
		StatusCode: 200,
	})

	return nil
}
