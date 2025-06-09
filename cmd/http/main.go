package main

import (
	"GeoService/config"
	_ "GeoService/internal/docs"
	"GeoService/run"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

// @title Geo Service
// @version 1.0
// @description Geo Service
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /

func main() {
	err := godotenv.Load()
	conf := config.NewAppConfig()
	logger, _ := zap.NewProduction()
	if err != nil {
		logger.Fatal("Error loading .env", zap.Error(err))
	}
	app := run.NewApp(*conf, logger)

	exitCode := app.
		Bootstrap().
		Run()
	os.Exit(exitCode)

}
