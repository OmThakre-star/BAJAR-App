package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
}

type Config struct {
	Database  DatabaseConfig `yaml:"database"`
	JWTSecret string         `yaml:"jwt_secret"`
	Port      string         `yaml:"port"`
}

func Load() *Config {
	cfg, err := LoadYamlConfig("./config/userConfig.yaml")
	if err != nil {
		log.Fatalf("Could not load userConfig.yaml: %v", err)
	}
	return cfg
}

func LoadYamlConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) BuildDSN() string {
	db := c.Database
	return "host=" + db.Host +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.Name +
		" port=" + fmt.Sprintf("%d", db.Port) +
		" sslmode=" + db.SSLMode +
		" TimeZone=" + db.TimeZone
}
