package controller

import (
	"GeoService/internal/infrastucture/components"
	"GeoService/internal/infrastucture/errors"
	"GeoService/internal/infrastucture/responder"
	"GeoService/internal/models"
	"GeoService/internal/modules/user/service"
	"encoding/json"
	"github.com/ptflp/godecoder"
	"net/http"
)

type UserController struct {
	service service.Userer
	responder.Responder
	godecoder.Decoder
}

func NewUserController(service service.Userer, components *components.Components) *UserController {
	return &UserController{service: service, Responder: components.Responder, Decoder: components.Decoder}
}

type Userer interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	GetByEmail(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	GetByUsername(w http.ResponseWriter, r *http.Request)
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	out := u.service.Create(r.Context(), service.UserCreateIn{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	})
	if out.ErrorCode != errors.NoError {
		u.OutputJSON(w, CreateUserResponse{
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: "retrieving user error",
			},
		})
	}
	u.OutputJSON(w, CreateUserResponse{
		ErrorCode: errors.NoError,
		Success:   true,
		Data: Data{
			User: out,
		},
	})
}

func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.ErrorBadRequest(w, err)
		return
	}

	out := u.service.Update(r.Context(), service.UserUpdateIn{
		User: *user,
	})

	if out.ErrorCode != errors.NoError {
		u.OutputJSON(w, UpdateUserRespone)
	}
}
