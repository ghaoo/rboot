package rboot

import "fmt"

type Memorizer interface {
	Save(bucket, key string, value []byte) error
	Find(bucket, key string) []byte
	FindAll(bucket string) map[string][]byte
	Update(bucket, key string, value []byte) error
	Delete(bucket, key string) error
}

var memorizers = make(map[string]func() Memorizer)

// 注册存储器
func RegisterMemorizer(name string, m func() Memorizer) {

	if name == "" {
		panic("RegisterMemorizer: memorizer must have a name")
	}
	if _, ok := memorizers[name]; ok {
		panic("RegisterMemorizer: memorizers named " + name + " already registered. ")
	}
	memorizers[name] = m
}

func DetectMemorizer(name string) (func() Memorizer, error) {
	if memo, ok := memorizers[name]; ok {
		return memo, nil
	}

	if len(memorizers) == 0 {
		return nil, fmt.Errorf("no memorizer available")
	}

	if name == "" {
		if len(memorizers) == 1 {
			for _, memo := range memorizers {
				return memo, nil
			}
		}
		return nil, fmt.Errorf("multiple memorizers available; must choose one")
	}
	return nil, fmt.Errorf("unknown memorizer '%s'", name)
}
