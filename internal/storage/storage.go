package storage

import "os"

const DataPath string = "data"

func CreateDataDirectoryIfNotExists() error {
	return os.MkdirAll(DataPath, os.ModePerm)
}
