package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache *cache.RedisCacheService
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"role"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"` // TTL
	Revoked   bool      `json:"revoked"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", ""))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", ""))
)

const (
	AccessTokenTTL  = 15 * time.Second
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(cache *cache.RedisCacheService) TokenService {
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET environment variable is required and cannot be empty")
	}
	return &JWTService{
		cache: cache,
	}
}

/************** Access Token ***************/
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
	if err != nil {
		return "", err
	}

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

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid claims", utils.ErrCodeUnauthorized)
	}

	return token, claims, nil
}

func (js *JWTService) DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError("Cannot parse token", utils.ErrCodeInternal, err)
	}

	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Encoded data not found", utils.ErrCodeUnauthorized)
	}

	decryptedBytes, err := utils.DecryptAES(encryptedData, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError("Cannot decode data", utils.ErrCodeInternal, err)
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(decryptedBytes, &payload); err != nil {
		return nil, utils.WrapError("Invalid data format", utils.ErrCodeInternal, err)
	}

	return &payload, nil
}

/************** Refresh Token ***************/
func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {

	// Tạo mã random dưới dạng bytes
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidateRefreshToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	// Kiểm tra đã hủy chưa, hết hạn chưa refreshToken.Revoked == true
	if err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError("Cannot get refresh token", utils.ErrCodeInternal, err)
	}

	return refreshToken, nil
}

// Revoke old refresh token
func (js *JWTService) RevokeRefreshToken(token string) error {
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil {
		return utils.WrapError("Cannot get refresh token", utils.ErrCodeInternal, err)
	}

	refreshToken.Revoked = true

	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))

}
