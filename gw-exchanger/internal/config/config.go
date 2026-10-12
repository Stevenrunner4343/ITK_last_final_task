package config

import "os"

var (
	Host     = os.Getenv("POSTGRES_HOST")
	Port     = os.Getenv("POSTGRES_PORT")
	User     = os.Getenv("POSTGRES_USER")
	Password = os.Getenv("POSTGRES_PASSWORD")
	DBName   = os.Getenv("POSTGRES_DB")
)
