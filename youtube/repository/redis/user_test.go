package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// - GetterSetter is an interface for a cache client that can get and set values. You can
// - extend this interface to add more methods, or create other interfaces that you can compose, rather
// - than using a single interface. From a testing perspective it's usually better to have multiple smaller interfaces,
// - rather than a single large one.
type GetterSetter interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

// - Client is used to implement the GetterSetter interface
type Client struct {
	//- Using an interface rather than a concrete type allows us to use a mock in our tests.
	//- In this case a concrete type would be *redis.Client.
	RedisClient redis.Cmdable
}

// - Get returns the value for the given key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.RedisClient.Get(ctx, key).Result()
}

// - Set sets the value for the given key.
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.RedisClient.Set(ctx, key, value, expiration).Err()
}

// - NewClient returns a GetterSetter that wraps a cache client.
func NewClient() GetterSetter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "", //- Update with your address
		Password: "", //- no password set
		DB:       0,  //- use default DB
	})

	return &Client{
		RedisClient: rdb,
	}
}

// - Mocking with miniredis to test more behaviour
func TestClient_GetWithMiniredis(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "key does not exist",
			key:     "no-value",
			want:    "",
			wantErr: true,
		},
		{
			name:    "key exists",
			key:     "key-exists",
			want:    "exists",
			wantErr: false,
		},
	}

	//- set up Miniredis
	mr := miniredis.RunT(t)
	//- Set key used in test
	mr.Set("key-exists", "exists")

	//- Cleanup registers a function to be called when the test (or subtest) and all its subtests complete.
	//- Cleanup functions will be called in last added, first called order.
	t.Cleanup(func() {
		mr.Close()
	})

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel() //- Run in parallel with other parallel tests

			//- Set up the client
			rc := redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			})
			defer rc.Close()

			cacheClient := &Client{
				RedisClient: rc,
			}

			got, err := cacheClient.Get(context.Background(), testCase.key)
			assert.Equal(t, testCase.want, got)
			if testCase.wantErr {
				assert.NotNil(t, err)
			}

		})
	}
}

func TestClient_SetWithMiniredis(t *testing.T) {
	t.Run("test Set using Miniredis", func(t *testing.T) {
		const wantKey = "new-key"
		const wantValue = "value"
		//- set up Miniredis
		mr := miniredis.RunT(t)
		//- Set up the client
		rc := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer func() {
			rc.Close()
			mr.Close()
		}()

		cacheClient := &Client{
			RedisClient: rc,
		}

		gotValue, _ := mr.Get(wantKey)
		t.Log("gotValue", gotValue)
		assert.Equal(t, gotValue, "")

		keyTTL := 1 * time.Minute
		err := cacheClient.Set(context.Background(), wantKey, wantValue, keyTTL)
		assert.Nil(t, err)

		gotValue, _ = mr.Get(wantKey)
		assert.Equal(t, gotValue, wantValue)

		//- Since miniredis is intended to be used in unittests TTLs don't decrease automatically.
		//- You can use TTL() to get the TTL (as a time.Duration) of a key.
		//- It will return 0 when no TTL is set.
		//
		//- m.FastForward(d) can be used to decrement all TTLs. All TTLs which become <= 0 will be removed.
		mr.FastForward(keyTTL)

		gotValue, _ = mr.Get(wantKey)
		assert.Equal(t, gotValue, "")

	})
}

//- Testing using mocks

type MockRedis struct {
	//- redis.Cmdable is embeded in the struct, so it implements the Cmdable interface,
	//- and we only need to implement the methods we care about.
	redis.Cmdable
	returnValue   string
	returnError   error
	receivedKey   string
	receivedValue any
	receivedTTL   time.Duration
}

// - Implement the Get method defined by the Cmdable interface
func (mr *MockRedis) Get(_ context.Context, key string) *redis.StringCmd {
	//- go-redis provides NewStringResult, as well as other similar methods that can be used for tests.
	return redis.NewStringResult(mr.returnValue, mr.returnError)
}

func (mr *MockRedis) Set(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	mr.receivedKey = key
	mr.receivedValue = value
	mr.receivedTTL = expiration
	return redis.NewStatusResult(mr.returnValue, mr.returnError)
}

func TestClient_GetWithMock(t *testing.T) {
	tests := []struct {
		name            string
		mockRedisClient *MockRedis
		key             string
		want            string
		wantErr         bool
	}{
		{
			name: "key does not exist",
			mockRedisClient: &MockRedis{
				returnError: errors.New("ERR no such key"),
			},
			key:     "no-value",
			want:    "",
			wantErr: true,
		},
		{
			name: "key exists",
			mockRedisClient: &MockRedis{
				returnValue: "exists",
			},
			key:     "key-exists",
			want:    "exists",
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cacheClient := &Client{
				RedisClient: testCase.mockRedisClient,
			}

			got, err := cacheClient.Get(context.Background(), testCase.key)
			assert.Equal(t, testCase.want, got)
			if testCase.wantErr {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestClient_SetWithMock(t *testing.T) {
	tests := []struct {
		name                string
		mockRedisClient     *MockRedis
		sequenceNumberKey   string
		sequenceNumberValue string
		wantTTL             time.Duration
		wantError           bool
	}{
		{
			name: "set with error",
			mockRedisClient: &MockRedis{
				returnError: errors.New("something went wrong"),
			},
			sequenceNumberKey:   "sequenceNumberKey",
			sequenceNumberValue: "sequenceNumberValue",
			wantTTL:             1 * time.Minute,
			wantError:           true,
		},
		{
			name:                "set successfully",
			mockRedisClient:     &MockRedis{},
			sequenceNumberKey:   "sequenceNumberKey",
			sequenceNumberValue: "sequenceNumberValue",
			wantTTL:             1 * time.Minute,
			wantError:           false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cacheClient := &Client{
				RedisClient: testCase.mockRedisClient,
			}
			err := cacheClient.Set(context.Background(), testCase.sequenceNumberKey, testCase.sequenceNumberValue, 1*time.Minute)
			assert.Equal(t, testCase.sequenceNumberKey, testCase.mockRedisClient.receivedKey)
			assert.Equal(t, testCase.sequenceNumberValue, testCase.mockRedisClient.receivedValue)
			assert.Equal(t, testCase.wantTTL, testCase.mockRedisClient.receivedTTL)
			if testCase.wantError {
				assert.NotNil(t, err)
			}
		})
	}
}
