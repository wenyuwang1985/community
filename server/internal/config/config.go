package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Snowflake SnowflakeConfig
	JWT       JWTConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr string
	DB   int
}

type SnowflakeConfig struct {
	NodeID int64 `mapstructure:"node_id"`
}

type JWTConfig struct {
	Secret        string
	AccessExpire  int `mapstructure:"access_expire"`
	RefreshExpire int `mapstructure:"refresh_expire"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
