ZebraDB监听redis中的dbq队列, 获取dbq中数据(redis部分写入指令的协议)并将其翻译为LevelDB格式保存

* Install: (directory: ZebraDB)
>1. ./all.bash
>2. mkdir log var
>3. 修改 start_zebra.sh stop_zebra.sh zebra_config.xml zebra_log.xml中的路径
* Run: (directory: bin)
>./start_zebra.sh (默认redis已启动)
* Stop: (directory: bin)
>./stop_zebra.sh
* Test: (directory: bin)
>./redisprotocol -i="HSET T a 1" | redis-cli -p 6381 -n 0 --pipe

ZebraDB支持的Redis指令

* HASH
>HGETALL  HKEYS HMGET HSET  HMSET
>HDEL DEL
* SET
>SMEMBERS SADD
>SREM
* SORTEDSET
>ZRANGE ZADD
>ZREM
