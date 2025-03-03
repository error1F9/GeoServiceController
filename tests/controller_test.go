package tests

//
//import (
//	"bytes"
//	"errors"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestHandleAddressGeocode(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockService := service.NewMockGeoProvider(ctrl)
//	mockCache := cacheproxy.NewMockCacher(ctrl)
//
//	controller := NewGeoServiceController(mockService, mockCache)
//
//	tests := []struct {
//		name           string
//		requestBody    GeocodeRequest
//		mockCacheFunc  func()
//		expectedStatus int
//		expectedBody   string
//	}{
//		{
//			name: "Valid request",
//			requestBody: GeocodeRequest{
//				Lat: "55.7558",
//				Lng: "37.6176",
//			},
//			mockCacheFunc: func() {
//				mockCache.EXPECT().GeoCodeWithCache(gomock.Any(), "55.7558", "37.6176").Return("Moscow", nil)
//			},
//			expectedStatus: http.StatusOK,
//			expectedBody:   `{"address":"Moscow"}`,
//		},
//		{
//			name: "Empty latitude",
//			requestBody: GeocodeRequest{
//				Lat: "",
//				Lng: "37.6176",
//			},
//			mockCacheFunc:  func() {},
//			expectedStatus: http.StatusBadRequest,
//			expectedBody:   "Empty Query\n",
//		},
//		{
//			name: "Cache error",
//			requestBody: GeocodeRequest{
//				Lat: "55.7558",
//				Lng: "37.6176",
//			},
//			mockCacheFunc: func() {
//				mockCache.EXPECT().GeoCodeWithCache(gomock.Any(), "55.7558", "37.6176").Return("", errors.New("cache error"))
//			},
//			expectedStatus: http.StatusInternalServerError,
//			expectedBody:   "cache error\n",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mockCacheFunc()
//
//			body, _ := json.Marshal(tt.requestBody)
//			req, _ := http.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer(body))
//			req.Header.Set("Content-Type", "application/json")
//
//			rr := httptest.NewRecorder()
//			controller.HandleAddressGeocode(rr, req)
//
//			assert.Equal(t, tt.expectedStatus, rr.Code)
//			assert.Equal(t, tt.expectedBody, rr.Body.String())
//		})
//	}
//}
//Тесты для HandleAddressSearch
//go
//Copy
//func TestHandleAddressSearch(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockService := service.NewMockGeoProvider(ctrl)
//	mockCache := cacheproxy.NewMockCacher(ctrl)
//
//	controller := NewGeoServiceController(mockService, mockCache)
//
//	tests := []struct {
//		name           string
//		requestBody    SearchRequest
//		mockCacheFunc  func()
//		expectedStatus int
//		expectedBody   string
//	}{
//		{
//			name: "Valid request",
//			requestBody: SearchRequest{
//				Query: "Moscow",
//			},
//			mockCacheFunc: func() {
//				mockCache.EXPECT().SearchAddressWithCache(gomock.Any(), "Moscow").Return([]string{"Moscow, Russia"}, nil)
//			},
//			expectedStatus: http.StatusOK,
//			expectedBody:   `{"addresses":["Moscow, Russia"]}`,
//		},
//		{
//			name: "Empty query",
//			requestBody: SearchRequest{
//				Query: "",
//			},
//			mockCacheFunc:  func() {},
//			expectedStatus: http.StatusBadRequest,
//			expectedBody:   "Empty Query\n",
//		},
//		{
//			name: "Cache error",
//			requestBody: SearchRequest{
//				Query: "Moscow",
//			},
//			mockCacheFunc: func() {
//				mockCache.EXPECT().SearchAddressWithCache(gomock.Any(), "Moscow").Return(nil, errors.New("cache error"))
//			},
//			expectedStatus: http.StatusInternalServerError,
//			expectedBody:   "cache error\n",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mockCacheFunc()
//
//			body, _ := json.Marshal(tt.requestBody)
//			req, _ := http.NewRequest("POST", "/api/address/search", bytes.NewBuffer(body))
//			req.Header.Set("Content-Type", "application/json")
//
//			rr := httptest.NewRecorder()
//			controller.HandleAddressSearch(rr, req)
//
//			assert.Equal(t, tt.expectedStatus, rr.Code)
//			assert.Equal(t, tt.expectedBody, rr.Body.String())
//		})
//	}
//}
