package redis

import (
	"os"

	"github.com/go-redis/redis"
)

var REDIS_ENDPOINT, _ = os.LookupEnv("REDIS_ENDPOINT")

type RedisService struct {
	Client *redis.Client
}

func New() RedisService {
	client := redis.NewClient(&redis.Options{
		Addr: REDIS_ENDPOINT,
		DB:   0, // use default DB
	})

	c := RedisService{client}
	return c
}
