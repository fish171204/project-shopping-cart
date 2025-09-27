package repository

type SqlUserRepository struct {
}

func NewSqlUserRepository() UserRepository {
	return &SqlUserRepository{}
}

// GET
func (ur *SqlUserRepository) FindAll() {}

func (ur *SqlUserRepository) FindByUUID(uuid string) {}

func (ur *SqlUserRepository) FindByEmail(email string) {}

// POST
func (ur *SqlUserRepository) Create() {

}

// PUT
func (ur *SqlUserRepository) Update(uuid string) {}

// DELETE
func (ur *SqlUserRepository) Delete(uuid string) {}
