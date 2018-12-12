package server

import "encoding/json"

type KeyValues struct {
	*KeyInfo	`json:"key_info"`
	Value   interface{} `json:"value"`  
}


func (values *KeyValues) String() string {
	bs,err := json.Marshal(values)
	if err != nil {
		return ""
	}
	return string(bs[:])
}



//type String struct {
//	Value string `json:"value"`
//}
//
//type Zset struct {
//
//}
//
//type ZsetValue struct {
//	Member string `json:"member"`
//	Score  int64  `json:"score"`
//}