package config

import (
	"flag"
	"github.com/gisvr/golib/config"
	"github.com/gisvr/golib/log"
)

var (
	cfg Config
)

type Config struct {
	ConfigFile string      `yaml:"configFile"`
	TimeFmt    string      `yaml:"timeFmt"`
	SysEnvVar  interface{} `yaml:"sysEnv"`
	Log        *log.Config `yaml:"log"`
}

func Get() *Config {
	c := config.Get()
	if c == nil {
		c = config.Init(&cfg)
	}
	cfg := c.(*Config)
	return cfg
}

func init() {
	flag.StringVar(&config.ConfigFile, "c", ".gdeploy.yml", "config file")
}
