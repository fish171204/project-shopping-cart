package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	FindByEmail(email string)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(uuid string)
}
