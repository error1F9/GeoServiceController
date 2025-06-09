package authService

import (
	"context"
)

type Auther interface {
	Register(ctx context.Context, in RegisterIn) RegisterOut
	Login(ctx context.Context, login, password string) AuthorizeOut
}

type AuthorizeOut struct {
	UserId       int64  `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Error        error
}

type RegisterIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type RegisterOut struct {
	Id    *int64 `json:"id"`
	Error error
}
