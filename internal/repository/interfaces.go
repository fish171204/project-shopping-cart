package repository

type UserRepository interface {
	FindAll()
	FindByUUID(uuid string)
	FindByEmail(email string)
	Create()
	Update(uuid string)
	Delete(uuid string)
}
