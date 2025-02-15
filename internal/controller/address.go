package controller

import (
	"GeoService/internal/entity"
	"GeoService/internal/service"
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

type GeoServiceController struct {
	service service.GeoService
}

func NewGeoServiceController(service service.GeoService) *GeoServiceController {
	return &GeoServiceController{service}
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
func (g *GeoServiceController) HandleAddressGeocode(w http.ResponseWriter, r *http.Request) {
	var geoReq GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&geoReq); err != nil || geoReq.Lng == "" || geoReq.Lat == "" {
		http.Error(w, "Empty Query", http.StatusBadRequest)
		return
	}

	geo, err := g.service.GeoCode(geoReq.Lat, geoReq.Lng)
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

	addresses, err := g.service.AddressSearch(req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SearchResponse{addresses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

// HandleRegister
// @Summary Register a user
// @Tags Login and Registration
// @Description Register a user with login and password
// @Accept json
// @Produce json
// @Param input body entity.User true "Username and Pass for registration"
// @Success 200 {string} string "User created"
// @failure 400 {string} string "Empty Query"
// @failure 500 {string} string "Internal Server Error"
// @Router /api/register [post]
func (g *GeoServiceController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := g.service.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created")
}

// HandleLogin
// @Summary Login with username and password
// @Tags Login and Registration
// @Description Login with username and password
// @Accept json
// @Produce json
// @Param input body entity.User true "Username and Pass for logining"
// @Success 200 {string} string "token string"
// @failure 400 {string} string "Empty Query"
// @failure 500 {string} string "Internal Server Error"
// @Router /api/login [post]
func (g *GeoServiceController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := g.service.Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenString)
}

func (g *GeoServiceController) JWTAuth() *jwtauth.JWTAuth {
	return g.service.GetJWTAuth()
}
