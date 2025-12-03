package flagsutil

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	jsoniter "github.com/json-iterator/go"
)

var Opts = &Flags{}

type Flags struct {
	ConfigFile    string `short:"c" long:"config-file" description:"specifying a config file"`
	Daemon        bool   `short:"d" long:"daemon" description:"running in background"`
	EnablePprof   bool   `long:"enable-pprof" description:"enable pprof profiling"`
	PprofAddr     string `long:"pprof-addr" default:":56060" description:"pprof listen addr"`
	DisableStdlog bool   `long:"disable-stdlog" description:"disable standard logging"`
	Mode          string `long:"mode" description:"run mode(debug, release)"`

	Version       bool `short:"v" long:"version" description:"show version number"`
	VersionDetail bool `long:"version-detail" description:"show version detail"`
}

func (f *Flags) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(f)
	if err != nil {
		return fmt.Sprintf("marshal metadata failure, nest error: %v", err)
	}
	return string(buf)
}

func Parse(opts *Flags) ([]string, error) {
	metadata, err := flags.NewParser(opts, flags.Default).Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
	return metadata, nil
}
