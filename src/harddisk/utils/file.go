package utils

import (
	"diyd/src/config"
	"os"
	"path"
)

// CreateNewTable creates new file in the given directory
func CreateNewTable(dirName string, table string) error {
	err := os.MkdirAll(path.Join(os.Getenv(config.DIYDStorageEnvVar), dirName), os.ModePerm)
	if err != nil {
		return err
	}

	fPath := GetTablePath(dirName, table)
	f, err := os.Create(fPath)
	defer f.Close()
	return err
}

// GetTablePath file path for the given directory
func GetTablePath(dirName string, table string) string {
	return path.Join(os.Getenv(config.DIYDStorageEnvVar), dirName, table+".data")
}
