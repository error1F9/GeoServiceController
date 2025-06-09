package authController

import (
	"GeoService/internal/infrastucture/responder"
	"GeoService/internal/models"
	"GeoService/internal/modules/auth/service"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	service authService.Auther
	responder.Responder
}

func NewAuthController(service authService.Auther, Responder responder.Responder) *AuthController {
	return &AuthController{service: service, Responder: Responder}
}

type Auther interface {
	HandleRegister(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
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
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	in := authService.RegisterIn{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Role:     user.Role,
	}

	out := c.service.Register(r.Context(), in)
	if out.Error != nil {
		http.Error(w, out.Error.Error(), http.StatusBadRequest)
		return
	}

	c.OutputJSON(w, map[string]interface{}{"id": out.Id})
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
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := c.service.Login(r.Context(), user.Username, user.Password)
	if out.Error != nil {
		http.Error(w, out.Error.Error(), http.StatusBadRequest)
	}

	c.OutputJSON(w, map[string]interface{}{
		"User_id":       out.UserId,
		"Access_token":  out.AccessToken,
		"Refresh_token": out.RefreshToken,
	})
}
