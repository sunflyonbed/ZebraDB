package main

import (
	"errors"
	"fmt"

	l4g "log4go"

	"levigo"
)

//LIKE REDIS READ CMD
var (
	//HASH CMD
	REDIS_OP_HGETALL = []byte("HGETALL")
	REDIS_OP_HKEYS   = []byte("HKEYS")
	REDIS_OP_HMGET   = []byte("HMGET")
	//SET CMD
	REDIS_OP_SMEMBERS = []byte("SMEMBERS")
	//SORTEDSET CMD
	REDIS_OP_ZRANGE = []byte("ZRANGE")
)

//LIKE REDIS WRITE CMD
var (
	//HASH CMD
	REDIS_OP_HSET  = []byte("HSET")
	REDIS_OP_HMSET = []byte("HMSET")
	REDIS_OP_HDEL  = []byte("HDEL")
	REDIS_OP_DEL   = []byte("DEL") //ONLY DELETE HASH
	//SET CMD
	REDIS_OP_SADD = []byte("SADD")
	REDIS_OP_SREM = []byte("SREM")
	//SORTEDSET CMD
	REDIS_OP_ZADD = []byte("ZADD")
	REDIS_OP_ZREM = []byte("ZREM")
)

var (
	CMD_OP_SELECT = []byte("SELECT")
	CMD_OP_INFO   = []byte("INFO")
	CMD_OP_DUMP   = []byte("DUMP")
	CMD_OP_SIZE   = []byte("SIZE")
)

func encodeHashKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = 'h'
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func encodeSetKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = 's'
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func encodeSortedSetKey(name, key []byte) []byte {
	ret := make([]byte, 1+1+len(name)+1+len(key))
	ret[0] = 'z'
	ret[1] = byte(len(name))
	copy(ret[2:], name)
	ret[2+len(name)] = '='
	copy(ret[3+len(name):], key)
	return ret
}

func decodeHashKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != 'h' {
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

func decodeSetKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != 's' {
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

func decodeSortedSetKey(data []byte) (name, key []byte, ret bool) {
	if len(data) < 5 {
		return nil, nil, false
	}
	if data[0] != 'z' {
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

//HASH FUNCTION
func (this *LevelDB) HSet(data [][]byte) bool {
	if len(data) != 3 {
		l4g.Error("hset len error: %d", len(data))
		return false
	}
	err := this.Put(encodeHashKey(data[0], data[1]), data[2])
	if err != nil {
		l4g.Error("hset %s %s %s write error: %s", data[0], data[1], data[2], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) HMSet(data [][]byte) bool {
	dl := len(data)
	if dl < 3 {
		l4g.Error("hmset len error: %d", dl)
		return false
	}
	if dl/2 == 0 {
		l4g.Error("hmset param error: %d", dl)
		return false
	}
	pairs := (dl - 1) / 2

	wb := levigo.NewWriteBatch()
	for i := 0; i < pairs; i++ {
		wb.Put(encodeHashKey(data[0], data[1+2*i]), data[2+2*i])
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s hmset write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) HDel(data [][]byte) bool {
	dl := len(data)
	if dl < 2 {
		if dl == 1 {
			l4g.Error("hdel %s len error: %d", data[0], dl)
		} else {
			l4g.Error("hdel len error: %d", dl)
		}
		return false
	}

	wb := levigo.NewWriteBatch()
	for i := 1; i < dl; i++ {
		wb.Delete(encodeHashKey(data[0], data[i]))
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s hdel write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) HClear(data [][]byte) bool {
	for _, v := range data {
		var tmp [][]byte
		tmp = append(tmp, v)
		tmp = append(tmp, this.HKeys(v)...)
		this.HDel(tmp)
	}
	return true
}

func (this *LevelDB) HKeys(data []byte) (retList [][]byte) {
	it := this.NewIterator()
	defer it.Close()
	start := encodeHashKey(data, []byte(nil))
	for it.Seek(start); it.Valid(); it.Next() {
		name, key, ret := decodeHashKey(it.Key())
		if ret {
			if string(name) == string(data) {
				retList = append(retList, key)
			} else {
				break
			}
		} else {
			break
		}
	}
	if err := it.GetError(); err != nil {
		l4g.Error("hgetall %s error: %s", data, err.Error())
		retList = nil
		retList = append(retList, []byte(err.Error()))
	}
	return
}

func (this *LevelDB) HGetAll(data []byte) (retList [][]byte) {
	it := this.NewIterator()
	defer it.Close()
	start := encodeHashKey(data, []byte(nil))
	for it.Seek(start); it.Valid(); it.Next() {
		name, key, ret := decodeHashKey(it.Key())
		if ret {
			if string(name) == string(data) {
				retList = append(retList, key)
				retList = append(retList, it.Value())
			} else {
				break
			}
		} else {
			break
		}
	}
	if err := it.GetError(); err != nil {
		l4g.Error("hgetall %s error: %s", data, err.Error())
		retList = nil
		retList = append(retList, []byte(err.Error()))
	}
	return
}

func (this *LevelDB) HMGet(data [][]byte) (retList [][]byte) {
	for _, v := range data[1:] {
		value, err := this.Get(encodeHashKey(data[0], v))
		if err != nil {
			retList = nil
			retList = append(retList, []byte(err.Error()))
			return
		} else {
			retList = append(retList, value)
		}
	}
	return
}

//SET FUNCTION
func (this *LevelDB) SAdd(data [][]byte) bool {
	dl := len(data)
	if dl < 2 {
		l4g.Error("sadd len error: %d", len(data))
		return false
	}

	wb := levigo.NewWriteBatch()
	for i := 1; i < dl; i++ {
		wb.Put(encodeSetKey(data[0], data[i]), nil)
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s sadd write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) SRem(data [][]byte) bool {
	dl := len(data)
	if dl < 2 {
		l4g.Error("srem len error: %d", len(data))
		return false
	}

	wb := levigo.NewWriteBatch()
	for i := 1; i < dl; i++ {
		wb.Delete(encodeSetKey(data[0], data[i]))
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s srem write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) SMembers(data []byte) (retList [][]byte) {
	it := this.NewIterator()
	defer it.Close()
	start := encodeSetKey(data, []byte(nil))
	for it.Seek(start); it.Valid(); it.Next() {
		name, key, ret := decodeSetKey(it.Key())
		if ret {
			if string(name) == string(data) {
				retList = append(retList, key)
			} else {
				break
			}
		} else {
			break
		}
	}
	if err := it.GetError(); err != nil {
		l4g.Error("smembers %s error: %s", data, err.Error())
		retList = nil
		retList = append(retList, []byte(err.Error()))
	}
	return
}

//SORTEDSET FUNCTION
func (this *LevelDB) ZAdd(data [][]byte) bool {
	dl := len(data)
	if dl < 3 {
		l4g.Error("zadd len error: %d", dl)
		return false
	}
	if dl/2 == 0 {
		l4g.Error("zadd param error: %d", dl)
		return false
	}
	pairs := (dl - 1) / 2

	wb := levigo.NewWriteBatch()
	for i := 0; i < pairs; i++ {
		wb.Put(encodeSortedSetKey(data[0], data[2+2*i]), data[1+2*i])
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s zadd write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) ZRem(data [][]byte) bool {
	dl := len(data)
	if dl < 2 {
		l4g.Error("zrem len error: %d", dl)
		return false
	}

	wb := levigo.NewWriteBatch()
	for i := 1; i < dl; i++ {
		wb.Delete(encodeSortedSetKey(data[0], data[i]))
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		l4g.Error("%s zrem write error: %s", data[0], err.Error())
		return false
	}
	return true
}

func (this *LevelDB) ZRange(data []byte) (retList [][]byte) {
	it := this.NewIterator()
	defer it.Close()
	start := encodeSortedSetKey(data, []byte(nil))
	for it.Seek(start); it.Valid(); it.Next() {
		name, key, ret := decodeSortedSetKey(it.Key())
		if ret {
			if string(name) == string(data) {
				retList = append(retList, key)
				retList = append(retList, it.Value())
			} else {
				break
			}
		} else {
			break
		}
	}
	if err := it.GetError(); err != nil {
		l4g.Error("zrange %s error: %s", data, err.Error())
		retList = nil
		retList = append(retList, []byte(err.Error()))
	}
	return
}

//COMMAND FUNCTION
func (this *LevelDB) Dump(data []byte) error {
	options := levigo.NewOptions()
	defer options.Close()

	options.SetCreateIfMissing(true)
	options.SetErrorIfExists(true)

	options.SetWriteBufferSize(gConf.LevelDB.WriteBufferSize * 1024 * 1024)
	options.SetCompression(levigo.SnappyCompression)

	newDB, err := levigo.Open(string(data), options)
	if err != nil {
		return err
	}
	defer newDB.Close()

	woptions := levigo.NewWriteOptions()
	defer woptions.Close()

	woptions.SetSync(false)

	roptions := levigo.NewReadOptions()
	defer roptions.Close()

	roptions.SetFillCache(false)

	it := this.NewIteratorWithReadOptions(roptions)
	defer it.Close()

	index := 0
	var wb *levigo.WriteBatch
	for it.SeekToFirst(); it.Valid(); it.Next() {
		if index == 0 {
			wb = levigo.NewWriteBatch()
		}
		wb.Put(it.Key(), it.Value())
		index++
		if index == 10000 {
			err := newDB.Write(woptions, wb)
			if err != nil {
				l4g.Error("dump write batch error: %s", err.Error())
				wb.Close()
				return err
			}
			index = 0
			wb.Close()
		}
	}
	if index > 0 {
		err := newDB.Write(woptions, wb)
		if err != nil {
			l4g.Error("dump write batch error: %s", err.Error())
			wb.Close()
			return err
		}
		wb.Close()
	}
	if err := it.GetError(); err != nil {
		l4g.Error("dump %s error: %s", data, err.Error())
		return errors.New(fmt.Sprintf("dump %s error: %s", data, err.Error()))
	}
	return nil
}

func (this *LevelDB) Info(key []byte) string {
	property := "leveldb." + string(key)
	prop := this.db.PropertyValue(property)
	if prop == "" {
		return "valid key:\n\tnum-files-at-level<N>\n\tstats\n\tsstables\n"
	}
	return prop + "\n"
}

func (this *LevelDB) Size(data [][]byte) []uint64 {
	rs := make([]levigo.Range, len(data)/2)
	for index, v := range data {
		if index%2 == 0 {
			rs[index/2].Start = v
		} else {
			rs[index/2].Limit = v
		}
	}
	return this.db.GetApproximateSizes(rs)
}
