package main

import (
	"diyd/src/database"
	"flag"
	"fmt"
	"strings"
)

func main() {
	run()
}

func run() {
	tableName := flag.String("table", "default_table", "use -table <name> to select the table")
	create := flag.Bool("create", false, "use -create to create a new table")
	getCmd := flag.String("get", "", "-get <key> to get the value of the key")
	setCmd := flag.String("set", "", "-set <key>=<value> to set the value of the key")
	deleteCmd := flag.String("delete", "", "-delete <key> to delete the key")
	flag.Parse()

	params := make([]database.Params, 0)
	if *create {
		params = append(params, database.WithNewTable(*tableName))
	} else {
		params = append(params, database.WithExistingTable(*tableName))
	}

	kvs, err := database.NewKVStore(params...)
	if err != nil {
		fmt.Printf("failed to create new kv store: %v", err)
		panic(err)
	}

	switch {
	case *getCmd != "":
		data, err := kvs.Get(*getCmd)
		if err != nil {
			fmt.Errorf("failed to get key %s: %e", *getCmd, err)
			panic(err)
		}
		fmt.Printf(string(data))
	case *setCmd != "":
		parts := strings.SplitN(*setCmd, "=", 2)
		key, val := parts[0], ""
		if len(parts) >= 2 {
			val = parts[1]
		}

		err = kvs.Set(key, []byte(val))
		if err != nil {
			fmt.Errorf("failed to set key %s: %e", key, err)
			panic(err)
		}
	case *deleteCmd != "":
		err = kvs.Delete(*deleteCmd)
		if err != nil {
			fmt.Errorf("failed to delete key %s: %e", *deleteCmd, err)
			panic(err)
		}
	default:
		if !*create {
			fmt.Errorf("no command provided, use -get, -set, or -delete")
		}
	}
}
