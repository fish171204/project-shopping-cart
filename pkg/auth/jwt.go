package auth

import (
	"time"
	"user-management-api/internal/db/sqlc"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}

type Claims struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) {
	claims := &Claims{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     string(rune(user.UserLevel)),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),                                   // jti - Unique token ID
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)), // exp - Expiration time (current + 15 minutes)
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // iat - Token issued time
			Issuer:    "nguyen-dang-khoa",                                 // iss - Token issuer
		},
	}
}

func (js *JWTService) GenerateRefreshToken() {

}
