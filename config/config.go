package config

import (
	"os"
	"github.com/naoina/toml"
)

type Config struct {
	SMS struct {
		Accesskey string
		Privatekey string
		Serviceid string
	}

	Token struct {
		Secret string
	}

	Server struct {
		Mode string
		Port string
		DBname string
		Host string
	}

	DB map[string]map[string]string
}

func GetConfig(fpath string) *Config {
	c := new(Config)
	if file, err := os.Open(fpath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return c
		}
	}
}