package main

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env  string `env:"ENV"  envDefault:"development"`
	Port string `env:"PORT" envDefault:"8080"`
	DB   struct {
		URL string `env:"DATABASE_URL"`
	}
}

func Load() Config {
	return env.Must(env.ParseAs[Config]())
}
