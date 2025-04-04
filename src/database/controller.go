package database

import (
	"fmt"
	"os/exec"
)

type kvstore struct {
	table *Table
	// delimiter is used to separate key-value pairs
	delimiter []byte
	// kv separator is used to separate key and value
	kvSeparator []byte
}

// Get finds the dats associated with the key
// current impm in very inefficient,
// TODO many things, use indexing
func (kv *kvstore) Get(key string) ([]byte, error) {
	// Implementation for getting a value by key
	return exec.Command("bash", "-c", fmt.Sprintf(
		"grep %s %s | tail -n 1 | awk -F '=' '{print $2}'", key+"=", kv.table.Directory),
	).Output()
}

// Set currently only appends to the file
// TODO many things, implement indexing
func (kv *kvstore) Set(key string, value []byte) error {
	// append to EOF
	return exec.Command("bash", "-c", fmt.Sprintf("echo %s=%s >> %s", key, value, kv.table.Directory)).Run()
}

// Delete removes the key-value pair from the store
func (kv *kvstore) Delete(key string) error {
	// Implementation for deleting a key-value pair
	return kv.Set(key, nil)
}
