package middleware

import (
	"GeoService/internal/infrastucture/responder"
	"GeoService/internal/infrastucture/tools/cryptography"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const authorization = "Authorization"

type Token struct {
	responder.Responder
	jwt cryptography.TokenManager
}

type UserRequest struct{}

func NewTokenManager(responder responder.Responder, jwt cryptography.TokenManager) *Token {
	return &Token{
		Responder: responder,
		jwt:       jwt,
	}
}

func (t *Token) CheckStrict(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRaw := r.Header.Get(authorization)
		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			t.ErrorUnauthorized(w, fmt.Errorf("wrong input data"))
			return
		}
		u, err := t.jwt.ParseToken(tokenParts[1], cryptography.AccessToken)
		if err != nil && err.Error() == "Token is expired" {
			t.ErrorUnauthorized(w, errors.New("token expired"))
			return
		}
		if err != nil {
			t.ErrorUnauthorized(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), UserRequest{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t *Token) CheckRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRaw := r.Header.Get(authorization)
		tokenParts := strings.Split(tokenRaw, " ")
		if len(tokenParts) < 2 && tokenParts[0] != "Bearer" {
			t.ErrorUnauthorized(w, fmt.Errorf("wrong input data"))
			return
		}
		u, err := t.jwt.ParseToken(tokenParts[1], cryptography.RefreshToken)
		if err != nil && err.Error() == "Token expired" {
			t.ErrorUnauthorized(w, err)
			return
		}
		if err != nil {
			t.ErrorUnauthorized(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), UserRequest{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t *Token) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(authorization)
		u, err := t.jwt.ParseToken(token, cryptography.AccessToken)
		if err != nil {
			u = cryptography.UserClaims{}
		}
		ctx := context.WithValue(r.Context(), UserRequest{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
