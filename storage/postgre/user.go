package postgre

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/user-service/model"
	"github.com/zhayt/user-service/model/dto"
	"go.uber.org/zap"
)

type UserStorage struct {
	db *sqlx.DB
	l  *zap.Logger
}

func (r *UserStorage) CreateUser(ctx context.Context, user *model.User) (uint64, error) {
	qr := `INSERT INTO web_user (name, email, password) VALUES ($1, $2, $3) RETURNING id`

	var userID uint64
	if err := r.db.GetContext(ctx, &userID, qr, user.Name, user.Email, user.Password); err != nil {
		return 0, fmt.Errorf("cannot create user: %w", err)
	}

	return userID, nil
}

func (r *UserStorage) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	qr := `SELECT id, name, email, password FROM web_user WHERE id = $1`

	var user model.User

	if err := r.db.GetContext(ctx, &user, qr, id); err != nil {
		return nil, fmt.Errorf("cannot get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserStorage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	qr := `SELECT id, name, email, password FROM web_user WHERE email = $1`

	var user model.User

	if err := r.db.GetContext(ctx, &user, qr, email); err != nil {
		return nil, fmt.Errorf("cannot get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserStorage) UpdateUserPassword(ctx context.Context, user *dto.ChangeUserPasswordDTO) error {
	qr := `UPDATE web_user SET password = $1 WHERE email = $2`

	if _, err := r.db.ExecContext(ctx, qr, user.NewPassword, user.Email); err != nil {
		return fmt.Errorf("cannot update user password: %w", err)
	}

	return nil
}

func (r *UserStorage) UpdateUserName(ctx context.Context, user *dto.ChangeUserNameDTO) error {
	qr := `UPDATE web_user SET name = $1 WHERE email = $2`

	if _, err := r.db.ExecContext(ctx, qr, user.Name, user.Email); err != nil {
		return fmt.Errorf("cannot update user name: %w", err)
	}

	return nil
}

func NewUserStorage(db *sqlx.DB, l *zap.Logger) *UserStorage {
	return &UserStorage{db: db, l: l}
}
