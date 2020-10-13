package user

import "context"

type User struct {
	ID         string
	Name       string
	PictureURL string
}

type Repository interface {
	GetByID(ctx context.Context, userID string) (User, error)
}
