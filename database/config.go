package database

import "github.com/pdcgo/shared/pkg/secret"

type Config struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func GetProductionClickhouseConfig() (*Config, error) {
	var cfg Config
	var sec *secret.Secret
	var err error

	sec, err = secret.GetSecret("clickhouse_credential", "latest")
	if err != nil {
		return nil, err
	}

	err = sec.YamlDecode(&cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
