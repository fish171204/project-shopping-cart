package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	FindByEmail(email string)
	Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDelete(ctx context.Context, uuid uuid.UUID)
	Restore(ctx context.Context, uuid uuid.UUID)
	Delete(ctx context.Context, uuid uuid.UUID)
}
