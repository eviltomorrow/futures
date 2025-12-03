package system

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eviltomorrow/futures/lib/netutil"
)

func LoadRuntime() error {
	executePath, err := os.Executable()
	if err != nil {
		return err
	}
	executePath, err = filepath.Abs(executePath)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(executePath)
	Directory.execDir = baseDir

	if strings.HasSuffix(baseDir, "/bin") {
		baseDir = filepath.Dir(baseDir)
	}
	Directory.rootDir = baseDir

	Directory.SetAppDir(filepath.Join(Directory.rootDir, "/app"))
	Directory.SetUsrDir(filepath.Join(Directory.rootDir, "/usr"))
	Directory.SetVarDir(filepath.Join(Directory.rootDir, "/var"))
	Directory.SetLogDir(filepath.Join(Directory.rootDir, "/log"))
	Directory.SetEtcDir(filepath.Join(Directory.rootDir, "/etc"))

	Process.name = filepath.Base(executePath)
	Process.args = os.Args[1:]
	Process.pid = os.Getpid()
	Process.ppid = os.Getppid()

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	Machine.hostname = hostname

	ipv4, err := netutil.GetInterfaceIPv4First()
	if err != nil {
		return err
	}
	Network.bindIP = ipv4

	return nil
}
