package fs

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/eviltomorrow/futures/lib/system"
	"github.com/eviltomorrow/futures/lib/zlog"
	"go.uber.org/zap"
)

func RewriteStderrToFile() error {
	if err := MkdirAll(system.Directory.LogDir()); err != nil {
		return err
	}

	panicFile, err := os.OpenFile(filepath.Join(system.Directory.LogDir(), "panic.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	if err = syscall.Dup3(int(panicFile.Fd()), int(os.Stderr.Fd()), 0); err != nil {
		return err
	}
	runtime.SetFinalizer(panicFile, func(fd *os.File) {
		fd.Close()
	})

	return nil
}

func RecoverFromPanic() {
	if err := recover(); err != nil {
		zlog.Error("Recover from panic", zap.Any("error", err))
	}
}
