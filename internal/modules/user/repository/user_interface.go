package repository

import (
	"GeoService/internal/models"
	"context"
)

type Userer interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserById(ctx context.Context, userId int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
}
