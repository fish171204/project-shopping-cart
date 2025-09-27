package repository

import "user-management-api/internal/db/sqlc"

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
func (ur *SqlUserRepository) Create() {
	ur.db.CreateUser()
}

// PUT
func (ur *SqlUserRepository) Update(uuid string) {}

// DELETE
func (ur *SqlUserRepository) Delete(uuid string) {}
