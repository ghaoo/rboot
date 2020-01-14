package bolt

import (
	"os"
	"path"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/ghaoo/rboot"
	"github.com/sirupsen/logrus"
)

const DefaultBoltDBFile = `db/rboot.db`

type boltMemory struct {
	bolt   *bolt.DB
	dbfile string

	mu sync.Mutex
}

// 检查存放db文件的文件夹是否存在。
// 如果db文件存在运行目录下，则无操作
func pathExist(dbfile string) error {

	dir, _ := path.Split(dbfile)

	if dir != `` {
		if _, err := os.Stat(dir); !os.IsExist(err) {
			return os.MkdirAll(dir, os.ModePerm)
		}
	}

	return nil
}

// new bolt brain ...
func Bolt() rboot.Brain {

	b := new(boltMemory)

	dbfile := os.Getenv(`BOLT_DB_FILE`)

	if dbfile == `` {
		logrus.Warningf(`BOLT_DB_FILE not set, using default: %s`, DefaultBoltDBFile)
		dbfile = DefaultBoltDBFile
	}
	err := pathExist(dbfile)

	if err != nil {
		return nil
	}

	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		return nil
	}

	b.bolt = db
	b.dbfile = dbfile
	return b
}

// save ...
func (b *boltMemory) Set(bucket, key string, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.bolt.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			logrus.Error("bolt: error saving:", e)
			return e
		}
		return b.Put([]byte(key), value)
	})

	return err
}

// find ...
func (b *boltMemory) Get(bucket, key string) []byte {
	b.mu.Lock()
	defer b.mu.Unlock()

	var found []byte
	b.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		found = b.Get([]byte(key))

		return nil
	})

	return found
}

// remove ...
func (b *boltMemory) Remove(bucket, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})

	return err
}

// register brain ...
func init() {
	rboot.RegisterBrain(`bolt`, Bolt)
}
