package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

const configPath = "./config/local.yaml"

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Storage    `yaml:"storage"`
}

type HTTPServer struct {
	Addres      string        `yaml:"addres" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost:8080"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"db_name" env-required:"true"`
	SslMode  string `yaml:"ssl_mode" env-default:"disable"`
}

func Load() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read konfig: %s", err)
	}
	return &cfg
}
