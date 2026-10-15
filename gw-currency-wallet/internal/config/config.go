package config

import "os"

type Config struct {
	JWT         []byte
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	KafkaBroker string
}

func New() Config {
	return Config{
		JWT:         []byte(os.Getenv("JWT_SECRET")),
		Host:        os.Getenv("POSTGRES_HOST"),
		Port:        os.Getenv("POSTGRES_PORT"),
		User:        os.Getenv("POSTGRES_USER"),
		Password:    os.Getenv("POSTGRES_PASSWORD"),
		DBName:      os.Getenv("POSTGRES_DB"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
	}
}
