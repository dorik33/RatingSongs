package songdb

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"restApi/internal/song"
	"restApi/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, song *song.Song) error {
	q := `
	INSERT INTO song(title, artist, album, release_date) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`
	if err := r.client.QueryRow(ctx, q, song.Title, song.Artist, song.Album, song.ReleaseDate).Scan(&song.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); !ok {
			newErr := fmt.Errorf("SQL error: %s, details: %s, where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]song.Song, error) {
	q := `SELECT id, title, artist, album, release_date FROM song`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error querying songs: %w", err)
	}
	defer rows.Close()

	songs := make([]song.Song, 0)
	for rows.Next() {
		var s song.Song
		err := rows.Scan(&s.ID, &s.Title, &s.Artist, &s.Album, &s.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		songs = append(songs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return songs, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (song.Song, error) {
	q := `SELECT id, title, artist, album, release_date FROM song WHERE id = $1`
	row := r.client.QueryRow(ctx, q, id)
	var s song.Song
	err := row.Scan(&s.ID, &s.Title, &s.Artist, &s.Album, &s.ReleaseDate)
	if err != nil {
		return song.Song{}, err
	}
	return s, nil
}

func (r *repository) Update(ctx context.Context, s song.Song) error {
	q := `
	UPDATE song 
	SET title = $1, artist = $2, album = $3, release_date = $4
	WHERE id = $5
 	`
	res, err := r.client.Exec(ctx, q, s.Title, s.Artist, s.Album, s.ReleaseDate, s.ID)
	if err != nil {
		return fmt.Errorf("error updating song: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("song with ID %d not found", s.ID)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM song WHERE id = $1`
	res, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("error deleting song: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("song with ID %d not found", id)
	}

	return nil
}

func NewRepository(client postgresql.Client) song.Repository {
	return &repository{client}
}
