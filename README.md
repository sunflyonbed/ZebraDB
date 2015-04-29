[ZebraDB - 基于LevelDB存储数据的Redis服务](https://bitbucket.org/ivanzt/zebradb)
===========================

####Why
1. dump耗内存
2. aof恢复慢

####Install: (directory: ZebraDB)
```
./all.bash
```
####Run: (directory: bin; 默认redis已启动,监听6381端口)
```
./start_zebra.sh 
```
####Stop: (directory: bin)
```
./stop_zebra.sh
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
