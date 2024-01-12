package main

import (
	"context"
	"github.com/degeboman/gas/constant"
	"github.com/degeboman/gas/internal/config"
	"github.com/degeboman/gas/internal/http-server/handlers/auth/refresh"
	"github.com/degeboman/gas/internal/http-server/handlers/auth/signin"
	"github.com/degeboman/gas/internal/http-server/handlers/auth/signup"
	"github.com/degeboman/gas/internal/http-server/handlers/auth/verify"
	mwLogger "github.com/degeboman/gas/internal/http-server/middleware/logger"
	"github.com/degeboman/gas/internal/lib/logger"
	"github.com/degeboman/gas/internal/lib/logger/sl"
	"github.com/degeboman/gas/internal/storage/mongodb"
	"github.com/degeboman/gas/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	@title			Swagger GAS API
//	@version		1.0
//	@description	GAS this is a general authorization service. Gas is developed for personal use, in particular for general use on pet projects

//	@host			localhost:2023
//	@BasePath		/api/v1

//	@securityDefinitions.basic JWT Auth

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"gas is running",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)

	storage := mongodb.New(cfg.MongoConnectionString)

	log.Info(
		"storage is running",
	)

	u := usecase.New(&storage)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route(constant.AuthRoute, func(r chi.Router) {
		r.Post(constant.SignUpRoute, signup.New(log, u))
		r.Post(constant.SignInRoute, signin.New(log, cfg, u))
		r.Post(constant.VerifyRoute, verify.New(log, cfg, u))
		r.Post(constant.RefreshRoute, refresh.New(log, cfg, u))
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}
