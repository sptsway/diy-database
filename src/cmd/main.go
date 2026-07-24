package main

import (
	"diyd/src/cmd/server"
	"diyd/src/config"
	"flag"
)

func parseArgs() config.CmdArgs {
	args := config.CmdArgs{}
	args.KVName = *flag.String("kvstore", "default_kvstore", "use -kvstore <name> to select the store")
	args.Create = *flag.Bool("create", false, "use -create to create a new kvstore")
	args.Port = *flag.Int("port", 8080, "use -port <no>")
	args.MaxCons = *flag.Int("maxconn", 1000, "use -maxconn <no>")
	flag.Parse()

	return args
}

func main() {
	args := parseArgs()
	s := &server.Server{
		Port: args.Port,
	}
	s.Start()
}
