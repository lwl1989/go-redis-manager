package server

import (
	"time"
	"crypto/sha256"
	"errors"
)

type RedisVhosts map[string]*RedisConfig

type RedisConfig struct {
	Host string `json:"host,omitempty"`
	Db   int `json:"db,omitempty"`
	Pw   string `json:"pw,omitempty"`
	hval string
}

type HostConfig struct {
	Host string `json:"host,omitempty"`
	Port int `json:"port,omitempty"`
}

var RedisHosts RedisVhosts

func (RConf *RedisConfig)  GetHval() string {
	if RConf.hval == "" {
		bs := sha256.Sum256([]byte(time.Now().String() + RConf.Host))
		RConf.hval = string(bs[:])
	}
	return RConf.hval
}

func (hosts RedisVhosts) GetConfig(hval string) (*RedisConfig,error) {
	conf,ok := hosts[hval]
	if ok {
		return conf,nil
	}

	return nil,errors.New("config not found with "+ hval)
}