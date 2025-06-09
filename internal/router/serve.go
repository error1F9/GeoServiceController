package router

import (
	"GeoService/internal/infrastucture/components"
	"GeoService/internal/infrastucture/middleware"
	"GeoService/internal/modules"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
)

func Serve(controllers *modules.Controllers, components *components.Components) *chi.Mux {
	r := chi.NewRouter()
	r.Use(components.Proxy.ReverseProxy)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})
		authCheck := middleware.NewTokenManager(components.Responder, components.TokenManager)
		//r.With(metrics.PrometheusMetrics.).Post("/register", c.AuthController.HandleRegister)
		//r.With(metrics.DurationAndCounterMW).Post("/login", c.AuthController.HandleLogin)
		r.Route("/auth", func(r chi.Router) {
			authController := controllers.AuthController
			r.Post("/register", authController.HandleRegister)
			r.Post("/login", authController.HandleLogin)
			r.Route("/refresh", func(r chi.Router) {
				r.Use(authCheck.CheckRefresh)
				//r.Post("/", authController.Refresh)
			})
		})

		r.Route("/user", func(r chi.Router) {
			//r.Use(authCheck.CheckStrict)
			userController := controllers.UserController
			r.Post("/", userController.CreateUser)
			//r.Put("/", userController.Up)
		})

		//r.Group(func(r chi.Router) {
		//	r.Mount("/debug/pprof", pprofHandler())
		//	r.Post("/address/search", c.GeoController.HandleAddressSearch)
		//	r.Post("/address/geocode", c.GeoController.HandleAddressGeocode)
		//
		//})
	})

	r.Mount("/metrics", promhttp.Handler())
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	return r
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
