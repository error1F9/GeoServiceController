package main

import (
	"GeoService/internal/cacheproxy"
	cache2 "GeoService/internal/infrastructure/cache"
	"GeoService/internal/modules"
	"GeoService/internal/modules/address/controller"
	"GeoService/internal/modules/address/service"
	authController "GeoService/internal/modules/auth/controller"
	auth "GeoService/internal/modules/auth/service"
	_ "GeoService/internal/pkg/docs"
	"GeoService/internal/pkg/metrics"
	"GeoService/internal/router"
	"log"
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

	c, err := cache2.NewRedisCache("redis-container:6379", "123456")
	if err != nil {
		log.Fatal(err)
	}

	geoService := service.NewGeoService("90a5dd26d0ba58ea94f25f085aa113ad67f2af27", "eb3066ce98823788c54dafb9e5e66d87a3c92d9d")
	cache := cacheproxy.NewCache(c, geoService)
	authService := auth.NewAuthService([]byte("fj0qwef089wsjg80sjg90dj$@UJF)!@JF"), c)

	geoController := controller.NewGeoServiceController(geoService, cache)
	authControl := authController.NewAuthController(authService)
	superController := modules.NewSuperController(*geoController, *authControl)

	if err = router.Serve(superController, metrics.PrometheusMetrics); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
