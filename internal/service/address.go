package service

import (
	"GeoService/internal"
	"GeoService/internal/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"strings"
)

var logPass = make(map[string]string)

type GeoService struct {
	jwtAuth   *jwtauth.JWTAuth
	api       *suggest.Api
	apiKey    string
	secretKey string
}

type GeoProvider interface {
	AddressSearch(input string) ([]*entity.Address, error)
	GeoCode(lat, lng string) ([]*entity.Address, error)
	Register(user entity.User) error
	Login(user entity.User) (string, error)
	GetJWTAuth() *jwtauth.JWTAuth
}

func NewGeoService(apiKey, secretKey string) *GeoService {
	var err error
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	tokenAuth := jwtauth.New("HS256", []byte("123456"), nil)

	return &GeoService{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
		jwtAuth:   tokenAuth,
	}
}

func (g *GeoService) GetJWTAuth() *jwtauth.JWTAuth {
	return g.jwtAuth
}

func (g *GeoService) AddressSearch(input string) ([]*entity.Address, error) {
	var res []*entity.Address
	rawRes, err := g.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &entity.Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}

	return res, nil
}

func (g *GeoService) GeoCode(lat, lng string) ([]*entity.Address, error) {
	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", g.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var geoCode internal.GeoCode

	err = json.NewDecoder(resp.Body).Decode(&geoCode)
	if err != nil {
		return nil, err
	}
	var res []*entity.Address
	for _, r := range geoCode.Suggestions {
		var address entity.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		res = append(res, &address)
	}

	return res, nil
}

func (g *GeoService) Register(user entity.User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New("username or password is empty")
	}

	if _, ok := logPass[user.Username]; ok {
		return errors.New("username exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	logPass[user.Username] = string(hashedPassword)

	return err
}

func (g *GeoService) Login(user entity.User) (string, error) {
	if _, ok := logPass[user.Username]; !ok {
		return "", errors.New("username not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(logPass[user.Username]), []byte(user.Password)); err != nil {
		return "", errors.New("password is wrong")
	}

	_, tokenString, err := g.jwtAuth.Encode(map[string]interface{}{"username": user.Username})

	return tokenString, err

}
