package components

import (
	"GeoService/config"
	"GeoService/internal/infrastucture/responder"
	"GeoService/internal/infrastucture/reverse"
	"GeoService/internal/infrastucture/tools/cryptography"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Config       config.AppConfig
	TokenManager cryptography.TokenManager
	Responder    responder.Responder
	Decoder      godecoder.Decoder
	Logger       *zap.Logger
	Proxy        reverse.ReverseProxier
}

func NewComponents(conf config.AppConfig, tokenManager cryptography.TokenManager, responder responder.Responder, decoder godecoder.Decoder, logger *zap.Logger, proxy reverse.ReverseProxier) *Components {
	return &Components{Config: conf, TokenManager: tokenManager, Responder: responder, Decoder: decoder, Logger: logger, Proxy: proxy}
}
