package rating

import "context"

type Repository interface {
	Create(ctx context.Context, rat *Rating) error
	FindAll(ctx context.Context) (rat []Rating, err error)
	FindAllBySongID(ctx context.Context, songID int) ([]Rating, error)
	FindOne(ctx context.Context, id int) (Rating, error)
	Update(ctx context.Context, rat Rating) error
	Delete(ctx context.Context, id int) error
}
