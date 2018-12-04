package server

import (
"github.com/go-redis/redis"
)

type redisConnection struct {
	*redis.Client
	conf *RedisConfig
}
var redisCon *redisConnection

func GetRedis(hval string)  *redisConnection{
	conf,ok := RedisHosts[hval]
	if  !ok {
		panic("hash value "+hval+" config not found")
	}

	redisCon = &redisConnection{
		conf:conf,
	}
	redisCon.Client = redis.NewClient(&redis.Options{
		Addr:    conf.Host,
		Password: conf.Pw, // no password set
		DB:      conf.Db,  // use default DB
	})

	_, err := redisCon.Ping().Result()
	if err != nil {
		redisCon.reConnection()
	}
	return redisCon
}

func (conn *redisConnection) reConnection() {
	conf := conn.conf
	redisCon.Client = redis.NewClient(&redis.Options{
		Addr:    conf.Host,
		Password: conf.Pw, // no password set
		DB:      conf.Db,  // use default DB
	})
}

func (conn *redisConnection) dConnection() {
	_, err := redisCon.Ping().Result()
	if err != nil {
		redisCon.reConnection()
	}
}
