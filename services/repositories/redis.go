package repositories

//region imports
import (
	"context"
	"time"
	"user/connection"

	redis "github.com/go-redis/redis/v8"
)

type Redis struct {
	client  redis.UniversalClient
	context context.Context
}

type IRedisRepository interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	SetNX(key string, value interface{}, ttl time.Duration) *redis.BoolCmd
	SetKeepTTL(key string, value interface{}) *redis.StatusCmd
	Delete(key string) *redis.IntCmd
	GetSetTTL(key string, ttl time.Duration) *redis.StringCmd
	PersistKey(key string) *redis.BoolCmd
}

func InitRedisRepository() *Redis {
	return &Redis{client: connection.RedisClient, context: context.Background()}
}

func (redis *Redis) Get(key string) *redis.StringCmd {
	result := redis.client.Get(redis.context, key)
	return result
}

func (redis *Redis) Set(key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	result := redis.client.Set(redis.context, key, value, ttl)
	return result
}

func (redis *Redis) Delete(key string) *redis.IntCmd {
	result := redis.client.Del(redis.context, key)
	return result
}

func (redis *Redis) GetSetTTL(key string, ttl time.Duration) *redis.StringCmd {
	result := redis.client.Get(redis.context, key)
	redis.client.Expire(redis.context, key, ttl)
	return result
}

func (redis *Redis) PersistKey(key string) *redis.BoolCmd {
	result := redis.client.Persist(redis.context, key)
	return result
}

func (redis *Redis) SetNX(key string, value interface{}, ttl time.Duration) *redis.BoolCmd {
	result := redis.client.SetNX(redis.context, key, value, ttl)
	return result
}

func (r *Redis) SetKeepTTL(key string, value interface{}) *redis.StatusCmd {
	result := r.client.SetArgs(r.context, key, value, redis.SetArgs{
		KeepTTL: true,
	})
	return result
}
