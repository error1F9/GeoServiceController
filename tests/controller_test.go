package tests

import (
	"GeoService/internal/cacheproxy"
	"GeoService/internal/modules/address/controller"
	"GeoService/internal/modules/address/entity"
	"GeoService/internal/modules/address/service"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAddressGeocode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockGeoProvider(ctrl)
	mockCache := cacheproxy.NewMockCacher(ctrl)

	control := controller.NewGeoServiceController(mockService, mockCache)

	tests := []struct {
		name           string
		requestBody    controller.GeocodeRequest
		mockCacheFunc  func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			requestBody: controller.GeocodeRequest{
				Lat: "55.7558",
				Lng: "37.6176",
			},
			mockCacheFunc: func() {
				mockCache.EXPECT().GeoCodeWithCache(gomock.Any(), "55.7558", "37.6176").Return([]*entity.Address{{City: "Moscow"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"addresses\":[{\"city\":\"Moscow\",\"street\":\"\",\"house\":\"\",\"lat\":\"\",\"lon\":\"\"}]}\n",
		},
		{
			name: "Empty latitude",
			requestBody: controller.GeocodeRequest{
				Lat: "",
				Lng: "37.6176",
			},
			mockCacheFunc:  func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Empty Query\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCacheFunc()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			control.HandleAddressGeocode(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestHandleAddressSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockGeoProvider(ctrl)
	mockCache := cacheproxy.NewMockCacher(ctrl)

	control := controller.NewGeoServiceController(mockService, mockCache)

	tests := []struct {
		name           string
		requestBody    controller.SearchRequest
		mockCacheFunc  func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			requestBody: controller.SearchRequest{
				Query: "Moscow",
			},
			mockCacheFunc: func() {
				mockCache.EXPECT().SearchAddressWithCache(gomock.Any(), "Moscow").Return([]*entity.Address{{City: "Moscow"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"addresses\":[{\"city\":\"Moscow\",\"street\":\"\",\"house\":\"\",\"lat\":\"\",\"lon\":\"\"}]}\n",
		},
		{
			name: "Empty query",
			requestBody: controller.SearchRequest{
				Query: "",
			},
			mockCacheFunc:  func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Empty Query\n",
		},
		{
			name: "Cache error",
			requestBody: controller.SearchRequest{
				Query: "Moscow",
			},
			mockCacheFunc: func() {
				mockCache.EXPECT().SearchAddressWithCache(gomock.Any(), "Moscow").Return(nil, errors.New("cache error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "cache error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCacheFunc()

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/address/search", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			control.HandleAddressSearch(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
