package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"strconv"
	"reflect"
)

func main()  {
	conn :=redis.NewClient(&redis.Options{
		Addr:    "127.0.0.1:6379",
		Password: "", // no password set
		DB:      0,  // use default DB
	})

	s := conn.ConfigGet("databases")
	val,_ := s.Result()
	if len(val) > 0 {
		db,_ := strconv.Atoi(val[1].(string))
		fmt.Println(db)
	}
	fmt.Println(val)
	fmt.Println(s.Val())

	s1 := conn.Keys("*")
	keys,err := s1.Result()
	if err == nil {
		p:=conn.Pipeline()
		for _, key := range keys {
			p.Process(conn.Type(key))
			p.Process(conn.TTL(key))

		}
		c,e := p.Exec()
		p.Close()



		if e != nil {
			fmt.Println(e)
		}else{
			for _,c1 :=range c {
				cc := reflect.TypeOf(c1)
				fmt.Println(cc)
				switch vvvv:=c1.(type){
				case *redis.StatusCmd:
					c2 := c1.(*redis.StatusCmd)
					fmt.Println(vvvv)
					fmt.Println(c2.Val())
				}
				//fmt.Println(c1.Name())
				//fmt.Println(c1.Args())
				//s2 := (*redis.StatusCmd)(unsafe.Pointer(&c1))
				////s2 := c1.(*redis.StatusCmd)
				//fmt.Println(s2.Val())
				//fmt.Println(c1)

				break
			}
		}

	}


}