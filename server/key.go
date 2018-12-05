package server

//redis type
const REDIS_NONE   = 0
const REDIS_STRING = 1
const REDIS_LIST   = 2
const REDIS_SET    = 3
const REDIS_ZSET   = 4
const REDIS_HASH   = 5
const REDIS_UNKNOW = 9

//db key map lists
type KeyMap map[int][]*KeyInfo

//key
type KeyInfo struct {
	Db		int    `json:"db"`   //which db
	KeyName string `json:"key_name"`  //keyname
	Type 	uint8  `json:"type"`      //key type
	TTl		int    `json:"ttl"`
	checkExists  bool
}

func GetKeyInfo() *KeyInfo  {
	return &KeyInfo{
		Db: 1,
		KeyName:"",
		Type:REDIS_UNKNOW,
		TTl: -1,
		checkExists:false,
	}
}

func GetKeyInfoWithInfo(key string, db int) *KeyInfo  {
	return &KeyInfo{
		Db: db,
		KeyName:key,
		Type:REDIS_UNKNOW, //default string
		TTl: -1,
		checkExists:true,
	}
}

func GetKeyMap() KeyMap  {
	//keys := make([]*KeyInfo,0)
	return make(KeyMap)
}

func (maps KeyMap) GetDbKeys(db int) []*KeyInfo {
	keyList,ok := maps[db]
	if ok {
		return keyList
	}

	return nil
}

func (info *KeyInfo) GetDb() int {
	return info.Db
}

func (info *KeyInfo) GetKeyName() string {
	return info.KeyName
}

func (info *KeyInfo) GetType() uint8 {
	return info.Type
}

func (info *KeyInfo) GetTypeString() string {
	switch info.Type {
	case REDIS_HASH:
		return "hash"
	case REDIS_LIST:
		return "list"
	case REDIS_SET:
		return "set"
	case REDIS_ZSET:
		return "zset"
	case REDIS_STRING:
		return "string"
	case REDIS_NONE:
		fallthrough  //fallthrough不会判断下一条case的expr结果是否为true。 就是没有break
		//但是如果几个条件都走一样的结果，使用 fallthrough串起来即可
	default:
		return "none"
	}
}

func (info *KeyInfo) GetTtl() int {
	return info.TTl
}

