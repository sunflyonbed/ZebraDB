package main

import (
	"levigo"

	l4g "log4go"
)

type LevelDB struct {
	env      *levigo.Env
	options  *levigo.Options
	roptions *levigo.ReadOptions
	woptions *levigo.WriteOptions
	db       *levigo.DB
}

func (this *LevelDB) Open(dbname string) (err error) {
	if this.db != nil {
		return
	}

	this.db, err = levigo.Open(dbname, this.options)
	return
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

func (this *LevelDB) Close() {
	if this.db != nil {
		this.db.Close()
	}

	if this.options != nil {
		this.options.Close()
	}

	if this.roptions != nil {
		this.roptions.Close()
	}

	if this.woptions != nil {
		this.woptions.Close()
	}

	if this.env != nil {
		this.env.Close()
	}
}

func (this *LevelDB) NewIterator() *levigo.Iterator {
	return this.db.NewIterator(this.roptions)
}

func (this *LevelDB) NewIteratorWithReadOptions(roptions *levigo.ReadOptions) *levigo.Iterator {
	return this.db.NewIterator(roptions)
}

func NewLevelDB(name string) (*LevelDB, error) {
	options := levigo.NewOptions()

	// options.SetComparator(cmp)
	options.SetCreateIfMissing(true)
	options.SetErrorIfExists(false)

	// set env
	env := levigo.NewDefaultEnv()
	options.SetEnv(env)

	// set cache
	cache := levigo.NewLRUCache(512 << 20)
	options.SetCache(cache)

	options.SetInfoLog(nil)
	options.SetParanoidChecks(false)
	options.SetWriteBufferSize(32 << 20)
	options.SetMaxOpenFiles(500)
	options.SetBlockSize(4 * 1024)
	options.SetBlockRestartInterval(16)
	options.SetCompression(levigo.SnappyCompression)

	// set filter
	filter := levigo.NewBloomFilter(10)
	options.SetFilterPolicy(filter)

	roptions := levigo.NewReadOptions()
	roptions.SetVerifyChecksums(false)
	roptions.SetFillCache(true)

	woptions := levigo.NewWriteOptions()
	// set sync false
	woptions.SetSync(false)

	db := &LevelDB{env,
		options,
		roptions,
		woptions,
		nil}
	if err := db.Open(name); err != nil {
		l4g.Error("open db faileddb, path: %s err: %s", name, err.Error())
		return nil, err
	}
	l4g.Info("open db succeed, path: %s", name)
	return db, nil
}
