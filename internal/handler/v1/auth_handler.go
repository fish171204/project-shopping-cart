package v1handler

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {

}

func (ah *AuthHandler) Logout(ctx *gin.Context) {

}
