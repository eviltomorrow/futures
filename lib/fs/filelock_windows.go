package fs

import (
	"os"
)

func CreateFlockFile(path string) (*os.File, error) {
	return createFlockFile(path)
}

func createFlockFile(_ string) (*os.File, error) {
	return nil, nil
}
