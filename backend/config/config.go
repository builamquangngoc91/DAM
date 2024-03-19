package config

import (
	"fmt"
	"os"
)

type ApplicationConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", d.Host, d.User, d.Password, d.Database, d.Port)
}

type RedisConfig struct {
	Host string
	Port string
}

func (d *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", d.Host, d.Port)
}

type Config struct {
	Database    *DatabaseConfig
	Redis       *RedisConfig
	Application *ApplicationConfig
}

var Cfg Config

func LoadConfig() {
	ApplicationConfig := ApplicationConfig{
		Port: os.Getenv("DAM_PORT"),
	}

	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DAM_DB_HOST"),
		Port:     os.Getenv("DAM_DB_PORT"),
		User:     os.Getenv("DAM_DB_USER"),
		Password: os.Getenv("DAM_DB_PASSWORD"),
		Database: os.Getenv("DAM_DB_DATABASE"),
	}

	redisConfig := RedisConfig{
		Host: os.Getenv("DAM_REDIS_HOST"),
		Port: os.Getenv("DAM_REDIS_PORT"),
	}

	Cfg = Config{
		Database:    &dbConfig,
		Redis:       &redisConfig,
		Application: &ApplicationConfig,
	}
}
