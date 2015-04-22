ZebraDB 监听 redis 中的dbq队列, 获取dbq中数据(redis中hash和set部分写入指令的协议)并将其翻译为LevelDB格式保存.

* Install: (dir-ZebraDB)
>make && make tools
* Run: (dir-bin)
>./start_zebra.sh
* Stop: (dir-bin)
>./stop_zebra.sh
* Test: (dir-bin)
>./tools -i="HSET T a 1" | redis-cli -p 6381 -n 0 --pipe




