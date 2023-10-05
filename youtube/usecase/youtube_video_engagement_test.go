package usecase_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestClient_GetWithMiniredis(t *testing.T) {
	mr := miniredis.RunT(t)
	mr.Set("key-exists", "exists")

	t.Cleanup(func() {
		mr.Close()
	})

	t.Run("asdasd", func(t *testing.T) {
		// Set up the client
		rc := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer rc.Close()

		err := rc.Get(context.TODO(), "key-exists").Err()
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
