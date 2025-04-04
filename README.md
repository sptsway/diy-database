# diy-database
do-it-yourself-database, build your own database

### <wip> todo, many things

build
```
$ export DIYD_STORAGE_DIR=$PWD
$ go build diyd/src/cmd
```

create table
```
go run diyd/src/cmd -table example -create
```

get key
```
go run diyd/src/cmd -table example -get <key>
```

set key
```
go run diyd/src/cmd -table example -get <key>=<value>
```

delete key
```
go run diyd/src/cmd -table example -delete <key>
```