package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func GetSeverConfig() ServerConfig {
	envCfg, err := ParseEnv()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := ParseArgs()
	if err != nil {
		log.Fatalln(err)
	}

	newConfig := ServerConfig{}

	if envCfg.IsDefaultAddress() && !cfg.IsDefaultAddress() {
		newConfig.Address = cfg.Address
	} else {
		newConfig.Address = envCfg.Address
	}

	if envCfg.IsDefaultDatabaseURI() && !cfg.IsDefaultDatabaseURI() {
		newConfig.DatabaseURI = cfg.DatabaseURI
	} else {
		newConfig.DatabaseURI = envCfg.DatabaseURI
	}

	if envCfg.IsDefaultAccrualAddress() && !cfg.IsDefaultAccrualAddress() {
		newConfig.AccrualAddress = cfg.AccrualAddress
	} else {
		newConfig.AccrualAddress = envCfg.AccrualAddress
	}

	if envCfg.IsDefaultHashKey() && !cfg.IsDefaultHashKey() {
		newConfig.HashKey = cfg.HashKey
	} else {
		newConfig.HashKey = envCfg.HashKey
	}

	return newConfig
}

func ParseArgs() (ServerConfig, error) {
	var addr, db, accrual, key string

	flag.StringVarP(&addr, "address", "a", "localhost:8081", "Host address")
	flag.StringVarP(&db, "db_uri", "d", "postgres://test:test@localhost:5432/gophermart?sslmode=disable", "Database URI")
	flag.StringVarP(&accrual, "accrual", "r", "localhost:8080", "Accrual system address")
	flag.StringVarP(&key, "key", "k", "super_secret_key", "System hashing key")
	flag.Parse()

	return ServerConfig{
		Address:        addr,
		DatabaseURI:    db,
		AccrualAddress: accrual,
		HashKey:        key,
	}, nil
}

func ParseEnv() (ServerConfig, error) {
	var cfg ServerConfig

	err := env.Parse(&cfg)
	if err != nil {
		return ServerConfig{}, err
	}

	return cfg, nil
}
