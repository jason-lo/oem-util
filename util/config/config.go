package util

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   string
	Port     int
	Email    string
	User     string
	Password string
}

func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
