package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/user-service/model"
	"github.com/zhayt/user-service/model/dto"
	"github.com/zhayt/user-service/storage/postgre"
	"go.uber.org/zap"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *model.User) (uint64, error)
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUserPassword(ctx context.Context, user *dto.ChangeUserPasswordDTO) error
	UpdateUserName(ctx context.Context, user *dto.ChangeUserNameDTO) error
}

type Storage struct {
	IStorage
}

func NewStorage(db *sqlx.DB, l *zap.Logger) *Storage {
	userStorage := postgre.NewUserStorage(db, l)
	return &Storage{userStorage}
}
