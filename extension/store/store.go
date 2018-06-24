package store

import (
	"time"
	"github.com/ghaoo/rboot"
)

type memorizer interface {
	Save(key string, value []byte)
	Read(key string) []byte
	Remove(key string) error
}

// stores 储存器集，用于注册自定义事件
var stores map[string]memorizer

// 注册储存器
func Register(name string, mem memorizer) {

	/*if name == "" {
		panic("Register memorizer: memorizer must have a name")
	}
	if _, ok := stores[name]; ok {
		panic("Register memorizer: memorizer named " + name + " already registered. ")
	}

	stores[name] = mem*/

}

// Store 储存器数据结构
type Store struct {
	Key      string
	Value    []byte
	CreateAt time.Time
}

//
func Put(name, key string, val []byte) {
	store := Store{
		Key: key,

	}

	rboot.SendCustomEvent("/store/put/"+name, )
}