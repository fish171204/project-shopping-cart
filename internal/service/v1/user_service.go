package v1service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo  repository.UserRepository
	cache cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:  repo,
		cache: cache.NewRedisCacheService(redisClient),
	}
}

// GET V1
func (us *userService) GetAllUsers(ctx *gin.Context, search, orderBy, sort string, page, limit int32) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	if sort == "" {
		sort = "desc"
	}

	if orderBy == "" {
		orderBy = "user_created_at"
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitInt := utils.GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(limitInt)
	}

	offset := (page - 1) * limit

	/** Get Cache Redis **/
	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}

	if err := us.cache.Get("getAllUsers", &cacheData); err == nil && cacheData.Users != nil {
		return cacheData.Users, cacheData.Total, nil
	}

	users, err := us.repo.GetAll(context, search, orderBy, sort, limit, offset)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("failed to fetch users", utils.ErrCodeInternal, err)
	}

	total, err := us.repo.CountUsers(context, search, false)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("failed to count users", utils.ErrCodeInternal, err)
	}

	// Create cache data
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: users,
		Total: int32(total),
	}
	us.cache.Set("getAllUsers", cacheData, 5*time.Second)

	return users, int32(total), nil
}

// GET V2
func (us *userService) GetAllUsersV2(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	/** Get Cache Redis **/
	cacheKey := us.generateCacheKey(search, orderBy, sort, page, limit, deleted)

	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}

	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Users != nil {
		return cacheData.Users, cacheData.Total, nil
	}

	if sort == "" {
		sort = "desc"
	}

	if orderBy == "" {
		orderBy = "user_created_at"
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitInt := utils.GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(limitInt)
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("failed to fetch users", utils.ErrCodeInternal, err)
	}

	total, err := us.repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("failed to count users", utils.ErrCodeInternal, err)
	}

	// Create cache data
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: users,
		Total: int32(total),
	}
	us.cache.Set(cacheKey, cacheData, 10*time.Minute)

	return users, int32(total), nil
}

func (us *userService) GetUserByUuid(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.GetByUuid(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("failed to get an user", utils.ErrCodeInternal, err)
	}

	return user, nil
}

// POST
func (us *userService) CreateUsers(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	input.UserEmail = utils.NormalizeString(input.UserEmail)

	hashesPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
	}

	input.UserPassword = string(hashesPassword)

	user, err := us.repo.Create(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
		}

		return sqlc.User{}, utils.WrapError("failed to create a new user", utils.ErrCodeInternal, err)
	}

	// Clear cache redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	return user, nil
}

// PUT
func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashesPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
		}

		hashed := string(hashesPassword)
		input.UserPassword = &hashed
	}

	updatedUser, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("failed to update user", utils.ErrCodeInternal, err)
	}

	// Clear cache redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	return updatedUser, nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	softDeleteUser, err := us.repo.SoftDelete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("failed to delete user", utils.ErrCodeInternal, err)
	}

	// Clear cache redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	return softDeleteUser, nil
}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	restoreUser, err := us.repo.Restore(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found or not marked as delete for restore", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("failed to restore user", utils.ErrCodeInternal, err)
	}

	// Clear cache redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	return restoreUser, nil
}

func (us *userService) DeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	_, err := us.repo.Delete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("user not found or not marked as delete for permenent removal", utils.ErrCodeNotFound)
		}
		return utils.WrapError("failed to restore user", utils.ErrCodeInternal, err)
	}

	// Clear cache redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	return nil
}

func (us *userService) generateCacheKey(search, orderBy, sort string, page, limit int32, deleted bool) string {
	search = strings.TrimSpace(search)
	if search == "" {
		search = "none"
	}

	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		orderBy = "user_created_at"
	}

	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}

	return fmt.Sprintf("users:%s:%s:%s:%d:%d:%t", search, orderBy, sort, page, limit, deleted)
}
