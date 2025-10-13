package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"
)

type SqlUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

// GET
func (ur *SqlUserRepository) FindAll() {}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) FindByEmail(email string) {}

// POST
func (ur *SqlUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

// PUT
func (ur *SqlUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.db.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

// DELETE
func (ur *SqlUserRepository) SoftDelete(uuid string) {}

func (ur *SqlUserRepository) Restore(uuid string) {}

func (ur *SqlUserRepository) Delete(uuid string) {}
