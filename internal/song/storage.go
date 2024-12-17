package song

import "context"

type Repository interface {
	Create(ctx context.Context, song *Song) error
	FindAll(ctx context.Context) (u []Song, err error)
	FindOne(ctx context.Context, id int) (Song, error)
	Update(ctx context.Context, song Song) error
	Delete(ctx context.Context, id int) error
}
