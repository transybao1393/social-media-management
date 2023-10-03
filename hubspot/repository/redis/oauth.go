package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tiktok_api/app/logger"
	"tiktok_api/domain"
	"tiktok_api/domain/dbInstance"

	"github.com/redis/go-redis/v9"
)

var clientInstance = dbInstance.GetRedisInstance()
var log = logger.NewLogrusLogger()
var ctx = context.Background()

func GetOneByTenantIdApiKeyType(tenantId string, apiKey string) (*domain.OAuth, error) {
	key := fmt.Sprintf("%s-%s", tenantId, apiKey)
	o := &domain.OAuth{}
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//- insert data
			oc := &domain.OAuth{
				TenantId: tenantId,
				ApiKey:   apiKey,
			}
			byte, err := json.Marshal(&oc)
			if err != nil {
				//- writing logs
				log.Fields(logger.Fields{
					"key":   key,
					"value": string(byte),
					"error": err,
				}).Errorf(err, "Error when json.Marshal into redis")
				return nil, err
			}

			err = clientInstance.Set(ctx, key, string(byte), 0).Err()
			if err != nil {
				//- writing logs
				log.Fields(logger.Fields{
					"key":   key,
					"value": string(byte),
					"error": err,
				}).Errorf(err, "Error when set into redis")
				return nil, err
			}

			return oc, nil
		}
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"error": err,
		}).Errorf(err, "Error when get into redis")
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &o)
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": val,
			"error": err,
		}).Errorf(err, "Error when json.Unmarshal into redis")
		return nil, err
	}

	return o, nil
}

func UpdateTenantDataBy(tenantId string, apiKey string, o *domain.OAuth) bool {
	key := fmt.Sprintf("%s-%s", tenantId, apiKey)
	//- check and update data
	byte, err := json.Marshal(&o)
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": string(byte),
			"error": err,
		}).Errorf(err, "Error when json.Marshal into redis")
		return false
	}

	err = clientInstance.Set(ctx, key, string(byte), 0).Err()
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": string(byte),
			"error": err,
		}).Errorf(err, "Error when Set into redis")
		return false
	}

	return true
}

func GetOneById(key string) (*domain.OAuth, error) {
	o := &domain.OAuth{}
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"error": err,
		}).Errorf(err, "Error when get into redis")
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &o)
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": val,
			"error": err,
		}).Errorf(err, "Error when json.Unmarshal into redis")
		return nil, err
	}

	return o, nil
}

func UpdateTokensById(key string, o *domain.OAuth) bool {
	//- check and update data
	byte, err := json.Marshal(&o)
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": string(byte),
			"error": err,
		}).Errorf(err, "Error when json.Marshal into redis")
		return false
	}

	err = clientInstance.Set(ctx, key, string(byte), 0).Err()
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": string(byte),
			"error": err,
		}).Errorf(err, "Error when set into redis")
		return false
	}

	return true
}
