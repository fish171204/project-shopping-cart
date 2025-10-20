package auth

import (
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"role"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", ""))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", ""))
)

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() *JWTService {
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET environment variable is required and cannot be empty")
	}
	return &JWTService{}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := &EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserLevel,
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.EncryptAES(rawData, jwtEncryptKey)

	// claims := &Claims{
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ID:        uuid.NewString(),                                   // jti - Unique token ID
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)), // exp - Expiration time (current + 15 minutes)
	// 		IssuedAt:  jwt.NewNumericDate(time.Now()),                     // iat - Token issued time
	// 		Issuer:    "nguyen-dang-khoa",                                 // iss - Token issuer
	// 	},
	// }

	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		"iat":  jwt.NewNumericDate(time.Now()),
		"iss":  "nguyen-dang-khoa",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (js *JWTService) GenerateRefreshToken() {

}
