package router

import (
	"GeoService/internal/controller"
	"GeoService/pkg/reverse"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(c *controller.GeoServiceController) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(reverse.NewReverseProxy("hugo", "1313").ReverseProxy)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})

		r.Post("/register", c.HandleRegister)
		r.Post("/login", c.HandleLogin)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(c.JWTAuth()))
			r.Use(jwtauth.Authenticator(c.JWTAuth()))

			r.Post("/address/search", c.HandleAddressSearch)
			r.Post("/address/geocode", c.HandleAddressGeocode)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
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

	log.Println("Waiting 5 seconds before shutdown...")
	time.Sleep(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
