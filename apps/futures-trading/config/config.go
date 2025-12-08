package config

import (
	"path/filepath"

	"github.com/eviltomorrow/futures/lib/flagsutil"
	"github.com/eviltomorrow/futures/lib/system"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type Config struct {
	Workspace Workspace `json:"workspace" toml:"workspace" mapstructure:"workspace"`
	Log       Log       `json:"log" toml:"log" mapstructure:"log"`
	IP        IP        `json:"ip" toml:"ip" mapstructure:"ip"`
}

type Workspace struct {
	AppDir string `json:"app_dir" toml:"app-dir" mapstructure:"app-dir"`
	UsrDir string `json:"usr_dir" toml:"usr-dir" mapstructure:"usr-dir"`
	VarDir string `json:"var_dir" toml:"var-dir" mapstructure:"var-dir"`
	LogDir string `json:"log_dir" toml:"log-dir" mapstructure:"log-dir"`
}

type Log struct {
	DisableStdlog bool `json:"disable_stdlog" toml:"disable-stdlog" mapstructure:"disable-stdlog"`
	MaxSize       int  `json:"max_size" toml:"max-size" mapstructure:"max-size"`
	MaxDays       int  `json:"max_days" toml:"max-days" mapstructure:"max-days"`
	MaxBackups    int  `json:"max_backups" toml:"max-backups" mapstructure:"max-backups"`

	Level string `json:"level" toml:"level" mapstructure:"level"`
}

type IP struct {
	AccessIP string `json:"access_ip" toml:"access-ip" mapstructure:"access-ip"`
	BindIP   string `json:"bind_ip" toml:"bind-ip" mapstructure:"bind-ip"`
}

func (c *Config) ResetSystemWithConfig() *Config {
	system.Directory.SetAppDir(c.Workspace.AppDir)
	system.Directory.SetUsrDir(c.Workspace.UsrDir)
	system.Directory.SetVarDir(c.Workspace.VarDir)
	system.Directory.SetLogDir(c.Workspace.LogDir)

	return c
}

func (c *Config) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func ReadConfigFromFile(opts *flagsutil.Flags) (*Config, error) {
	path := filepath.Join(system.Directory.ExecDir(), "config.toml")
	if opts.ConfigFile != "" {
		path = opts.ConfigFile
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	if opts.DisableStdlog {
		c.Log.DisableStdlog = true
	}

	return mergeConfigWithDefault(c, initDefaultConfig()), nil
}

type presettable interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~bool
}

func presetValue[T presettable](targetValue, defaultValue T) T {
	var zero T
	if targetValue == zero {
		return defaultValue
	}
	return targetValue
}

func mergeConfigWithDefault(targetConfig, defaultConfig *Config) *Config {
	targetConfig.Workspace.AppDir = presetValue(targetConfig.Workspace.AppDir, defaultConfig.Workspace.AppDir)
	targetConfig.Workspace.UsrDir = presetValue(targetConfig.Workspace.UsrDir, defaultConfig.Workspace.UsrDir)
	targetConfig.Workspace.VarDir = presetValue(targetConfig.Workspace.VarDir, defaultConfig.Workspace.VarDir)
	targetConfig.Workspace.LogDir = presetValue(targetConfig.Workspace.LogDir, defaultConfig.Workspace.LogDir)

	targetConfig.Log.Level = presetValue(targetConfig.Log.Level, defaultConfig.Log.Level)
	targetConfig.Log.MaxSize = presetValue(targetConfig.Log.MaxSize, defaultConfig.Log.MaxSize)
	targetConfig.Log.MaxDays = presetValue(targetConfig.Log.MaxDays, defaultConfig.Log.MaxDays)
	targetConfig.Log.MaxBackups = presetValue(targetConfig.Log.MaxBackups, defaultConfig.Log.MaxBackups)
	targetConfig.Log.DisableStdlog = presetValue(targetConfig.Log.DisableStdlog, defaultConfig.Log.DisableStdlog)

	targetConfig.IP.AccessIP = presetValue(targetConfig.IP.AccessIP, defaultConfig.IP.AccessIP)
	targetConfig.IP.BindIP = presetValue(targetConfig.IP.BindIP, defaultConfig.IP.BindIP)

	return targetConfig
}

func initDefaultConfig() *Config {
	c := &Config{
		Workspace: Workspace{
			AppDir: system.Directory.AppDir(),
			UsrDir: system.Directory.UsrDir(),
			VarDir: system.Directory.VarDir(),
			LogDir: system.Directory.LogDir(),
		},
		Log: Log{
			Level:         "info",
			MaxSize:       100,
			MaxDays:       180,
			MaxBackups:    180,
			DisableStdlog: false,
		},
		IP: IP{
			AccessIP: "",
			BindIP:   "0.0.0.0",
		},
	}
	return c
}
