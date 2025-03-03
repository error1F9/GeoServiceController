package tests

import (
	"GeoService/internal/modules/address/entity"
	"GeoService/internal/modules/address/service"
	"bou.ke/monkey"
	"context"
	"errors"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var ErrEmptyQuery = errors.New("empty query")

func AddressSearch200Test(_ *suggest.Api, ctx context.Context, params *suggest.RequestParams) ([]*suggest.AddressSuggestion, error) {
	return []*suggest.AddressSuggestion{
		{
			Value: "Test Value 1",
			Data: &model.Address{
				City:   "Test City 1",
				Street: "Test Street 1",
				House:  "Test House 1",
				GeoLat: "Test GeoLat 1",
				GeoLon: "Test GeoLon 1",
			},
		}}, nil
}

func AddressSearchEmptyResult(_ *suggest.Api, ctx context.Context, params *suggest.RequestParams) ([]*suggest.AddressSuggestion, error) {
	return []*suggest.AddressSuggestion{}, ErrEmptyQuery
}

func TestAddressSearchService(t *testing.T) {
	tests := []struct {
		Name           string
		InputFunc      func(_ *suggest.Api, ctx context.Context, params *suggest.RequestParams) ([]*suggest.AddressSuggestion, error)
		ExpectedOutput []*entity.Address
		ExpectedError  error
		Input          string
	}{
		{
			Name:      "Test 1: no errors",
			InputFunc: AddressSearch200Test,
			ExpectedOutput: []*entity.Address{
				{
					City:   "Test City 1",
					Street: "Test Street 1",
					House:  "Test House 1",
					Lat:    "Test GeoLat 1",
					Lon:    "Test GeoLon 1",
				},
			},
			ExpectedError: nil,
			Input:         "Test Value 1",
		},
		{
			Name:           "Test 2: empty query",
			InputFunc:      AddressSearchEmptyResult,
			ExpectedOutput: nil,
			ExpectedError:  ErrEmptyQuery,
			Input:          "",
		},
	}

	var api *suggest.Api

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(api), "Address", test.InputFunc)
			defer monkey.UnpatchAll()

			geoService := service.NewGeoService("test", "test")
			addresses, err := geoService.AddressSearch(test.Input)
			assert.ErrorIs(t, test.ExpectedError, err)
			assert.Equal(t, test.ExpectedOutput, addresses)

		})
	}

}

func TestGeoCodeService(t *testing.T) {
	t.Run("GeoCodeService", func(t *testing.T) {

		geoService := service.NewGeoService("90a5dd26d0ba58ea94f25f085aa113ad67f2af27", "eb3066ce98823788c54dafb9e5e66d87a3c92d9d")

		addresses, err := geoService.GeoCode("55.878", "37.653")

		assert.NoError(t, err)
		assert.NotEmpty(t, addresses)
		for _, addr := range addresses {
			assert.NotEmpty(t, addr.City)
			assert.NotEmpty(t, addr.Lat)
			assert.NotEmpty(t, addr.Lon)
		}
	})
}
