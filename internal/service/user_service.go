package service

import (
	"strings"
	"user-management-api/internal/models"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(search string, page, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError("failed to fetch users", utils.ErrCodeInternal, err)
	}

	var filteredUsers []models.User
	// Search
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			if strings.Contains(name, search) || strings.Contains(email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		// No search results
		filteredUsers = users
	}

	start := (page - 1) * limit
	if start >= len(filteredUsers) {
		return []models.User{}, nil
	}

	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}

func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	return user, nil
}

// POST
func (us *userService) CreateUsers(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exists := us.repo.FindByEmail(user.Email); exists {
		return models.User{}, utils.NewError("email already exists", utils.ErrCodeConflict)
	}

	user.UUID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
	}

	user.Password = string(hashedPassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError("failed to create password", utils.ErrCodeInternal, err)
	}

	return user, nil

}

// PUT
func (us *userService) UpdateUser(uuid string, updatedUser models.User) (models.User, error) {
	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	// Update fields conditionally
	utils.UpdateStringField(&currentUser.Name, updatedUser.Name)
	utils.UpdateIntField(&currentUser.Age, updatedUser.Age)
	utils.UpdateIntField(&currentUser.Status, updatedUser.Status)
	utils.UpdateIntField(&currentUser.Level, updatedUser.Level)

	// Handle email with conflict check
	if updatedUser.Email != "" {
		normalizedEmail := utils.NormalizeString(updatedUser.Email)
		if normalizedEmail != currentUser.Email {
			if u, exists := us.repo.FindByEmail(normalizedEmail); exists && u.UUID != uuid {
				return models.User{}, utils.NewError("email already exists", utils.ErrCodeConflict)
			}
		}
		currentUser.Email = normalizedEmail
	}

	// Handle password with hashing
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
		}
		currentUser.Password = string(hashedPassword)
	}

	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError("failed to update user", utils.ErrCodeInternal, err)
	}

	return currentUser, nil
}

func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError("failed to delete user", utils.ErrCodeInternal, err)

	}
	return nil
}
