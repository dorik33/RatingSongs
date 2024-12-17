package main

//TODO: init config: cleanenv
//TODO: init logger: slog
//TODO: init storage: psql
//TODO: init router: chi, render
//TODO: run server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"restApi/internal/config"
	"restApi/internal/song"
	"restApi/pkg/client/postgresql"
	"restApi/pkg/logger"
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

	newSong := song.Song{
		Title:       "Another One Bites the Dust",
		Artist:      "Queen",
		Album:       "The Game",
		ReleaseDate: time.Date(1980, 10, 10, 0, 0, 0, 0, time.UTC),
	}

	err = InsertSong(pool, newSong)
	if err != nil {
		log.Error(fmt.Sprintf("%+v", err))
	}

	fmt.Println("Song inserted successfully!")
}

func InsertSong(pool *pgxpool.Pool, s song.Song) error {
	query := `INSERT INTO song (title, artist, album, release_date) VALUES ($1, $2, $3, $4)`
	_, err := pool.Exec(context.Background(), query, s.Title, s.Artist, s.Album, s.ReleaseDate)
	return err
}
