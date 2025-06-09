package modules

import (
	"GeoService/internal/infrastucture/components"
	geoService "GeoService/internal/modules/address/service"
	authService "GeoService/internal/modules/auth/service"
	userService "GeoService/internal/modules/user/service"
)

type Services struct {
	GeoService  geoService.GeoProvider
	AuthService authService.Auther
	UserService userService.Userer
}

func NewServices(repositories *Repositories, components *components.Components) *Services {
	return &Services{
		GeoService:  geoService.NewGeoService(components.Config.DadataCreds.ApiKey, components.Config.DadataCreds.SecretKey),
		AuthService: authService.NewAuthService(repositories.UserRepository, components.TokenManager),
		UserService: userService.NewUserService(repositories.UserRepository, components.Logger),
	}
}
