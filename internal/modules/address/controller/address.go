package controller

import (
	"GeoService/internal/cacheproxy"
	"GeoService/internal/modules/address/service"
	"encoding/json"
	"net/http"
)

type GeoServiceController struct {
	service service.GeoProvider
	cache   cacheproxy.Cacher
}

func NewGeoServiceController(service service.GeoProvider, cache cacheproxy.Cacher) *GeoServiceController {
	return &GeoServiceController{service: service, cache: cache}
}

type GeoServiceControllerInterface interface {
	HandleAddressGeocode(w http.ResponseWriter, r *http.Request)
	HandleAddressSearch(w http.ResponseWriter, r *http.Request)
}

// HandleAddressGeocode
// @Summary receive Address by GeoCode
// @Tags GeoCode
// @Security ApiKeyAuth
// @Description Request structure for geocoding addresses
// @ID geo
// @Accept json
// @Produce json
// @Param input body GeocodeRequest true "Handle Address by GeoCode"
// @Success 200 {object} GeocodeResponse
// @failure 400 {string} string "Empty Query"
// @failure 500 {string} string "Internal Server Error"
// @Router /api/address/geocode [post]

//go:generate mockgen -source=address.go -destination=mocks/mock_controller.go -package=mocks
func (g *GeoServiceController) HandleAddressGeocode(w http.ResponseWriter, r *http.Request) {
	var geoReq GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&geoReq); err != nil || geoReq.Lng == "" || geoReq.Lat == "" {
		http.Error(w, "Empty Query", http.StatusBadRequest)
		return
	}

	geo, err := g.cache.GeoCodeWithCache(r.Context(), geoReq.Lat, geoReq.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GeocodeResponse{geo}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HandleAddressSearch
// @Summary receive Information by Address
// @Tags AddressSearch
// @Security ApiKeyAuth
// @Description Receive Information by Address
// @ID addSearch
// @Accept json
// @Produce json
// @Param input body SearchRequest true "Receive information by Address"
// @Success 200 {object} SearchResponse
// @failure 400 {string} string "Empty Query"
// @failure 500 {string} string "Internal Server Error"
// @Router /api/address/search [post]
func (g *GeoServiceController) HandleAddressSearch(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Query == "" {
		http.Error(w, "Empty Query", http.StatusBadRequest)
		return
	}

	addresses, err := g.cache.SearchAddressWithCache(r.Context(), req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SearchResponse{addresses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
