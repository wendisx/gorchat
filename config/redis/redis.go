package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/wendisx/gorchat/config"
	"github.com/wendisx/gorchat/internal/constant"
)

type RedisClient = redis.Client

func NewRedisClient(env config.Env) *redis.Client {
	var db int
	db, err := strconv.Atoi(env[constant.REDIS_DATABASE])
	if err != nil {
		// default db
		db = 6
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     env[constant.REDIS_URL],
		Password: env[constant.REDIS_PASSWORD],
		DB:       db,
	})
	c := context.Background()
	v, err := rdb.Ping(c).Result()
	if err != nil {
		log.Printf("[init] -- (config/redis) status: fail")
	} else {
		log.Printf("[init] -- (config/redis) status: success res: %s", v)
	}
	return rdb
}
