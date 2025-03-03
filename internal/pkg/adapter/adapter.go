package adapter

import (
	"GeoService/internal/modules/address/entity"
	"errors"
)

func GeoCodeAdapter(geoCodeFunc func(lat, lng string) ([]*entity.Address, error)) func(params ...string) ([]*entity.Address, error) {
	return func(params ...string) ([]*entity.Address, error) {
		if len(params) != 2 {
			return nil, errors.New("expected exactly 2 parameters: lat and lng")
		}
		return geoCodeFunc(params[0], params[1])
	}
}

func GeoAddressAdapter(geoAddress func(address string) ([]*entity.Address, error)) func(params ...string) ([]*entity.Address, error) {
	return func(params ...string) ([]*entity.Address, error) {
		if len(params) != 1 {
			return nil, errors.New("expected exactly 1 parameters: address")
		}
		return geoAddress(params[0])
	}
}
