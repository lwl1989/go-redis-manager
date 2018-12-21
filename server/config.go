package server

import (
	"time"
	"crypto/sha256"
	"errors"
)

type RedisVhosts map[string]*RedisConfig

type RedisConfig struct {
	Name string `json:"name,omitempty"`
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

func (RConf *RedisConfig)  getName() string {
	if RConf.Name == "" {
		return RConf.Host
	}
	return RConf.Name
}


func (hosts RedisVhosts) GetConfig(hval string) (*RedisConfig,error) {
	conf,ok := hosts[hval]
	if ok {
		return conf,nil
	}

	return nil,errors.New("config not found with "+ hval)
}

func (hosts RedisVhosts) GetConfigByName(name string) (*RedisConfig,error) {
	for _,conf := range hosts {
		if conf.getName() == name {
			return conf,nil
		}
	}
	return nil,errors.New("config not found with "+ name)
}

func (hosts RedisVhosts) GetName(hval string) (string ,error) {
	conf, err := hosts.GetConfig(hval)
	if err != nil {
		return "",err
	}

	return conf.getName(), nil
}

func (hosts RedisVhosts) AddHost(RConf *RedisConfig) (error) {
	name := RConf.getName()
	for _,conf := range hosts {
		if conf.getName() == name {
			return  errors.New("config already exists with name(if name is nil, name is host:port):"+name)
		}
	}

	return  nil
}