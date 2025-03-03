package controller

import (
	"GeoService/internal/modules/address/entity"
)

type GeocodeRequest struct {
	Lat string `json:"lat" example:"55.878"`
	Lng string `json:"lng" example:"37.653"`
}

type GeocodeResponse struct {
	Addresses []*entity.Address `json:"addresses"`
}

type SearchRequest struct {
	Query string `json:"query" example:"мск сухонска 11/-89"`
}
type SearchResponse struct {
	Addresses []*entity.Address `json:"addresses"`
}
