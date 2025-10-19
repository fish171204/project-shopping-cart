package v1service

import (
	"user-management-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		userRepo: repo,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) error {

	return nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
