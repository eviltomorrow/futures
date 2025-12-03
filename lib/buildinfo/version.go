package buildinfo

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

// -X with build
var (
	MainVersion string
	GoVersion   = runtime.Version()
	GoOSArch    = runtime.GOOS + "/" + runtime.GOARCH
	GitSha      string
	BuildTime   string
)

func Version() string {
	return fmt.Sprintf("%s", MainVersion)
}

var (
	bluebold  = color.New(color.FgBlue, color.Bold)
	whitebold = color.New(color.FgWhite, color.Bold)
)

func VersionDetail() string {
	s1 := bluebold.Sprintf("Version: ")
	s2 := whitebold.Sprintf("%s %s (commit-id=%s)", AppName, MainVersion, GitSha)
	s3 := bluebold.Sprintf("Runtime: ")
	s4 := whitebold.Sprintf("%s %s RELEASE.%s", GoVersion, GoOSArch, BuildTime)
	return s1 + s2 + "\r\n" + s3 + s4
}
