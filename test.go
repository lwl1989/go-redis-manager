package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"strconv"
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

	for i := 0 ; i < 16 ; i++ {
		//command := redis.NewCmd()
		c := conn.Do("select", i)
		r,e := c.Result()
		fmt.Println(r.(string),e)
		s1 := conn.Keys("*")
		keys, err := s1.Result()
		if err == nil {
			p := conn.Pipeline()
			for _, key := range keys {
				p.Process(conn.Type(key))
				p.Process(conn.TTL(key))
			}
			c, e := p.Exec()
			p.Close()

			if e != nil {
				fmt.Println(e)
			} else {
				i := 0
				for _, c1 := range c {
					switch c1.(type) {
					case *redis.StatusCmd:
						c2 := c1.(*redis.StatusCmd)
						fmt.Println(c2.Val())
						fmt.Println(c2.Err())
					case *redis.DurationCmd:
						c2 := c1.(*redis.DurationCmd)
						fmt.Println(c2.Val())
						fmt.Println(c2.Args())
					}
					//fmt.Println(c1.Name())
					//fmt.Println(c1.Args())
					//s2 := (*redis.StatusCmd)(unsafe.Pointer(&c1))
					////s2 := c1.(*redis.StatusCmd)
					//fmt.Println(s2.Val())
					//fmt.Println(c1)
					i++
					if i > 10 {
						break
					}
				}
			}

		}
	}


}