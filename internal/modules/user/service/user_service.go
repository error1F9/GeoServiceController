package service

import (
	"GeoService/internal/infrastucture/errors"
	"GeoService/internal/infrastucture/tools/cryptography"
	"GeoService/internal/models"
	"GeoService/internal/modules/user/repository"
	"context"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserService struct {
	repository repository.Userer
	logger     *zap.Logger
}

func NewUserService(repository repository.Userer, logger *zap.Logger) *UserService {
	return &UserService{repository: repository, logger: logger}
}

func (u *UserService) Create(ctx context.Context, in UserCreateIn) UserCreateOut {
	user := &models.User{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
		Role:     in.Role,
	}
	userID, err := u.repository.CreateUser(ctx, user)
	u.logger.Info("userID ", zap.Any("user_id", userID), zap.Any("email", user.Email))
	if err != nil {
		if v, ok := err.(*pq.Error); ok && v.Code == "23505" {
			return UserCreateOut{
				ErrorCode: errors.UserServiceUserAlreadyExists,
			}
		}
		return UserCreateOut{
			ErrorCode: errors.UserServiceCreateUserErr,
		}
	}
	return UserCreateOut{
		UserID: userID,
	}
}

func (u *UserService) Update(ctx context.Context, in UserUpdateIn) UserUpdateOut {
	userID := in.User.Id
	user, err := u.repository.GetUserById(ctx, int(userID))
	if err != nil {
		return UserUpdateOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}

	err = u.repository.UpdateUser(ctx, user)
	if err != nil {
		return UserUpdateOut{
			ErrorCode: errors.UserServiceUpdateErr,
		}
	}
	return UserUpdateOut{
		Success: true,
	}
}

func (u *UserService) GetById(ctx context.Context, in GetByIDIn) UserOut {
	user, err := u.repository.GetUserById(ctx, in.Id)
	if err != nil {
		u.logger.Error("user: getById error", zap.Error(err))
		return UserOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}
	return UserOut{
		User: *user,
	}
}

func (u *UserService) GetByUsername(ctx context.Context, in GetByUsernameIn) UserOut {
	user, err := u.repository.GetUserByUsername(ctx, in.Username)
	if err != nil {
		u.logger.Error("user: getByUsername error", zap.Error(err))
		return UserOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}
	return UserOut{
		User: *user,
	}
}

func (u *UserService) GetByEmail(ctx context.Context, in GetByEmailIn) UserOut {
	user, err := u.repository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		u.logger.Error("user: getByEmail error", zap.Error(err))
		return UserOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}
	return UserOut{
		User: *user,
	}
}

func (u *UserService) ChangePassword(ctx context.Context, in ChangePasswordIn) ChangePasswordOut {
	user, err := u.repository.GetUserById(ctx, in.UserID)
	if err != nil {
		u.logger.Error("user: getById error", zap.Error(err))
		return ChangePasswordOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}

	if cryptography.CheckPassword(user.Password, in.OldPassword) {
		user.Password = in.NewPassword
	} else {
		return ChangePasswordOut{
			ErrorCode: errors.UserServiceWrongPasswordErr,
		}
	}

	err = u.repository.UpdateUser(ctx, user)
	if err != nil {
		u.logger.Error("user: changePassword error", zap.Error(err))
		return ChangePasswordOut{
			ErrorCode: errors.UserServiceUpdateErr,
		}
	}
	return ChangePasswordOut{
		Success: true,
	}
}

func (u *UserService) Delete(ctx context.Context, in UserDeleteIn) UserDeleteOut {
	user, err := u.repository.GetUserById(ctx, in.UserID)
	if err != nil {
		u.logger.Error("user: getById error", zap.Error(err))
		return UserDeleteOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}

	err = u.repository.DeleteUser(ctx, user)
	if err != nil {
		u.logger.Error("user: delete error", zap.Error(err))
		return UserDeleteOut{
			ErrorCode: errors.UserServiceDeleteUserErr,
		}
	}
	return UserDeleteOut{
		Success: true,
	}
}
