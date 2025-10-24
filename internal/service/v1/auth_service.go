package v1service

import (
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	// Compare hashed password in database with pass input
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError("Unable to create access token", utils.ErrCodeInternal, err)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError("Unable to create refresh token", utils.ErrCodeInternal, err)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError("Cannot save refresh token", utils.ErrCodeInternal, err)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {
	context := ctx.Request.Context()

	// Kiểm tra refresh token. Trả thông tin user_uuid
	token, err := as.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	// Lấy thông tin user
	userUuid, _ := uuid.Parse(token.UserUUID)
	user, err := as.userRepo.GetByUuid(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// Tạo access token mới
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError("Unable to create access token", utils.ErrCodeInternal, err)
	}

	// Tạo refresh token mới
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError("Unable to create refresh token", utils.ErrCodeInternal, err)
	}

	// Vô hiệu hóa refresh token cũ

	// Lưu refresh token mới
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError("Cannot save refresh token", utils.ErrCodeInternal, err)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
