package controller

import (
	"GeoService/internal/modules/user/service"
)

type CreateUserResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Data struct {
	Message string                `json:"message,omitempty"`
	User    service.UserCreateOut `json:"user,omitempty"`
}
