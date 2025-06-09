package repository

import (
	"GeoService/internal/models"
	"context"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) Userer {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	var id int64
	_, err := r.db.NewInsert().Model(user).Returning("id").Exec(ctx, &id)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (r *UserRepository) GetUserById(ctx context.Context, userId int) (*models.User, error) {
	var user *models.User
	err := r.db.NewSelect().Model(user).Where("id = ?", userId).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)

	return user, err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	return user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewDelete().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
