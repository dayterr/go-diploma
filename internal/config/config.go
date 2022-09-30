package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	RunAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8081"`
	DatabaseURI   string        `env:"DATABASE_URI" envDefault:""`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
}

type FlagStruct struct {
	RunAddress string
	DatabaseURI   string
	AccrualSystemAddress string
}

func GetConfig() (Config, error) {
	log.Println("reading config")
	cfg := Config{}
	fs := FlagStruct{}
	flag.StringVar(&fs.RunAddress, "a", "", "RunAddress")
	flag.StringVar(&fs.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&fs.AccrualSystemAddress, "r", "", "Accrual System Address")
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