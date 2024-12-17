package ratingdb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"restApi/internal/rating"
	"restApi/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, rating *rating.Rating) error {
	q := `
	INSERT INTO rating(song_id, rating, description) 
	VALUES ($1, $2, $3) 
	RETURNING id
	`
	if err := r.client.QueryRow(ctx, q, rating.SongID, rating.Rating, rating.Description).Scan(&rating.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); !ok {
			newErr := fmt.Errorf("SQL error: %s, details: %s, where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) (rat []rating.Rating, err error) {
	q := `SELECT id, song_id, rating, description, created_at FROM rating`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error querying songs: %w", err)
	}
	defer rows.Close()

	ratings := make([]rating.Rating, 0)

	for rows.Next() {
		var rat rating.Rating
		err := rows.Scan(&rat.ID, &rat.SongID, &rat.Rating, &rat.Description, &rat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}
		ratings = append(ratings, rat)
	}
	return ratings, nil
}

func (r *repository) FindAllBySongID(ctx context.Context, songID int) ([]rating.Rating, error) {
	q := `SELECT id, song_id, rating, description, created_at FROM rating WHERE song_id = $1`
	rows, err := r.client.Query(ctx, q, songID)
	if err != nil {
		return nil, fmt.Errorf("error querying ratings: %w", err)
	}
	defer rows.Close()

	ratings := make([]rating.Rating, 0)
	for rows.Next() {
		var rat rating.Rating
		err := rows.Scan(&rat.ID, &rat.SongID, &rat.Rating, &rat.Description, &rat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}
		ratings = append(ratings, rat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return ratings, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (rating.Rating, error) {
	q := `SELECT id, song_id, rating, description, created_at FROM rating WHERE id = $1`
	row := r.client.QueryRow(ctx, q, id)
	var rat rating.Rating
	err := row.Scan(&rat.ID, &rat.SongID, &rat.Rating, &rat.Description, &rat.CreatedAt)
	if err != nil {
		return rating.Rating{}, fmt.Errorf("error scanning row: %w", err)
	}
	return rat, nil
}

func (r *repository) Update(ctx context.Context, song rating.Rating) error {
	q := `UPDATE rating SET song_id = $1, rating = $2, description = $3 WHERE id = $4`
	_, err := r.client.Exec(ctx, q, song.SongID, song.Rating, song.Description, song.ID)
	if err != nil {
		return fmt.Errorf("error updating rating: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM rating WHERE id = $1`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("error deleting rating: %w", err)
	}
	return nil
}

func NewRepository(client postgresql.Client) rating.Repository {
	return &repository{client}
}
