###ZebraDB监听redis中的dbq队列,将dbq中数据(redis部分写入指令的协议)翻译为LevelDB格式保存

####Aim
1. 替代redis提供的数据落地方式
2. 支持redis协议访问ZebraDB中数据

####Why
1. dump方式耗内存
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