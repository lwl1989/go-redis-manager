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
}