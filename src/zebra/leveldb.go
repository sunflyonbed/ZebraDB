package main

import (
	"levigo"

	l4g "log4go"
)

type LevelDB struct {
	options  *levigo.Options
	cache    *levigo.Cache
	roptions *levigo.ReadOptions
	woptions *levigo.WriteOptions
	db       *levigo.DB
}

func NewLevelDB(name string) (*LevelDB, error) {
	options := levigo.NewOptions()

	options.SetCreateIfMissing(true)

	cache := levigo.NewLRUCache(gConf.LevelDB.CacheSize * 1024 * 1024)
	options.SetCache(cache)

	options.SetBlockSize(gConf.LevelDB.BlockSize * 1024)
	options.SetWriteBufferSize(gConf.LevelDB.WriteBufferSize * 1024 * 1024)
	options.SetMaxOpenFiles(gConf.LevelDB.MaxOpenFiles)
	options.SetCompression(levigo.SnappyCompression)

	filter := levigo.NewBloomFilter(10)
	options.SetFilterPolicy(filter)

	roptions := levigo.NewReadOptions()
	roptions.SetFillCache(true)

	woptions := levigo.NewWriteOptions()
	woptions.SetSync(false)

	db, err := levigo.Open(name, options)

	if err != nil {
		l4g.Error("open db failed, path: %s err: %s", name, err.Error())
		return nil, err
	}
	l4g.Info("open db succeed, path: %s", name)
	ret := &LevelDB{options,
		cache,
		roptions,
		woptions,
		db}
	return ret, nil
}

func (this *LevelDB) Put(key, value []byte) error {
	return this.db.Put(this.woptions, key, value)
}

func (this *LevelDB) Get(key []byte) ([]byte, error) {
	return this.db.Get(this.roptions, key)
}

func (this *LevelDB) Delete(key []byte) error {
	return this.db.Delete(this.woptions, key)
}

func (this *LevelDB) Write(wb *levigo.WriteBatch) error {
	return this.db.Write(this.woptions, wb)
}

func (this *LevelDB) NewIterator() *levigo.Iterator {
	return this.db.NewIterator(this.roptions)
}

func (this *LevelDB) NewIteratorWithReadOptions(roptions *levigo.ReadOptions) *levigo.Iterator {
	return this.db.NewIterator(roptions)
}

func (this *LevelDB) Close() {
	if this.db != nil {
		this.db.Close()
	}

	if this.roptions != nil {
		this.roptions.Close()
	}

	if this.woptions != nil {
		this.woptions.Close()
	}

	if this.cache != nil {
		this.cache.Close()
	}

	if this.options != nil {
		this.options.Close()
	}
}
