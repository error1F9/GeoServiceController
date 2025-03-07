package cacheproxy

import (
	"GeoService/internal/infrastructure/cache"
	"GeoService/internal/modules/address/entity"
	"GeoService/internal/modules/address/service"
	"GeoService/internal/pkg/adapter"
	"GeoService/internal/pkg/hash"
	"GeoService/internal/pkg/metrics"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type CacheProxy struct {
	cache   cache.Cacher
	service service.GeoProvider
}

func NewCache(cache cache.Cacher, service service.GeoProvider) *CacheProxy {
	return &CacheProxy{
		cache:   cache,
		service: service,
	}
}

type Cacher interface {
	SearchAddressWithCache(ctx context.Context, query string) ([]*entity.Address, error)
	GeoCodeWithCache(ctx context.Context, lat, lng string) ([]*entity.Address, error)
}

//go:generate mockgen -source=cacheproxy.go -destination=mocks/mock_controller.go -package=mocks
func (c *CacheProxy) Set(ctx context.Context, key string, data interface{}) error {
	addresses, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.cache.Set(ctx, key, addresses, time.Minute)
}

func (c *CacheProxy) GetData(ctx context.Context, function func(params ...string) ([]*entity.Address, error), params ...string) ([]*entity.Address, error, bool) {
	addresses := make([]*entity.Address, 0)
	hashed := hash.Body([]byte(strings.Join(params, "")))

	val, err := c.cache.Get(ctx, hashed)
	if errors.Is(err, redis.Nil) {
		addresses, err = function(params...)

		if err != nil {
			return nil, err, false
		}
		if err = c.Set(ctx, hashed, addresses); err != nil {
			return nil, err, false
		}
		return addresses, nil, false
	} else if err != nil {
		return nil, err, false
	}
	err = json.Unmarshal(val.([]byte), &addresses)
	if err != nil {
		return nil, err, false
	}

	return addresses, nil, true
}

func (c *CacheProxy) SearchAddressWithCache(ctx context.Context, query string) ([]*entity.Address, error) {
	wrapped := metrics.PrometheusMetrics.MethodRequestDuration("Search address", func(ctx context.Context, params ...string) ([]*entity.Address, error, bool) {
		return c.GetData(ctx, adapter.GeoAddressAdapter(c.service.AddressSearch), query)
	})
	return wrapped(ctx, query)
}

func (c *CacheProxy) GeoCodeWithCache(ctx context.Context, lat, lng string) ([]*entity.Address, error) {
	wrapped := metrics.PrometheusMetrics.MethodRequestDuration("GeoCode", func(ctx context.Context, params ...string) ([]*entity.Address, error, bool) {
		return c.GetData(ctx, adapter.GeoCodeAdapter(c.service.GeoCode), lat, lng)
	})
	return wrapped(ctx, lat, lng)
}
