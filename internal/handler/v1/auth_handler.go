package v1handler

import (
	v1service "user-management-api/internal/service/v1"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {

}

func (ah *AuthHandler) Logout(ctx *gin.Context) {

}
