package rboot

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

// Brain 是Rboot缓存器实现的接口
type Brain interface {
	Set(bucket, key string, value []byte) error
	Get(bucket, key string) ([]byte, error)
	GetBucket() ([]string, error)
	GetKeys(bucket string) ([]string, error)
	Remove(bucket, key string) error
}

var brains = make(map[string]func() Brain)

// RegisterBrain 注册存储器，名称须唯一
// 需实现Brain接口
func RegisterBrain(name string, m func() Brain) {

	if name == "" {
		panic("RegisterBrain: brain must have a name")
	}
	if _, ok := brains[name]; ok {
		panic("RegisterBrain: brains named " + name + " already registered. ")
	}
	brains[name] = m
}

// DetectBrain 获取名称为 name 的缓存器
func DetectBrain(name string) (func() Brain, error) {
	if brain, ok := brains[name]; ok {
		return brain, nil
	}

	if len(brains) == 0 {
		return nil, fmt.Errorf("no Brain available")
	}

	if name == "" {
		if len(brains) == 1 {
			for _, brain := range brains {
				return brain, nil
			}
		}
		return nil, fmt.Errorf("multiple brains available; must choose one")
	}
	return nil, fmt.Errorf("unknown brain '%s'", name)
}

// Store 向储存器中存入信息
func (bot *Robot) Store(bucket, key string, value []byte) error {
	return bot.Brain.Set(bucket, key, value)
}

// Find 从储存器中获取指定的bucket和key对应的信息
func (bot *Robot) Find(bucket, key string) ([]byte, error) {
	return bot.Brain.Get(bucket, key)
}

// Remove 从储存器中移除指定的bucket和key对应的信息
func (bot *Robot) Remove(bucket, key string) error {
	return bot.Brain.Remove(bucket, key)
}

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
func Bolt() Brain {

	b := new(boltMemory)

	dbfile := os.Getenv(`BOLT_DB_FILE`)

	if dbfile == `` {
		logrus.Warningf(`BOLT_DB_FILE not set, using default: %s`, DefaultBoltDBFile)
		dbfile = DefaultBoltDBFile
	}

	dbfile = path.Join(os.Getenv("DATA_PATH"), dbfile)
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
			logrus.WithFields(logrus.Fields{
				"mod": `rboot`,
			}).Error("bolt: error saving:", e)
			return e
		}
		return b.Put([]byte(key), value)
	})

	return err
}

// find ...
func (b *boltMemory) Get(bucket, key string) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var found []byte
	err := b.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		found = b.Get([]byte(key))

		return nil
	})

	return found, err
}

// GetBucket ...
func (b *boltMemory) GetBucket() ([]string, error) {
	var buckets []string
	err := b.bolt.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	return buckets, err
}

// GetKeys ...
func (b *boltMemory) GetKeys(bucket string) ([]string, error) {
	var numKeys = 0
	var keys []string
	err := b.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket does not exist")
		}
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			numKeys++
		}

		keys = make([]string, numKeys)
		numKeys = 0
		c = b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys[numKeys] = string(k)
			numKeys++
		}
		return nil
	})

	return keys, err
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
	RegisterBrain(`bolt`, Bolt)
}
