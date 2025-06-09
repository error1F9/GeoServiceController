package authService

import (
	"GeoService/config"
	"GeoService/internal/infrastucture/tools/cryptography"
	"GeoService/internal/models"
	"GeoService/internal/modules/user/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"strconv"
)

type AuthService struct {
	config  config.AppConfig
	db      repository.Userer
	jwtAuth cryptography.TokenManager
	logger  *zap.Logger
}

func NewAuthService(db repository.Userer, jwtauth cryptography.TokenManager) *AuthService {
	return &AuthService{db: db, jwtAuth: jwtauth}
}

func (a *AuthService) Register(ctx context.Context, in RegisterIn) RegisterOut {
	if in.Username == "" || in.Password == "" {
		return RegisterOut{
			Id:    nil,
			Error: errors.New("username or password is empty"),
		}
	}

	hashedPassword, err := cryptography.HashPassword(in.Password)
	if err != nil {
		return RegisterOut{
			nil,
			errors.New("failed to hash password"),
		}
	}
	in.Password = hashedPassword

	user := &models.User{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
		Role:     in.Role,
	}
	id, err := a.db.CreateUser(ctx, user)

	return RegisterOut{
		Id:    &id,
		Error: err,
	}
}

func (a *AuthService) Login(ctx context.Context, login, password string) AuthorizeOut {
	user, err := a.db.GetUserByUsername(ctx, login)
	if err != nil {
		return AuthorizeOut{
			Error: err,
		}
	}

	if !cryptography.CheckPassword(user.Password, password) {
		return AuthorizeOut{
			Error: err,
		}
	}

	accessToken, refreshToken, err := a.generateTokens(user)
	if err != nil {
		return AuthorizeOut{
			Error: err,
		}
	}

	return AuthorizeOut{
		UserId:       user.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

}

func (a *AuthService) generateTokens(user *models.User) (string, string, error) {
	accessToken, err := a.jwtAuth.CreateToken(
		strconv.Itoa(int(user.Id)),
		user.Role,
		"",
		a.config.JWTKey.AccessTTL,
		cryptography.AccessToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", errors.New("failed to generate access token")
	}
	refreshToken, err := a.jwtAuth.CreateToken(
		strconv.Itoa(int(user.Id)),
		user.Role,
		"",
		a.config.JWTKey.RefreshTTL,
		cryptography.RefreshToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}
