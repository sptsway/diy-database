package main

import (
	"diyd/src/cmd/server"
	"diyd/src/database"
	"diyd/src/database/worker"
	"flag"
)

type cmdArgs struct {
	kvName   string
	create bool
	port     int
	maxConns int
}

func parseArgs() cmdArgs {
	args := cmdArgs{}
	args.kvName = *flag.String("kvstore", "default_kvstore", "use -kvstore <name> to select the store")
	args.create = *flag.Bool("create", false, "use -create to create a new kvstore")
	args.port = *flag.Int("port", 8080, "use -port <no>")
	args.maxConns = *flag.Int("maxconn", 1000, "use -maxconn <no>")
	flag.Parse()

	return args
}

func main() {
	args := parseArgs()

	dbparams := []database.Params{database.WithTable(args.kvName)}
	if args.create {
		dbparams = append(dbparams, database.WithCreateTable(args.kvName))
	}
	kvs, _ := database.NewKVStore(dbparams...)
	kvw := worker.NewKVWorker(worker.Params{
		WC:      10,
		KVStore: kvs,
	})
	s := &server.Server{
		Port:   args.port,
		Worker: kvw,
	}
	s.Start()
}
