package leveldb

import (
	"os"
	"path"

	"github.com/ghaoo/rboot"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

const DefaultDBFile = `.data/db/level_weasyad.db`

type goleveldb struct {
	db *leveldb.DB
}

func NewLevelDB() rboot.Brain {
	dbfile := os.Getenv(`LEVELDB_FILE`)

	if dbfile == "" {
		dbfile = DefaultDBFile
	}

	err := pathExist(dbfile)

	if err != nil {
		return nil
	}

	db, err := leveldb.OpenFile(dbfile, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func": `newDB`,
		}).WithError(err).Error("leveldb 初始化失败")
		return nil
	}

	return &goleveldb{
		db: db,
	}
}

// Set ...
func (level *goleveldb) Set(bucket, key string, val []byte) error {
	key = bucket + key
	return level.db.Put([]byte(key), val, nil)
}

// Get ...
func (level *goleveldb) Get(bucket, key string) ([]byte, error) {
	key = bucket + key
	found, err := level.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func (level *goleveldb) Remove(bucket, key string) error {
	key = bucket + key
	return level.db.Delete([]byte(key), nil)
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

// register brain ...
func init() {
	rboot.RegisterBrain(`leveldb`, NewLevelDB)
}
