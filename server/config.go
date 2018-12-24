package server

import (
	"crypto/sha256"
	"errors"
	"time"
)

type RedisVhosts map[string]*RedisConfig

type RedisConfig struct {
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
	Db   int    `json:"db,omitempty"`
	Pw   string `json:"pw,omitempty"`
	Hval string `json:"hval"`
}

type HostConfig struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

var RedisHosts RedisVhosts

func InitConfig(host, name, pw string) *RedisConfig {
	conf := &RedisConfig{
		Name:name,
		Host:host,
		Db:0,
		Pw:pw,
	}
	conf.GetHval()
	return conf
}

func (RConf *RedisConfig) GetHval() string {
	if RConf.Hval == "" {
		bs := sha256.Sum256([]byte(RConf.Host + time.Now().String()))
		RConf.Hval = string(bs[:])
	}

	return RConf.Hval
}

func (RConf *RedisConfig) GetName() string {
	if RConf.Name == "" {
		return RConf.Host
	}
	return RConf.Name
}

func (hosts RedisVhosts) GetConfig(Hval string) (*RedisConfig, error) {
	conf, ok := hosts[Hval]
	if ok {
		return conf, nil
	}

	return nil, errors.New("config not found with " + Hval)
}

func (hosts RedisVhosts) GetConfigByName(name string) (*RedisConfig, error) {
	for _, conf := range hosts {
		if conf.GetName() == name {
			return conf, nil
		}
	}

	return nil, errors.New("config not found with " + name)
}

func (hosts RedisVhosts) GetName(Hval string) (string, error) {
	conf, err := hosts.GetConfig(Hval)
	if err != nil {
		return "", err
	}

	return conf.GetName(), nil
}

func (hosts RedisVhosts) Add(RConf *RedisConfig) (error) {
	for _, conf := range hosts {
		if conf.GetHval() == RConf.GetHval() {
			return errors.New("config already exists with name(if name is nil, name is host:port):" + RConf.GetName())
		}
	}

	hosts[RConf.GetHval()] = RConf
	return nil
}

func (hosts RedisVhosts) Remove(hval string) error {
	deleted := false

	for key, conf := range hosts {
		if conf.GetHval() == hval {
			delete(hosts, key)
			deleted = true
			break
		}
	}

	if !deleted {
		return errors.New("delete config error with hash:" + hval)
	}

	return nil
}

func (hosts RedisVhosts) Edit(hval string, New *RedisConfig) error {
	changed := false

	for key, conf := range hosts {
		if conf.GetHval() == hval {
			delete(hosts, key)
			hosts[New.GetHval()] = New
			changed = true
			break
		}
	}

	if !changed {
		return errors.New("change config error with hval:" + hval)
	}

	return nil
}
