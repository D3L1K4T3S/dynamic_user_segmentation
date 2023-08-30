package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
	"time"
)

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"log"`
		Storage `yaml:"storage"`
		JWT     `yaml:"jwt"`
		Hash    `yaml:"hash"`
	}

	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Env     string `yaml:"env"`
	}

	Storage struct {
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size"`
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Database    string `yaml:"database"`
	}

	JWT struct {
		SignKey  string        `yaml:"sign_key"`
		TokenTTL time.Duration `yaml:"token_ttl"`
	}

	Hash struct {
		Salt string `yaml:"salt"`
	}
)

func MustLoadConfig(configPath string) *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		panic(err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func PgUrl(user, password, host, port, database string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)
}
