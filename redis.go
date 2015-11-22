package margelet

import (
	"gopkg.in/redis.v3"
)

var Redis *redis.Client

func InitRedis(addr string, password string, db int64) {

	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
