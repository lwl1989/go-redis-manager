package server

import (
	"github.com/go-redis/redis"
	"sync"
	"strconv"
)


type redisConnections map[string]*redisConnection
type redisConnection struct {
	*redis.Client
	conf *RedisConfig
	kMap KeyMap
}

var onceRedis sync.Once
var connections *redisConnections
var redisCon *redisConnection


//初始化不同配置的连接池
func Init() *redisConnections {
	onceRedis.Do(func() {
		connections := make(map[string]*redisConnection)
		for _,conf := range RedisHosts {
			hval := conf.GetHval()
			connections[hval] = GetRedis(hval)
		}
	})

	return connections
}


func GetRedis(hval string)  *redisConnection{

	conf,err := RedisHosts.GetConfig(hval)
	if  err != nil {
		panic(err.Error())
	}

	redisCon = &redisConnection{
		conf:conf,
	}
	redisCon.Client = redis.NewClient(&redis.Options{
		Addr:    conf.Host,
		Password: conf.Pw, // no password set
		DB:      conf.Db,  // use default DB
	})

	_, err1 := redisCon.Ping().Result()
	if err1 != nil {
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

//select dbs and keys to mem
func (conn *redisConnection) initKeys() {
	s := conn.ConfigGet("database")

	val,err := s.Result()

	db := 1
	if err == nil {
		if len(val) > 0 {
			db,err = strconv.Atoi(val[1].(string))
			if err != nil {
				db = 1
			}
		}
	}

	conn.kMap = GetKeyMap()
	for i:=1; i<=db; i++ {

		s := conn.Keys("*")
		keys,err := s.Result()
		kList := make([]*KeyInfo,0)
		if err == nil {
			//use pipe line to get keys info
			p:=conn.Pipeline()

			for _,key := range keys {
				info := GetKeyInfoWithInfo(key, i)
				kList = append(kList, info)
				//p.Process(conn.Type(key))
				p.Process(conn.TTL(key))
			}
			c,e := p.Exec()
			p.Close()


			if e == nil {
				i := 0
				for _,c1 :=range c {
					switch c1.(type) {
					case *redis.StatusCmd:
						//c2 := c1.(*redis.StatusCmd)
						//fmt.Println(c2.Val())
					case *redis.DurationCmd:
						c2 := c1.(*redis.DurationCmd)
						info := kList[i]
						info.SetTtlWithTime(c2.Val())
						i++
					}
				}
			}

		}
		conn.kMap[i] = kList
	}
}