package throttler

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	redisService "gitlab.com/ncent/throttler/services/redis"
)

type Throttler struct {
	RedisService redisService.RedisService
}

func New() Throttler {

	t := Throttler{redisService.New()}
	return t
}

func (t Throttler) LimitToNTimesByNHours(key string, times int, hours int) error {
	val, err := t.RedisService.Client.Get(key).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get key %+v", err)
		return err
	}
	log.Printf("Val: %+v", val)
	valInt := 0
	if val != "" {
		valInt, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("Failed to convert val to int %+v", err)
			return err
		}
	}
	log.Printf("ValInt: %+v", valInt)
	if valInt < times {
		t.RedisService.Client.Incr(key)
		t.RedisService.Client.Expire(key, time.Duration(hours)*time.Hour)
	} else {
		return errors.New("Limit reached")
	}
	return nil
}
