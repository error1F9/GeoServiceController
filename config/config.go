package config

import (
	"os"
	"strconv"
	"time"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Server struct {
	Port string
}

type JWTKey struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type DadataCreds struct {
	ApiKey    string
	SecretKey string
}

type AppConfig struct {
	DB          DB
	Server      Server
	JWTKey      JWTKey
	DadataCreds DadataCreds
}

func NewAppConfig() *AppConfig {
	accessTTL, _ := strconv.Atoi(os.Getenv("ACCESS_TTL"))
	refreshTTL, _ := strconv.Atoi(os.Getenv("REFRESH_TTL"))
	return &AppConfig{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Server: Server{
			Port: os.Getenv("PORT"),
		},
		JWTKey: JWTKey{
			AccessSecret:  os.Getenv("JWT_PRIVATE_KEY"),
			RefreshSecret: os.Getenv("JWT_REFRESH_KEY"),
			AccessTTL:     time.Duration(accessTTL) * time.Minute,
			RefreshTTL:    time.Duration(refreshTTL) * time.Hour,
		},
		DadataCreds: DadataCreds{
			ApiKey:    os.Getenv("API_KEY"),
			SecretKey: os.Getenv("SECRET_KEY"),
		},
	}
}
