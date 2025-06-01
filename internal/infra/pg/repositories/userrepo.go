package repositories

import (
	"context"

	"github.com/victorsvart/gofin/internal/domain/apperror"
	"github.com/victorsvart/gofin/internal/domain/structs/user"
	"github.com/victorsvart/gofin/internal/infra/pg"
)

type userRepositoryImpl struct {
	pg pg.PostgresConn
}

func NewUserRepository(pg pg.PostgresConn) user.UserRepository {
	return &userRepositoryImpl{pg}
}

func (u *userRepositoryImpl) InsertUser(ctx context.Context, c user.CreateUserInput) (*user.CreateUserResponse, *apperror.AppError) {
	qb := u.pg.Qb
	query, args, err := qb.Insert("users").
		Columns("sub", "name", "username", "email_address").
		Values(c.ID, c.Name, c.Username, c.EmailAddress).
		ToSql()
	if err != nil {
		return nil, apperror.NewAppError(
			c,
			apperror.QUERY,
			apperror.INTERNAL,
			err,
		)
	}

	_, err = u.pg.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.NewAppError(
			c,
			apperror.DB,
			apperror.INTERNAL,
			err,
		)
	}

	return &user.CreateUserResponse{
		ID:       c.ID,
		Name:     c.Name,
		Username: c.Username,
	}, nil
}
