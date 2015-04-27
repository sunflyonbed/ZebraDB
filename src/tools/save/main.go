package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"

	"redis"
)

var (
	help     = flag.String("h", "false", "help")
	leveldb  = flag.String("l", "ZebraDB_PATH/var", "leveldb data path")
	redisdb  = flag.String("r", "127.0.0.1:6381", "redis ip:port")
	selectdb = flag.String("n", "0", "redis select db number")

	gDB *LevelDB
)

const CONNECT_NUM int = 64

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *help == "true" {
		println("restore help:")
		println("-l=leveldb数据库文件路径     default: -l=ZebraDB_PATH/var")
		println("-r=redis ip:port             default: -r=127.0.0.1:6381")
		println("-n=redis select db number    default: -n=0")
		return
	}

	println("param: -l=", *leveldb)
	println("param: -r=", *redisdb)
	println("param: -n=", *selectdb)

	if db, err := NewLevelDB(*leveldb); err != nil {
		panic(err.Error())
	} else {
		gDB = db
	}

	var redisCs [CONNECT_NUM]*redis.Client
	for index := 0; index < CONNECT_NUM; index++ {
		if client, err := redis.Dial("tcp", *redisdb); err != nil {
			panic(err.Error())
		} else {
			redisCs[index] = client
			reply := redisCs[index].Cmd("select", *selectdb)
			if reply.Err != nil {
				panic(reply.Err.Error())
			}
			defer redisCs[index].Close()
		}
	}

	reply := redisCs[0].Cmd("DEL", "dbq")
	if reply.Err != nil {
		panic(reply.Err.Error())
	}

	allKey, err := redisCs[0].Cmd("KEYS", "*").ListBytes()
	if err != nil {
		panic(err.Error())
	}

	var msg_wg sync.WaitGroup
	var msgs chan string = make(chan string, CONNECT_NUM)
	msg_wg.Add(1)
	go func() {
		for msg := range msgs {
			println(msg)
		}
		msg_wg.Done()
	}()

	var wg sync.WaitGroup
	keysCount := len(allKey)
	println("count keys:", keysCount)
	if keysCount <= CONNECT_NUM {
		wg.Add(1)
		go process(redisCs, allKey, &wg, msgs, 0, 0, keysCount)
	} else {
		segment := keysCount / CONNECT_NUM
		startIndex, endIndex := 0, segment
		for index := 0; index < CONNECT_NUM; index++ {
			wg.Add(1)
			if index != CONNECT_NUM-1 {
				go process(redisCs, allKey, &wg, msgs, index, startIndex, endIndex)
			} else {
				go process(redisCs, allKey, &wg, msgs, index, startIndex, keysCount)
			}
			startIndex, endIndex = endIndex, endIndex+segment
		}
	}
	wg.Wait()
	close(msgs)
	msg_wg.Wait()
}

func process(redisCs [CONNECT_NUM]*redis.Client, keys [][]byte, wg *sync.WaitGroup, msgs chan string, connectIndex, startIndex, endIndex int) {
	defer wg.Done()
	redisC := redisCs[connectIndex]
	allKey := keys[startIndex:endIndex]
	var hashSize, setSize, zsetSize int
	for _, key := range allKey {
		if keyType, err := redisC.Cmd("TYPE", key).Str(); err != nil {
			panic(err.Error())
		} else {
			if keyType == "hash" {
				if hashInfo, err := redisC.Cmd("HGETALL", key).ListBytes(); err != nil {
					panic(err.Error())
				} else {
					gDB.HMSet(key, hashInfo)
					hashSize++
				}
			} else if keyType == "set" {
				if setInfo, err := redisC.Cmd("SMEMBERS", key).ListBytes(); err != nil {
					panic(err.Error())
				} else {
					gDB.SAdd(key, setInfo)
					setSize++
				}
			} else if keyType == "zset" {
				if zsetInfo, err := redisC.Cmd("ZRANGE", key, 0, -1, "WITHSCORES").ListBytes(); err != nil {
					panic(err.Error())
				} else {
					gDB.ZAdd(key, zsetInfo)
					zsetSize++
				}
			}
		}
	}
	msgs <- fmt.Sprintf("index %d start %d end %d sum %d hash %d set %d zset %d", connectIndex, startIndex, endIndex, endIndex-startIndex, hashSize, setSize, zsetSize)
}
