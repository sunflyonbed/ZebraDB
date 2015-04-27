package main

import (
	"common"
	"levigo"
)

type LevelDB struct {
	options  *levigo.Options
	woptions *levigo.WriteOptions
	db       *levigo.DB
}

func NewLevelDB(name string) (*LevelDB, error) {
	options := levigo.NewOptions()

	options.SetCreateIfMissing(true)

	options.SetWriteBufferSize(64 * 1024 * 1024)
	options.SetCompression(levigo.SnappyCompression)

	woptions := levigo.NewWriteOptions()
	woptions.SetSync(false)

	db, err := levigo.Open(name, options)

	if err != nil {
		return nil, err
	}
	ret := &LevelDB{options,
		woptions,
		db}
	return ret, nil
}

func (this *LevelDB) Write(wb *levigo.WriteBatch) error {
	return this.db.Write(this.woptions, wb)
}

func (this *LevelDB) Close() {
	if this.db != nil {
		this.db.Close()
	}

	if this.woptions != nil {
		this.woptions.Close()
	}

	if this.options != nil {
		this.options.Close()
	}
}

func (this *LevelDB) HMSet(key []byte, data [][]byte) {
	dl := len(data)
	pairs := dl / 2

	wb := levigo.NewWriteBatch()
	for i := 0; i < pairs; i++ {
		wb.Put(common.EncodeHashKey(key, data[2*i]), data[1+2*i])
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		panic(err.Error())
	}
}

func (this *LevelDB) SAdd(key []byte, data [][]byte) {
	dl := len(data)

	wb := levigo.NewWriteBatch()
	for i := 0; i < dl; i++ {
		wb.Put(common.EncodeSetKey(key, data[i]), nil)
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		panic(err.Error())
	}
}

func (this *LevelDB) ZAdd(key []byte, data [][]byte) {
	dl := len(data)
	pairs := dl / 2

	wb := levigo.NewWriteBatch()
	for i := 0; i < pairs; i++ {
		wb.Put(common.EncodeSortedSetKey(key, data[2*i]), data[1+2*i])
	}
	err := this.Write(wb)
	wb.Close()

	if err != nil {
		panic(err.Error())
	}
}
