package authController

import (
	"GeoService/internal/modules/auth/entity"
	"GeoService/internal/modules/auth/service"
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

type AuthController struct {
	service authService.Auther
}

func NewAuthController(service authService.Auther) *AuthController {
	return &AuthController{service}
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
func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.Register(r.Context(), user); err != nil {
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
func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := c.service.Login(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenString)
}

func (c *AuthController) JWTAuth() *jwtauth.JWTAuth {
	return c.service.GetJWTAuth()
}
