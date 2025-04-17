package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Cfg struct {
	DB         DbCfg         `yaml:"db"`
	HTTP       HTTPServerCfg `yaml:"http_server"`
	GRPC       GRPSCfg       `yaml:"grps"`
	Prometheus PrometheusCfg `yaml:"prometheus"`
	Auth       AuthCfg       `yaml:"auth"`
	Limits     LimitsCfg     `yaml:"limits"`
}

type DbCfg struct {
	Connection string `yaml:"connection"`
}

type HTTPServerCfg struct {
	Port    string `yaml:"http_port"`
	Timeout int    `yaml:"timeout_ms"`
}

type GRPSCfg struct {
	Port string `yaml:"grps_port"`
}

type PrometheusCfg struct {
	Port string `yaml:"prometheus_port"`
}

type AuthCfg struct {
	DummyTokenPrefix     string `yaml:"dummy_token_prefix"`
	JWTSecret            string `yaml:"jwt_secret"`
	JWTExpirationMinutes int    `yaml:"jwt_expiration_minutes"`
}

type LimitsCfg struct {
	PaginationLimit int `yaml:"pagination_limit"`
}

func GetConfig(path string) (Cfg, error) {
	var cfg Cfg

	file, err := os.Open(path)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to open config file")
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return cfg, errors.Wrap(err, "failed to decode config")
	}

	return cfg, nil
}
