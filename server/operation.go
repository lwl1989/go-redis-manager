package server

import (
	"net/url"
	"strconv"
)

//config
const SetConfig = "set_config"

//server
const Info		= "get_server_info"

//key
const GetKey    = "get_key"
const DelKey    = "del_key"
const SetTtl    = "set_ttl"

//string
const SetValue  = "set_value"

//set use it
const AddListValue  = "add_list_value"
const DelListValue  = "del_list_value"

const SetZsetField  = "set_zset_field"
const DelZsetField  = "del_zset_field"

const SetHashField  = "set_hash_field"
const DelHashField  = "del_hash_field"


func getServerInfo(conn *redisConnection) string {
	cmd := conn.Info()
	str,_ := cmd.Result()
	return str
}

func setConfig(values url.Values) error {
	host := values.Get("host")
	name := values.Get("name")
	pw   := values.Get("pw")

	return  AddConnections(&RedisConfig{
		Name:name,
		Host:host,
		Db:0,
		Pw:pw,
	})
}

func KeyOperation(conn *redisConnection, action string, values url.Values) (*KeyValues) {
	dbstr := values.Get("db")
	db := 0
	if dbstr != "" {
		db,_ = strconv.Atoi(dbstr)
	}

	conn.Do("select", db)

	key := values.Get("key_name")
	info := &KeyInfo{
		KeyName:key,
		Db:db,
		TTl:getTtl(conn, key),
		Type:getType(conn, key),
	}
	switch action {
		case GetKey:
			return getKey(conn, info)
		case DelKey:
		case SetTtl:
		case SetValue:
		case AddListValue:
		case DelListValue:
		case SetZsetField:
		case DelZsetField:
		case SetHashField:
		case DelHashField:
		default:
			return nil
	}
	return nil
}

func DoOperation(values url.Values) (bool, interface{}) {
	redisName := values.Get("redis_name")
	conf,err := RedisHosts.GetConfigByName(redisName)
	if err != nil {
		return false,nil
	}

	conn := GetRedisInConnections(conf.GetHval())

	action := values.Get("action")
	switch action {
		case SetConfig:
		case Info:
			return true,getServerInfo(conn)
		default:
			result := KeyOperation(conn, action, values)
			if result != nil {
				return true, result
			}
			return false, nil
	}
	return false, nil
}

func getType(conn *redisConnection, key string) uint8  {
	cmd := conn.Type(key)
	str,err := cmd.Result()
	if err != nil {
		return RedisUnknow
	}
	return getTypeWithString(str)
}

func getTtl(conn *redisConnection, key string) int64{
	cmd := conn.TTL(key)
	d,err := cmd.Result()
	if err != nil {
		return -1
	}

	return d.Nanoseconds()
}

func getKey(con *redisConnection, info *KeyInfo) *KeyValues {
	cmd := con.Get(info.GetKeyName())
	str,_ := cmd.Result()
	return &KeyValues{
		KeyInfo:info,
		Value: str,
	}
}



