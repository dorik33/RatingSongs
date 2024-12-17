package main

//TODO: init router: chi, render
//TODO: run server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"restApi/internal/config"
	"restApi/internal/http-server/handlers/rating_handlers"
	"restApi/internal/http-server/handlers/song_handlers"
	mwLogger "restApi/internal/http-server/middleware/logger"
	ratingdb "restApi/internal/rating/db"
	songdb "restApi/internal/song/db"
	"restApi/pkg/client/postgresql"
	"restApi/pkg/logger"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)

	log := logger.SetupLogger(cfg.Env)
	log.Info("Logger is init", slog.String("env", cfg.Env))
	log.Debug("Debug is enabled")

	pool, err := postgresql.NewClient(context.Background(), cfg.Storage)
	if err != nil {
		log.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	defer pool.Close()
	log.Debug("Connected to postgresql")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	songRep := songdb.NewRepository(pool)
	ratingRep := ratingdb.NewRepository(pool)
	song_handlers.RegisterRoutes(router, songRep, log, context.Background())
	rating_handlers.RegisterRoutes(router, ratingRep, log, context.Background())

	if err = RunServer(cfg, router, log); err != nil {
		log.Error(fmt.Sprintf("%+v", err))
	}

}

func RunServer(cfg *config.Config, router *chi.Mux, log *slog.Logger) error {
	server := &http.Server{
		Addr:         cfg.HTTPServer.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info("Server starting", "address", cfg.HTTPServer.Addres)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", err)
		return err
	}

	log.Info("Server exiting")
	return nil
}
