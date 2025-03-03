package router

import (
	"GeoService/internal/modules"
	"GeoService/internal/pkg/metrics"
	"GeoService/internal/reverse"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(c *modules.SuperController, metrics *metrics.Metrics) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(reverse.NewReverseProxy("hugo", "1313").ReverseProxy)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})

		r.With(metrics.DurationAndCounterMW).Post("/register", c.AuthController.HandleRegister)
		r.With(metrics.DurationAndCounterMW).Post("/login", c.AuthController.HandleLogin)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(c.AuthController.JWTAuth()))
			r.Use(jwtauth.Authenticator(c.AuthController.JWTAuth()))

			r.Mount("/debug/pprof", pprofHandler())
			r.With(metrics.DurationAndCounterMW).Post("/address/search", c.GeoController.HandleAddressSearch)
			r.With(metrics.DurationAndCounterMW).Post("/address/geocode", c.GeoController.HandleAddressGeocode)

		})
	})

	r.Mount("/metrics", promhttp.Handler())
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

func pprofHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", pprof.Index)
	r.Get("/cmdline", pprof.Cmdline)
	r.Get("/profile", pprof.Profile)
	r.Get("/symbol", pprof.Symbol)
	r.Get("/trace", pprof.Trace)
	r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.Get("/heap", pprof.Handler("heap").ServeHTTP)
	r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	r.Get("/block", pprof.Handler("block").ServeHTTP)
	r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
	return r
}
