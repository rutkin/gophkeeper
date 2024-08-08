package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/adapter/config"
	httpserver "github.com/rutkin/gophkeeper/internal/server/adapter/http_server"
	repositry "github.com/rutkin/gophkeeper/internal/server/adapter/repository/file"
	"github.com/rutkin/gophkeeper/internal/server/adapter/repository/postgress"
	"github.com/rutkin/gophkeeper/internal/server/adapter/token"
	"github.com/rutkin/gophkeeper/internal/server/core/service"
)

func initLogger(cfg config.Config) {
	switch cfg.LogLevel {
	case config.LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case config.LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func initService(cfg config.Config) {
	userRepository, err := postgress.NewUserRepo(cfg.DatabaseDSN)
	if err != nil {
		log.Err(err).Msg("filed to create user repository")
		os.Exit(1)
	}
	defer userRepository.Close()
	keeperRepository, err := repositry.NewKeeper()
	if err != nil {
		log.Err(err).Msg("filed to create keeper repository")
		os.Exit(1)
	}
	tokenService := token.New(time.Hour * time.Duration(cfg.TokenExpiration))
	authService := service.NewAuthService(userRepository, tokenService)
	keeperService := service.NewKeeperService(keeperRepository)
	handler := httpserver.NewHandler(authService, keeperService, tokenService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("failed to shutdown service")
		}
		close(idleConnsClosed)
	}()
	err = srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Err(err).Msg("failed to start keeper service")
		os.Exit(1)
	}
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Err(err).Msg("failed to create config")
		os.Exit(1)
	}

	initLogger(cfg)
	log.Info().Msg("Starting keeper service")
	initService(cfg)
	log.Info().Msg("keeper service stopped")
}
