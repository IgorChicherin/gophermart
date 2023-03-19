package config

type ServerConfig struct {
	Address        string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseURI    string `env:"DATABASE_URI" envDefault:"postgres://test:test@localhost:5432/gophermart?sslmode=disable"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"localhost:8081"`
	HashKey        string `env:"HASH_KEY" envDefault:"super_secret_key"`
}

func (sc ServerConfig) IsDefaultAddress() bool {
	return sc.Address == "localhost:8080"
}

func (sc ServerConfig) IsDefaultDatabaseURI() bool {
	return sc.DatabaseURI == "postgres://test:test@localhost:5432/gophermart?sslmode=disable"
}

func (sc ServerConfig) IsDefaultAccrualAddress() bool {
	return sc.AccrualAddress == "localhost:8000"
}

func (sc ServerConfig) IsDefaultSettings() bool {
	return sc.IsDefaultAddress() &&
		sc.IsDefaultDatabaseURI() &&
		sc.IsDefaultAccrualAddress()
}

func (sc ServerConfig) IsDefaultHashKey() bool {
	return sc.HashKey == "super_secret_key"
}
