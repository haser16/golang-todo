package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/haser16/golang-todo/internal/core/domain"
	core_errors "github.com/haser16/golang-todo/internal/core/errors"
	core_postgres_pool "github.com/haser16/golang-todo/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeOut())
	defer cancel()

	query := `SELECT id, version, full_name, phone_number FROM todoapp.users WHERE id = $1;`
	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan errors: %w", err)
	}
	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil
}
