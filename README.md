[ZebraDB-基于LevelDB存储数据的Redis服务](https://github.com/ivanabc/ZebraDB)
===========================

####Why
1. dump耗内存
2. aof恢复慢

####Install
```
./all.bash
```
####Run Redis (关闭redis需要检查队列长度)
```
./deps/redis-2.8.19/src/redis-server ./deps/redis-2.8.19/redis.conf
```
####Run Zebra
```
./bin/start_zebra.sh 
```
####Stop Zebra
```
./bin/stop_zebra.sh
```
###ZebraDB支持的Redis协议操作

| HASH       | SET       | SORTEDSET  | SERVER |
| --------   | --------- | ---------  | ------ |
| HGETALL    | SMEMBERS  | ZRANGE     | DUMP   |
| HKEYS      | SADD      | ZADD       | INFO   |
| HMGET      | SREM      | ZREM       | SIZE   |
| HSET       |           |            |        |
| HMSET      |           |            |        |
| HDEL       |           |            |        |
| DEL        |           |            |        |
