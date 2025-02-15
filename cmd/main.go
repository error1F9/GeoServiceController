package main

import (
	"GeoService/internal/controller"
	"GeoService/internal/service"
	_ "GeoService/pkg/docs"
	"GeoService/pkg/router"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"os"
	"time"
)

var TokenAuth *jwtauth.JWTAuth

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
	geoService := service.NewGeoService("90a5dd26d0ba58ea94f25f085aa113ad67f2af27", "eb3066ce98823788c54dafb9e5e66d87a3c92d9d")
	geoController := controller.NewGeoServiceController(*geoService)
	if err := router.Serve(geoController); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

const content = ``

func WorkerTest() {
	t := time.NewTicker(1 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}
