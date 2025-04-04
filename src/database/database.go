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

func WithNewTable(tname string) Params {
	return func(kv *kvstore) error {
		err := utils.CreateNewTable(config.DefaultDirectoryPath, tname)
		if err != nil {
			return err
		}
		return WithExistingTable(tname)(kv)
	}
}

func WithExistingTable(tname string) Params {
	return func(kv *kvstore) error {
		kv.table.Directory = utils.GetTablePath(config.DefaultDirectoryPath, tname)
		kv.table.Name = tname
		return nil
	}
}
