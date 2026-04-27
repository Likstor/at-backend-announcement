package main

import (
	"at-backend-announcement/internal/handlers"
	"at-backend-announcement/internal/handlers/middleware"
	"at-backend-announcement/internal/infrastructure/db/postgres"
	"at-backend-announcement/internal/pkg/config"
	"at-backend-announcement/internal/pkg/health"
	"at-backend-announcement/internal/pkg/logs"
	"at-backend-announcement/internal/pkg/pgclient"
	"at-backend-announcement/internal/pkg/tokens"
	usecase "at-backend-announcement/internal/usecases"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	handler := logs.NewHandlerMiddleware(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(slog.New(handler))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := Config{}
	if err := config.Load(&cfg); err != nil {
		slog.Error(err.Error())
		return
	}

	pool, err := pgclient.NewClient(ctx, 5, pgclient.Config{
		Host:     cfg.StorageHost,
		Port:     cfg.StoragePort,
		Database: cfg.StorageName,
		Username: cfg.StorageUser,
		Password: cfg.StoragePassword,
		SSLMode:  cfg.StorageSSLMode,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer pool.Close()

	tokensVerifier := tokens.NewTokensVerifier(
		15*time.Minute,
		cfg.JWKSURI,
		cfg.TokensIssuer,
		cfg.TokensKeyID,
		5*time.Second,
	)

	announcementRepo := postgres.NewAnnouncementRepository(pool)

	announcementUsecase := usecase.NewAnnouncementUsecase(announcementRepo, 100)

	announcementHandler := handlers.NewAnnouncementHandler(announcementUsecase)

	announcementHandlerForAdmins := handlers.NewAnnouncementHandlerForAdmins(announcementUsecase)

	mainMux := http.NewServeMux()

	announcementHandler.Setup("/announcements", tokensVerifier.VerifyJWT, mainMux)
	announcementHandlerForAdmins.Setup("/admins/announcements", tokensVerifier.VerifyJWT, mainMux)

	mainMuxWrapped := middleware.Cors(mainMux, cfg.AllowOriginCors)
	mainMuxWrapped = middleware.Logger(mainMux)
	mainMuxWrapped = middleware.Correlation(mainMuxWrapped)

	finMux := http.NewServeMux()
	finMux.Handle("/v1/", http.StripPrefix("/v1", mainMuxWrapped))
	health.Setup(finMux)

	server := &http.Server{
		Addr:              ":" + cfg.BackendPort,
		Handler:           finMux,
		ReadTimeout:       cfg.BackendReadTimeout,
		ReadHeaderTimeout: cfg.BackendReadHeaderTimeout,
		WriteTimeout:      cfg.BackendWriteTimeout,
		IdleTimeout:       cfg.BackendIdleTimeout,
	}

	run(context.Background(), server)

	slog.Info("app closed")
}

func run(ctx context.Context, server *http.Server) {
	var wg sync.WaitGroup
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	wg.Add(1)
	go func() {
		slog.Info("HTTP server started")
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			slog.Info("http server closed")
			wg.Done()
			return
		}

		slog.Error("http server closed with error: " + err.Error())
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()

		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("http server shutdown error: " + err.Error())
		}

		wg.Done()
	}()

	wg.Wait()
}
