package rating_handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log/slog"
	"net/http"
	"restApi/internal/rating"
	"strconv"
)

func RatingCreateHandler(w http.ResponseWriter, r *http.Request, rep rating.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Error getting all songs:", err)
	}
	var newRating rating.Rating

	err = json.Unmarshal(data, &newRating)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Error getting all songs:", err)
	}

	err = rep.Create(ctx, &newRating)
	if err != nil {
		log.Error("Error creating songs:", err)
		http.Error(w, "Error creating rating", http.StatusInternalServerError)
	}
}

func RatingGetAllHandler(w http.ResponseWriter, r *http.Request, rep rating.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ratings, err := rep.FindAll(ctx)
	if err != nil {
		log.Error("Error getting all songs:", err)
		http.Error(w, "Error getting all songs", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(ratings)
	if err != nil {
		log.Error("Error getting all songs:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func RatingGetAllBySongIDHandler(w http.ResponseWriter, r *http.Request, rep rating.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Bad id:", err)
		return
	}
	ratings, err := rep.FindAllBySongID(ctx, id)

	data, err := json.Marshal(ratings)
	if err != nil {
		log.Error("Error getting all songs:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func RatingGetHandler(w http.ResponseWriter, r *http.Request, rep rating.Repository, log *slog.Logger, ctx context.Context) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Bad id:", err)
		return
	}

	newRating, err := rep.FindOne(ctx, id)
	if err != nil {
		log.Error("Error getting song:", err)
		http.Error(w, "Error getting songs", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newRating)
	if err != nil {
		log.Error("Error json Marshal:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func RegisterRoutes(r chi.Router, rep rating.Repository, log *slog.Logger, ctx context.Context) {
	r.Post("/rating", func(w http.ResponseWriter, r *http.Request) {
		RatingCreateHandler(w, r, rep, log, ctx)
	})
	r.Get("/ratingbysong", func(w http.ResponseWriter, r *http.Request) {
		RatingGetAllBySongIDHandler(w, r, rep, log, ctx)
	})
	r.Get("/ratings", func(w http.ResponseWriter, r *http.Request) {
		RatingGetAllHandler(w, r, rep, log, ctx)
	})
	r.Get("/rating", func(w http.ResponseWriter, r *http.Request) {
		RatingGetHandler(w, r, rep, log, ctx)
	})
}
