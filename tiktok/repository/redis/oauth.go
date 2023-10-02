package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"tiktok_api/domain"

	"github.com/redis/go-redis/v9"
)

func GetOneByTenantIdApiKeyType(tenantId string, apiKey string) (*domain.OAuth, error) {
	key := fmt.Sprintf("%s-%s", tenantId, apiKey)
	o := &domain.OAuth{}
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//- insert data
			byte, err := json.Marshal(&o)
			if err != nil {
				return nil, err
			}

			err = clientInstance.Set(ctx, key, string(byte), 0).Err()
			if err != nil {
				return nil, err
			}

			return o, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func UpdateTenantDataBy(tenantId string, apiKey string, o *domain.OAuth) bool {
	key := fmt.Sprintf("%s-%s", tenantId, apiKey)
	//- check and update data
	byte, err := json.Marshal(&o)
	if err != nil {
		return false
	}

	err = clientInstance.Set(ctx, key, string(byte), 0).Err()
	if err != nil {
		return false
	}

	return true
}

func GetOneById(key string) (*domain.OAuth, error) {
	o := &domain.OAuth{}
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func UpdateTokensById(key string, o *domain.OAuth) bool {
	//- check and update data
	byte, err := json.Marshal(&o)
	if err != nil {
		return false
	}

	err = clientInstance.Set(ctx, key, string(byte), 0).Err()
	if err != nil {
		return false
	}

	return true
}
