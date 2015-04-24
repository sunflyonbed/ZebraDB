package main

import (
	"bytes"
	"fmt"
	"net"
	"runtime/debug"

	l4g "log4go"

	"redis"
)

func RedisServer(ip string) {
	addr, err := net.ResolveTCPAddr("tcp", ip)
	if err != nil {
		l4g.Error("resolve tcp addr error: %s", err.Error())
		return
	}
	l, e := net.ListenTCP("tcp", addr)
	if e != nil {
		l4g.Error("listen tcp error: %s", e.Error())
		return
	}

	defer l.Close()

	for {
		rw, e := l.AcceptTCP()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}
			l4g.Error("accept tcp error: %s", e.Error())
			return
		}
		go RedisProcess(rw)
	}
}

func RedisProcess(rw *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {
			l4g.Error("redis client recover: %s \nstack\n %s", err, debug.Stack())
		}
	}()
	client := redis.NewClient(rw)
	defer client.Close()
	for {
		reply := client.ReadRequest()
		data, err := reply.ListBytes()
		if err != nil {
			l4g.Error("redis client read request error: %s", err.Error())
			return
		}
		dl := len(data)
		CMD := bytes.ToUpper(data[0])
		var e error
		if bytes.Equal(CMD, REDIS_OP_HGETALL) { //HASH
			if dl < 2 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.HGetAll(data[1]))
			}
		} else if bytes.Equal(CMD, REDIS_OP_HKEYS) {
			if dl < 2 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.HKeys(data[1]))
			}
		} else if bytes.Equal(CMD, REDIS_OP_HMGET) {
			if dl < 3 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.HMGet(data[1:]))
			}
		} else if bytes.Equal(CMD, REDIS_OP_SMEMBERS) { //SET
			if dl < 2 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.SMembers(data[1]))
			}
		} else if bytes.Equal(CMD, REDIS_OP_ZRANGE) { //SORTEDSET
			if dl < 2 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.ZRange(data[1]))
			}
		} else if bytes.Equal(CMD, CMD_OP_SELECT) { //CMD
			e = client.WriteRespone("ok")
		} else if bytes.Equal(CMD, CMD_OP_INFO) {
			if dl < 2 {
				e = client.WriteRespone("param no enough\n")
			} else {
				e = client.WriteRespone(gDB.Info(data[1]))
			}
		} else if bytes.Equal(CMD, CMD_OP_DUMP) {
			if dl < 2 {
				e = client.WriteRespone("param no enough")
			} else {
				derr := gDB.Dump(data[1])
				if derr != nil {
					e = client.WriteRespone(derr.Error())
				} else {
					e = client.WriteRespone("ok")
				}
			}
		} else if bytes.Equal(CMD, CMD_OP_SIZE) {
			if dl%2 == 0 {
				e = client.WriteRespone("param no enough")
			} else {
				e = client.WriteRespone(gDB.Size(data[1:]))
			}
		} else {
			l4g.Error("redis %s cmd no found", CMD)
			client.WriteRespone(fmt.Sprintf("redis %s cmd no found", CMD))
			return
		}
		if e != nil {
			l4g.Error("%s write respone error: %s", CMD, e.Error())
			return
		}
	}
}
