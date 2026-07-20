package database

import (
	"diyd/src/config"
	"diyd/src/harddisk/utils"
)

// KeyValueStore interface type
type KeyValueStore interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type Params func(kv *kvstore) error

func NewKVStore(params ...Params) (KeyValueStore, error) {
	kv := &kvstore{
		table: &Table{},
	}
	for _, param := range params {
		err := param(kv)
		if err != nil {
			return nil, err
		}
	}

	return kv, nil
}

func WithCreateTable(tname string) Params {
	return func(kv *kvstore) error {
		return utils.CreateNewTable(config.DefaultDirectoryPath, tname)
	}
}

func WithTable(tname string) Params {
	return func(kv *kvstore) error {
		kv.table.Directory = utils.GetTablePath(config.DefaultDirectoryPath, tname)
		kv.table.Name = tname
		return nil
	}
}
