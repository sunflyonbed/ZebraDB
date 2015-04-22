package main

import (
	"flag"
	// "strings"
	"fmt"

	"redis"
	"redis/resp"
)

var (
	addr  = flag.String("a", "127.0.0.1:6381", "")
	index = flag.Int("i", 0, "")
	op    = flag.String("o", "", "")
)

func main() {
	flag.Parse()

	var dbClient *redis.Client
	if client, err := redis.Dial("tcp", *addr); err != nil {
		panic(fmt.Sprintf("connect redis %v", err))
	} else {
		dbClient = client
	}
	defer dbClient.Close()
	reply := dbClient.Cmd("SELECT", *index)
	if reply.Err != nil {
		panic(fmt.Sprintf("SELECT %s %d error", *addr, *index))
	}

	//println(string(resp.Format(data)))
	for i := 1000000; i < 2000000; i++ {
		//for i := 2; i < 3; i++ {
		var data []string
		data = append(data, "HSET")
		data = append(data, "T1")
		data = append(data, "a"+fmt.Sprintf("%d", i))
		data = append(data, fmt.Sprintf("%d", i))
		var cmd []string
		cmd = append(cmd, "RPUSH")
		cmd = append(cmd, "dbq")
		cmd = append(cmd, string(resp.Format(data)))
		print(string(resp.Format(cmd)))
		//dbClient.Cmd("RPUSH", "dbq", resp.Format(data))
	}
}
