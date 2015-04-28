[ZebraDB - Redis数据保存服务](https://bitbucket.org/ivanzt/zebradb)
===========================

####Why
1. dump耗内存
2. aof恢复慢

####Install: (directory: ZebraDB)
```
./all.bash
```
####Run: (directory: bin,默认redis已启动)
```
 ./start_zebra.sh 
```
####Stop: (directory: bin)
```
./stop_zebra.sh
```
####Tools: (directory: bin)
```
./save -h=true
```
```
./restore -h=true
```
####Test: (directory: bin)
```
./redisprotocol -i="HSET T a 1" | redis-cli -p 6381 -n 0 --pipe
```
```
redis-cli -p 9999
```
###ZebraDB支持的Redis指令

####HASH
```
HGETALL  HKEYS HMGET HSET  HMSET
HDEL DEL
```
####SET
```
SMEMBERS SADD
SREM
```
####SORTEDSET
```
ZRANGE ZADD
ZREM
```
####SERVER
```
DUMP
INFO
SIZE
```