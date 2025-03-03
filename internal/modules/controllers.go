package modules

import (
	"GeoService/internal/modules/address/controller"
	authController "GeoService/internal/modules/auth/controller"
)

type SuperController struct {
	GeoController  controller.GeoServiceController
	AuthController authController.AuthController
}

func NewSuperController(geoController controller.GeoServiceController, authController authController.AuthController) *SuperController {
	return &SuperController{
		GeoController:  geoController,
		AuthController: authController,
	}
}
