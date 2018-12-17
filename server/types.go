package server

import (
	"encoding/json"
	"strconv"
)

type KeyValues struct {
	*KeyInfo	`json:"key_info"`
	Value   interface{} `json:"value"`
}

type ValueOf interface {
	Value() interface{}
	String() string
}

func (values *KeyValues) String() string {
	bs,err := json.Marshal(values)
	if err != nil {
		return ""
	}
	return string(bs[:])
}

type String struct {
	Values string `json:"value"`
}

type Zset struct {
	Values map[string]int64 `json:"value"`
}

type Set struct {
	Values []string `json:"value"`
}

type Hash struct {
	Values map[string]string `json:"value"`
}

type List struct {
	Values []string `json:"value"`
}

type Geo struct {
	Values []map[string]*GeoValue `json:"values"`
}

type GeoValue struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func GetValue(redisType int, value interface{}) ValueOf  {
	switch redisType {
		case RedisSet:
			m := value.(map[string]int64)
			return &Zset{
				Values:m,
			}
		case RedisZset:
			l := value.([]string)
			return &Set{
				Values:l,
			}
		case RedisString:
			s := value.(string)
			return &String{
				Values:s,
			}
		case RedisHash:
			h := value.(map[string]string)
			return &Hash{
				Values:h,
			}
		case RedisList:
			l := value.([]string)
			return &List{
				Values:l,
			}
	}

	return nil
}


func (zset *Zset) Value() interface{} {
	return zset.Values
}

func (zset *Zset) String() string  {
	str := ""
	for k,v := range zset.Values {
		str += "key:"+k+"\tvalue:"+strconv.FormatInt(v, 10)
	}
	return str
}


func (set *Set) Value() interface{} {
	return set.Values
}

func (set *Set) String() string  {
	str := ""
	for k,v := range set.Values {
		str += "key:"+strconv.FormatInt(int64(k), 10)+"\tvalue:"+v
	}
	return str
}


func (s *String) Value() interface{} {
	return s.Values
}

func (s *String) String() string  {
	return s.Values
}


func (hash *Hash) Value() interface{} {
	return hash.Values
}

func (hash *Hash) String() string  {
	str := ""
	for k,v := range hash.Values {
		str += "key:"+k+"\tvalue:"+v
	}
	return str
}


func (l *List) Value() interface{} {
	return l.Values
}

func (l *List) String() string  {
	str := ""
	for k,v := range l.Values {
		str += "key:"+strconv.FormatInt(int64(k), 10)+"\tvalue:"+v
	}
	return str
}


func (g *Geo) Value() interface{} {
	return g.Values
}

func (g *Geo) String() string  {
	str := ""
	for k,v := range g.Values {
		str += "key:"+strconv.FormatInt(int64(k), 10)
		for k1,v1  := range v {
			str += "\tname:"+k1
			str += "\tlat:"+strconv.FormatFloat(v1.Lat, 'f', -1, 64)
			str += "\tlng:"+strconv.FormatFloat(v1.Lng, 'f', -1, 64)
		}
	}
	return str
}