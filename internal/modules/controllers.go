package modules

import (
	"GeoService/internal/infrastucture/components"
	geoController "GeoService/internal/modules/address/controller"
	authController "GeoService/internal/modules/auth/controller"
	userController "GeoService/internal/modules/user/controller"
)

type Controllers struct {
	GeoController  geoController.Geoer
	AuthController authController.Auther
	UserController userController.Userer
}

func NewControllers(services *Services, components *components.Components) *Controllers {
	return &Controllers{
		GeoController:  geoController.NewGeoController(services.GeoService),
		AuthController: authController.NewAuthController(services.AuthService, components.Responder),
		UserController: userController.NewUserController(services.UserService, components),
	}
}
