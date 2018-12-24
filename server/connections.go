package server

import (
	"github.com/go-redis/redis"
	"sync"
	"strconv"
	"errors"
)


type redisConnections map[string]*redisConnection
type redisConnection struct {
	*redis.Client
	Conf *RedisConfig 	`json:"conf"`
	AllKeys [][]string	`json:"all_keys"`
	Err error			`json:"err"`
}

var onceRedis sync.Once
var connections redisConnections
var redisCon *redisConnection


//初始化不同配置的连接池
func Init() map[string]*redisConnection {
	onceRedis.Do(func() {
		connections := make(redisConnections)
		for _,conf := range RedisHosts {
			hval := conf.GetHval()
			connections[hval] = GetRedis(hval)
			conn := connections[hval]
			if conn.Err == nil {
				conn.initKeys()
			}
		}
	})

	return connections
}

func (rcons redisConnections) Remove(conn *redisConnection) error {
	deleted := false
	for key,conn := range rcons {
		if conn.Conf.GetHval() == conn.Conf.GetHval() {
			delete(rcons, key)
			deleted = true
			break
		}
	}

	if !deleted {
		return errors.New("remove connections error:"+conn.Conf.GetName())
	}

	return nil
}

func GetRedisInConnections(hval string) *redisConnection {
	conn,ok := connections[hval]
	if ok {
		return conn
	}

	connections[hval] = GetRedis(hval)

	return connections[hval]
}

func GetRedis(hval string)  *redisConnection{
	conn,ok := connections[hval]
	if ok {
		return conn
	}

	conf,err := RedisHosts.GetConfig(hval)
	if  err != nil {
		return &redisConnection{
			Err:errors.New("config not found with hval:"+hval),
		}
	}

	redisCon = &redisConnection{
		Conf:conf,
	}

	redisCon.Client = redis.NewClient(&redis.Options{
		Addr:    conf.Host,
		Password: conf.Pw, // no password set
		DB:      conf.Db,  // use default DB
	})

	i := 0
	//try again
	for ;i<3; i++ {
		_, err1 := redisCon.Ping().Result()
		if err1 != nil {
			redisCon.reConnection()
			redisCon.Err = err1
		}else{
			redisCon.Err = nil
		}
	}
	connections[conf.GetHval()] = redisCon
	return redisCon
}

func (conn *redisConnection) reConnection() {
	conf := conn.Conf
	redisCon.Client = redis.NewClient(&redis.Options{
		Addr:    conf.Host,
		Password: conf.Pw, // no password set
		DB:      conf.Db,  // use default DB
	})
}

//select dbs and keys to mem
func (conn *redisConnection) initKeys() {
	s := conn.ConfigGet("databases")

	val,err := s.Result()
	db := 0
	if err == nil {
		if len(val) > 0 {
			db,err = strconv.Atoi(val[1].(string))
			if err != nil {
				db = 0
			}
		}
	}

	conn.AllKeys = make([][]string,0)

	for i:=0; i<=db; i++ {
		conn.Do("select", i)
		s := conn.Keys("*")
		keys,err := s.Result()
		if err != nil {
			conn.Err = err
		}else {
			conn.AllKeys[i] = keys
		}
	}
}

func RemoveConnections(hval string) error {
	err := connections.Remove(GetRedis(hval))

	if err != nil {
		return err
	}

	err = RedisHosts.Remove(hval)

	if err != nil {
		return err
	}


	return nil
}

func AddConnections(config *RedisConfig) error {
	err := RedisHosts.Add(config)

	if err != nil {
		return err
	}

	client := GetRedis(config.GetHval())
	client.initKeys()

	if client.Err != nil {
		return client.Err
	}

	return nil
}

func ChangeConnection(hval string, conf *RedisConfig) error {
	err := RemoveConnections(hval)
	if err != nil {
		return err
	}

	return AddConnections(conf)
}