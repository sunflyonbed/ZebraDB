package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"levigo"
	"redis"
)

const (
	MAX_LINK         = 64
	MAX_REDIS_OP_NUM = 128
)

var (
	help     = flag.String("h", "false", "help")
	leveldb  = flag.String("l", "ZebraDB_PATH/var", "leveldb data path")
	redisdb  = flag.String("r", "127.0.0.1:6381", "redis ip:port")
	selectdb = flag.Int("n", 0, "redis select db number")
)

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

	println("******       database restore          ******")
	println(time.Now().String())
	println("param: -l=", *leveldb)
	println("param: -r=", *redisdb)
	println("param: -s=", *selectdb)

	var redisCs [MAX_LINK]*redis.Client

	for i := 0; i < MAX_LINK; i++ {
		if client, err := redis.Dial("tcp", *redisdb); err != nil {
			panic(err.Error())
		} else {
			reply := client.Cmd("select", *selectdb)
			if reply.Err != nil {
				panic(reply.Err.Error())
			}
			redisCs[i] = client
			defer redisCs[i].Close()
		}
	}

	println("******   start leveldb sync  redis     ******")
	println(time.Now().String())
	CopyLevelDBToRedis(*leveldb, redisCs)
	println("******   finish leveldb sync redis     ******")

	println(time.Now().String())
	println("******   database restore  successful  ******")
}

func writeRedis(client *redis.Client, queue chan *redis.Requests, wg *sync.WaitGroup) {
	defer wg.Done()
	for rs := range queue {
		client.SetPending(rs)
		reply := client.GetReply()
		for ; reply.Err == nil; reply = client.GetReply() {
		}
		if reply.Err != nil {
			if reply.Err != redis.PipelineQueueEmptyError {
				panic(fmt.Sprintf("get reply error: %s", reply.Err.Error()))
			}
		}
	}
}

func CopyLevelDBToRedis(dir string, clients [MAX_LINK]*redis.Client) {
	var wg sync.WaitGroup
	var chans [MAX_LINK]chan *redis.Requests
	for index, client := range clients {
		ch := make(chan *redis.Requests, 1024)
		chans[index] = ch
		wg.Add(1)
		go writeRedis(client, ch, &wg)
	}

	options := levigo.NewOptions()
	defer options.Close()
	options.SetCreateIfMissing(false)
	db, err := levigo.Open(dir, options)
	if err != nil {
		panic(fmt.Sprintf("open leveldb fail: %s %s", dir, err.Error()))
	}
	defer db.Close()

	ropt := levigo.NewReadOptions()
	defer ropt.Close()
	ropt.SetFillCache(false)
	it := db.NewIterator(ropt)
	defer it.Close()

	//flush redis
	println("******   start flush redis   ******")
	if reply := clients[0].Cmd("FLUSHDB"); reply.Err != nil {
		panic(reply.Err.Error())
	}
	println("******   finish flush redis  ******")

	hcount, scount, zcount, num, index := 0, 0, 0, 0, 0
	var rs *redis.Requests
	for it.SeekToFirst(); it.Valid(); it.Next() {
		if num == 0 {
			rs = &redis.Requests{}
		}
		if it.Key()[0] == 'h' {
			k, f, ret := decodeKey(it.Key())
			if ret == false {
				return
			}
			rs.Append("hset", k, f, it.Value())
			hcount++
			num++
		} else if it.Key()[0] == 's' {
			k, f, ret := decodeKey(it.Key())
			if ret == false {
				return
			}
			rs.Append("sadd", k, f)
			scount++
			num++
		} else if it.Key()[0] == 'z' {
			k, f, ret := decodeKey(it.Key())
			if ret == false {
				return
			}
			rs.Append("zadd", k, it.Value(), f)
			zcount++
			num++
		} else {
			continue
		}
		if num == MAX_REDIS_OP_NUM {
			tmp := rs
			chans[index] <- tmp
			index++
			if index == MAX_LINK {
				index = 0
			}
			num = 0
		}

		if hcount%1000000 == 0 {
			println("leveldb hash field:      ", hcount)
		}
		if scount%10000000 == 0 {
			println("leveldb set field:       ", scount)
		}
		if zcount%10000000 == 0 {
			println("leveldb sortedset field: ", zcount)
		}
	}
	if num > 0 {
		chans[index] <- rs
	}

	println("leveldb hash field:      ", hcount)
	println("leveldb set field:       ", scount)
	println("leveldb sortedset field: ", zcount)

	if err := it.GetError(); err != nil {
		panic(fmt.Sprintf("leveldb  iterator error: %s", err.Error()))
	}
	for _, ch := range chans {
		close(ch)
	}
	wg.Wait()
}

func decodeKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	nameLen := int(data[1])
	if len(data)-4 < nameLen {
		return nil, nil, false
	}
	name = data[2 : 2+nameLen]
	key = data[3+nameLen:]
	return name, key, true
}
