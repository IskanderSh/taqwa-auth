package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string     `yaml:"env" env-default:"local"`
	LogLevel string     `yaml:"log_level" env-default:"debug"`
	Token    Token      `yaml:"token"`
	GRPC     GRPCConfig `yaml:"grpc"`
	DB       DB         `yaml:"db"`
}

type Token struct {
	TTL    time.Duration `yaml:"ttl" env-required:"true"`
	Secret string        `yaml:"secret" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DB struct {
	Port int `yaml:"port"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	fmt.Println(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file is not exists")
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config file")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		fmt.Println("getting env value")
		res = os.Getenv("CONFIG_PATH")
		fmt.Printf("get env value: %s\n", res)
	}

	return res
}
