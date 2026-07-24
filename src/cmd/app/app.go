package app

import (
	"diyd/src/config"
	"diyd/src/database"
	"diyd/src/database/worker"
	"errors"
	"io"
	"net/http"
)

type App interface {
	Ping(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Set(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewApp(cfg config.CmdArgs) (App, error) {

	dbparams := []database.Params{database.WithTable(cfg.KVName)}
	if cfg.Create {
		dbparams = append(dbparams, database.WithCreateTable(cfg.KVName))
	}
	kvs, err := database.NewKVStore(dbparams...)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to setup kvstore"))
	}
	kvw := worker.NewKVWorker(worker.Params{
		KVStore: kvs,
	})

	var closers []io.Closer
	for _, resource := range []any{kvw, kvs} {
		if c, ok := resource.(io.Closer); ok {
			closers = append(closers, c)
		}
	}

	return &kvapp{
		kvw:     kvw,
		closers: closers,
	}, nil
}

type kvapp struct {
	kvw     *worker.KVWorker
	closers []io.Closer
}

func (k *kvapp) Ping(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (k *kvapp) Get(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (k *kvapp) Set(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (k *kvapp) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (k *kvapp) Close() error {
	errs := []error{}
	for _, closer := range k.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
