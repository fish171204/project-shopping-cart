package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	FindByEmail(email string)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	Update(uuid string)
	Delete(uuid string)
}
