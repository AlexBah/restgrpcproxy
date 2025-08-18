package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string        `yaml:"env" env-required:"true"`
	TlsPath    string        `yaml:"tls_path" env-required:"true"`
	Port       int           `yaml:"port" env-required:"true"`
	Timeout    time.Duration `yaml:"timeout"`
	GRPCServer string        `yaml:"grpcserver" env-required:"true"`
}

// returns config from *.yaml
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

// returns config from *.yaml
func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	return &cfg
}

// returns path to config.yaml from the arguments of the executable file
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
