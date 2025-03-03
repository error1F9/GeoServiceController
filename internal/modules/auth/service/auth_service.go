package authService

import (
	"GeoService/internal/infrastructure/cache"
	"GeoService/internal/modules/auth/entity"
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtAuth *jwtauth.JWTAuth
	redis   cache.Cacher
}

type Auther interface {
	Register(ctx context.Context, user entity.User) error
	Login(ctx context.Context, user entity.User) (string, error)
	GetJWTAuth() *jwtauth.JWTAuth
}

func NewAuthService(authKey []byte, redis cache.Cacher) *AuthService {
	tokenAuth := jwtauth.New("HS256", authKey, nil)
	return &AuthService{jwtAuth: tokenAuth, redis: redis}
}

func (a *AuthService) GetJWTAuth() *jwtauth.JWTAuth {
	return a.jwtAuth
}

func (a *AuthService) Register(ctx context.Context, user entity.User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New("username or password is empty")
	}

	if _, err := a.redis.Get(ctx, user.Username); !errors.Is(err, redis.Nil) {
		return errors.New("username exist")
	} else if errors.Is(err, redis.Nil) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}

		err = a.redis.Set(ctx, user.Username, string(hashedPassword), 0)
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Login(ctx context.Context, user entity.User) (string, error) {
	val, err := a.redis.Get(ctx, user.Username)
	if errors.Is(err, redis.Nil) {
		return "", errors.New("username doesn't exist")
	} else if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(val.([]byte), []byte(user.Password)); err != nil {
		return "", errors.New("password is wrong")
	}

	_, tokenString, err := a.jwtAuth.Encode(map[string]interface{}{"username": user.Username})

	return tokenString, err

}
