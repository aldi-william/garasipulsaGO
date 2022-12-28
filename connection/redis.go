package connection

import (
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var (
	RedisClient redis.UniversalClient
)

func InitRedis() error {
	var (
		redisHost     string = os.Getenv("REDIS_HOST")
		redisPassword string = os.Getenv("REDIS_PASSWORD")
	)

	RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(redisHost, ","),
		Password: redisPassword,
	})

	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		return errors.Wrap(err, "failed ping redis")
	}

	return nil
}
