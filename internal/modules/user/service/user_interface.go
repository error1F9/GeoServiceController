package service

import (
	"GeoService/internal/models"
	"context"
)

type Userer interface {
	Create(ctx context.Context, in UserCreateIn) UserCreateOut
	Update(ctx context.Context, in UserUpdateIn) UserUpdateOut
	Delete(ctx context.Context, in UserDeleteIn) UserDeleteOut
	ChangePassword(ctx context.Context, in ChangePasswordIn) ChangePasswordOut
	GetByEmail(ctx context.Context, in GetByEmailIn) UserOut
	GetById(ctx context.Context, in GetByIDIn) UserOut
	GetByUsername(ctx context.Context, in GetByUsernameIn) UserOut
}

type UserCreateIn struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserCreateOut struct {
	UserID    int64 `json:"user_id"`
	ErrorCode int   `json:"error_code"`
}

type UserUpdateIn struct {
	User models.User `json:"user"`
}

type UserUpdateOut struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code"`
}

type ChangePasswordIn struct {
	UserID      int    `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordOut struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code"`
}

type GetByEmailIn struct {
	Email string `json:"email"`
}

type UserOut struct {
	User      models.User `json:"user"`
	ErrorCode int         `json:"error_code"`
}

type GetByIDIn struct {
	Id int `json:"id"`
}

type GetByUsernameIn struct {
	Username string `json:"username"`
}

type UserDeleteIn struct {
	UserID int `json:"user_id"`
}

type UserDeleteOut struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code"`
}
