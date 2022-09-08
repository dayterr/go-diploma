package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	DatabaseURI   string        `env:"DATABASE_URI" envDefault:""`
}

type FlagStruct struct {
	DatabaseURI   string
}

func GetConfig() (Config, error) {
	log.Println("reading config")
	cfg := Config{}
	fs := FlagStruct{}
	flag.StringVar(&fs.DatabaseURI, "d", "", "Database URI")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println("parsing env error", err)
		return Config{}, err
	}

	if cfg.DatabaseURI == "" && fs.DatabaseURI != "" {
		cfg.DatabaseURI = fs.DatabaseURI
	}
	return cfg, nil
}