package song_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"restApi/internal/song"
	"strconv"
)

func SongCreateHandler(w http.ResponseWriter, r *http.Request, repo song.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var newSong song.Song
	err = json.Unmarshal(body, &newSong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Println(newSong)
	err = repo.Create(ctx, &newSong)
	if err != nil {
		log.Error("Error creating song:", err)
		http.Error(w, "Error creating song", http.StatusInternalServerError)
	}
}

func SongGetAllHandler(w http.ResponseWriter, r *http.Request, repo song.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	songs, err := repo.FindAll(ctx)
	if err != nil {
		log.Error("Error getting all songs:", err)
		http.Error(w, "Error getting all songs", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(songs)
	if err != nil {
		log.Error("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func SongGetHandler(w http.ResponseWriter, r *http.Request, repo song.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Error("Error converting id to int:", err)
	}

	song, err := repo.FindOne(ctx, id)
	if err != nil {
		log.Error("Error getting song:", err)
	}

	jsonData, err := json.Marshal(song)
	if err != nil {
		log.Error("Error marshaling JSON:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func RegisterRoutes(r chi.Router, repo song.Repository, log *slog.Logger, ctx context.Context) {
	r.Post("/song", func(w http.ResponseWriter, r *http.Request) {
		SongCreateHandler(w, r, repo, log, ctx)
	})
	r.Get("/songs", func(w http.ResponseWriter, r *http.Request) {
		SongGetAllHandler(w, r, repo, log, ctx)
	})
	r.Get("/song", func(w http.ResponseWriter, r *http.Request) {
		SongGetHandler(w, r, repo, log, ctx)
	})
}
