package v1service

import (
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

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
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)

	return nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
