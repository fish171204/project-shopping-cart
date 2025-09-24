package repository

type InMemoryUserRepository struct {
}

func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{}
}

// GET
func (ur *InMemoryUserRepository) FindAll() {}

func (ur *InMemoryUserRepository) FindByUUID(uuid string) {}

func (ur *InMemoryUserRepository) FindByEmail(email string) {}

// POST
func (ur *InMemoryUserRepository) Create() {}

// PUT
func (ur *InMemoryUserRepository) Update(uuid string) {}

// DELETE
func (ur *InMemoryUserRepository) Delete(uuid string) {}
