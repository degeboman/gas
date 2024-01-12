package config

import (
	"flag"
	"github.com/degeboman/gas/constant"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	MongoConnectionString string
	Env                   string `yaml:"env" env-default:"local"`
	HTTPServer            `yaml:"http_server"`
	JwtSettings           `yaml:"jwt_settings"`
}

type JwtSettings struct {
	SigningKey      []byte
	AccessDuration  time.Duration `yaml:"access_duration" env-default:"300s"`
	RefreshDuration time.Duration `yaml:"access_duration" env-default:"604800s"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:2023"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() Config {
	configPath := flag.String(
		constant.ConfigPathFlag,
		"config/local.yml",
		constant.ConfigPathFlagUsage,
	)

	// check if file exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	// reading flags
	mongoConnectionString := flag.String(
		constant.MongoDBTokenFlagName,
		"",
		constant.MongoDBTokenFlagUsage,
	)

	jwtSigningKey := flag.String(
		constant.SigningKeyFlagName,
		"",
		constant.SigningKeyFlagUsage,
	)

	flag.Parse()

	// checking for flags
	if *mongoConnectionString == "" {
		log.Fatal("mongo connection string is not specified")
	}

	if *jwtSigningKey == "" {
		log.Fatal("jwt signing key is not specified")
	}

	var cfg Config

	cfg.MongoConnectionString = *mongoConnectionString
	cfg.JwtSettings.SigningKey = []byte(*jwtSigningKey)

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
