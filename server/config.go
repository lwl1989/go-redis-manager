package server

import (
	"time"
	"crypto/sha256"
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
		RConf.hval = string(sha256.Sum256([]byte(time.Now().String() + RConf.Host))[:])
	}
	return RConf.hval
}