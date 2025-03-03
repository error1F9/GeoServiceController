package modules

import (
	"GeoService/internal/modules/address/service"
	authService "GeoService/internal/modules/auth/service"
)

type SuperService struct {
	geoService  service.GeoService
	authService authService.AuthService
}

func NewSuperService(geoService service.GeoService, authService authService.AuthService) *SuperService {
	return &SuperService{geoService: geoService, authService: authService}
}
