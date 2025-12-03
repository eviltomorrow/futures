package procutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eviltomorrow/futures/lib/buildinfo"
	"github.com/eviltomorrow/futures/lib/fs"
	"github.com/eviltomorrow/futures/lib/system"
)

func CreatePidFile() (func() error, error) {
	file, err := fs.CreateFlockFile(filepath.Join(system.Directory.VarDir(), "run", fmt.Sprintf("%s.pid", buildinfo.AppName)))
	if err != nil {
		return nil, err
	}
	if file == nil {
		return func() error { return nil }, nil
	}

	file.WriteString(fmt.Sprintf("%d", os.Getpid()))
	if err := file.Sync(); err != nil {
		file.Close()
		return nil, err
	}

	return func() error {
		if file != nil {
			if err := file.Close(); err != nil {
				return err
			}
			return os.Remove(filepath.Join(system.Directory.VarDir(), "run", fmt.Sprintf("%s.pid", buildinfo.AppName)))
		}
		return fmt.Errorf("panic: pid file is nil")
	}, nil
}
