package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

//PG_URL=postgres://{user}:{password}@{host}:{port}/{database}

type (
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"log"`
		Storage `yaml:"storage"`
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
