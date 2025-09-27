package v1dto

import (
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
)

// Data Transfer Object (DTO): chuyển dữ liệu ra response
type UserDTO struct {
	UUID      string `json:"uuid"`
	Name      string `json:"full_name"`
	Email     string `json:"email_address"`
	Age       int    `json:"age"`
	Status    string `json:"status"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gt=18"`
	Password string `json:"password" binding:"required,min=8,password_strong"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}

type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Age      int    `json:"age" binding:"omitempty,gt=18"`
	Password string `json:"password" binding:"omitempty,min=8,password_strong"`
	Status   int    `json:"status" binding:"omitempty,oneof=1 2"`
	Level    int    `json:"level" binding:"omitempty,oneof=1 2"`
}

// Request
func (input *CreateUserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserEmail:    input.Email,
		UserPassword: input.Password,
		UserFullname: input.Name,
		UserAge:      utils.ConvertToInt32Pointer(input.Age),
		UserStatus:   int32(input.Status),
		UserLevel:    int32(input.Level),
	}
}

func (input *UpdateUserInput) MapUpdateInputToModel() {

}

// Response
func MapUserToDTO(user sqlc.User) *UserDTO {
	dto := &UserDTO{
		UUID:      user.UserUuid.String(),
		Name:      user.UserFullname,
		Email:     user.UserEmail,
		Status:    mapStatusText(int(user.UserStatus)),
		Level:     mapLevelText(int(user.UserLevel)),
		CreatedAt: user.UserCreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.UserAge != nil {
		dto.Age = int(*user.UserAge)
	}

	return dto
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}

func mapLevelText(status int) string {
	switch status {
	case 1:
		return "Admin"
	case 2:
		return "Member"
	default:
		return "None"
	}
}
