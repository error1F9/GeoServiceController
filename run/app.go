package run

import (
	"GeoService/config"
	"GeoService/internal/db"
	"GeoService/internal/infrastucture/components"
	"GeoService/internal/infrastucture/errors"
	"GeoService/internal/infrastucture/responder"
	"GeoService/internal/infrastucture/reverse"
	"GeoService/internal/infrastucture/server"
	"GeoService/internal/infrastucture/tools/cryptography"
	"GeoService/internal/modules"
	"GeoService/internal/router"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

type Application interface {
	Runner
	Bootstraper
}

// Runner - интерфейс запуска приложения
type Runner interface {
	Run() int
}

// Bootstraper - интерфейс инициализации приложения
type Bootstraper interface {
	Bootstrap(options ...interface{}) Runner
}

// App - структура приложения
type App struct {
	conf     config.AppConfig
	logger   *zap.Logger
	srv      server.Server
	Sig      chan os.Signal
	Repos    *modules.Repositories
	Services *modules.Services
}

// NewApp - конструктор приложения
func NewApp(conf config.AppConfig, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

// Run - запуск приложения
func (a *App) Run() int {
	// на русском
	// создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	// запускаем горутину для graceful shutdown
	// при получении сигнала SIGINT
	// вызываем cancel для контекста
	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return errors.GeneralError
	}

	return errors.NoError
}

func (a *App) Bootstrap(options ...interface{}) Runner {
	tokenManager := cryptography.NewTokenJWT(a.conf.JWTKey)

	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})

	responseManager := responder.NewResponder(decoder, a.logger)

	dbx := db.NewBunDB(a.conf)

	newRepositories := modules.NewRepositories(dbx)
	a.Repos = newRepositories

	reverseProxy := reverse.NewReverseProxy("localhost", a.conf.Server.Port)

	components := components.NewComponents(a.conf, tokenManager, responseManager, decoder, a.logger, reverseProxy)

	services := modules.NewServices(newRepositories, components)
	a.Services = services

	controllers := modules.NewControllers(services, components)

	r := router.Serve(controllers, components)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}

	a.srv = server.NewHttpServer(a.conf.Server, srv, a.logger)

	return a
}
