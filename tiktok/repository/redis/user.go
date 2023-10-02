package redis

import (
	"context"
	"tiktok_api/app/logger"
	"tiktok_api/domain/dbInstance"
	"time"
)

var clientInstance = dbInstance.GetRedisInstance()
var log = logger.NewLogrusLogger()
var ctx = context.Background()

func AddNew(key string, value string) bool {
	err := clientInstance.Set(ctx, key, value, 0).Err()
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": value,
			"error": err,
		}).Errorf(err, "Error when set new key-value into redis")
		return false
	}
	return true
}

func GetByKey(key string) (string, error) {
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	log.Fields(logger.Fields{
		"key":  key,
		"date": time.Now(),
	}).Info("Get by key from Redis success")
	return val, nil
}

func RemoveByKey(key string) (int64, error) {
	val, err := clientInstance.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	log.Fields(logger.Fields{
		"key":  key,
		"date": time.Now(),
	}).Info("Remove a key from Redis success")
	return val, nil
}
